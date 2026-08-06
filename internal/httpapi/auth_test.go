package httpapi

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"simon-jp-api/internal/models"
	"simon-jp-api/internal/service"
)

type testUserStore struct {
	byUsername map[string]*models.User
	byID       map[int64]*models.User
}

func (f *testUserStore) FindByUsername(_ context.Context, username string) (*models.User, error) {
	u, ok := f.byUsername[username]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (f *testUserStore) FindByID(_ context.Context, id int64) (*models.User, error) {
	u, ok := f.byID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

type testSessionStore struct {
	sessions map[string]int64
	ttl      time.Duration
}

func (f *testSessionStore) Create(_ context.Context, token string, userID int64) error {
	f.sessions[token] = userID
	return nil
}

func (f *testSessionStore) Get(_ context.Context, token string) (int64, error) {
	userID, ok := f.sessions[token]
	if !ok {
		return 0, service.ErrSessionNotFound
	}
	return userID, nil
}

func (f *testSessionStore) Delete(_ context.Context, token string) error {
	delete(f.sessions, token)
	return nil
}

func (f *testSessionStore) TTL() time.Duration { return f.ttl }

type testThrottler struct {
	locked map[string]bool
}

func (f *testThrottler) Check(_ context.Context, key string) error {
	if f.locked[key] {
		return service.ErrTooManyAttempts
	}
	return nil
}

func (f *testThrottler) RecordFailure(_ context.Context, _ string) error { return nil }
func (f *testThrottler) Reset(_ context.Context, _ string) error         { return nil }

type testPegawaiStore struct {
	byUser map[int64]*models.Pegawai
}

func (f *testPegawaiStore) FindByUserID(_ context.Context, userID int64) (*models.Pegawai, error) {
	p, ok := f.byUser[userID]
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}

func newTestRouter(t *testing.T) (http.Handler, *testSessionStore) {
	t.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	user := &models.User{
		BaseID:   models.BaseID{ID: 1, UUID: "11111111-1111-1111-1111-111111111111"},
		Username: "alice",
		Password: string(hash),
		Role:     "pegawai",
	}
	users := &testUserStore{
		byUsername: map[string]*models.User{"alice": user},
		byID:       map[int64]*models.User{user.ID: user},
	}
	pegawai := &testPegawaiStore{byUser: map[int64]*models.Pegawai{
		user.ID: {BaseID: models.BaseID{UUID: "22222222-2222-2222-2222-222222222222"}, Nama: "Alice"},
	}}
	sessions := &testSessionStore{sessions: make(map[string]int64), ttl: time.Hour}
	throttle := &testThrottler{locked: make(map[string]bool)}
	auth := service.NewAuthService(users, pegawai, sessions, throttle)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewRouter(Deps{Auth: auth}, logger), sessions
}

func doRequest(t *testing.T, handler http.Handler, method, path, body, token string) *httptest.ResponseRecorder {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = bytes.NewBufferString(body)
	}
	req := httptest.NewRequest(method, path, reader)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestLoginHandlerSuccess(t *testing.T) {
	handler, _ := newTestRouter(t)
	rec := doRequest(t, handler, http.MethodPost, "/api/login",
		`{"username":"alice","password":"secret"}`, "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if rec.Body.String() == "" || rec.Body.String()[0] != '{' {
		t.Fatalf("expected JSON response, got %q", rec.Body.String())
	}
}

func TestLoginHandlerInvalidCredentials(t *testing.T) {
	handler, _ := newTestRouter(t)
	rec := doRequest(t, handler, http.MethodPost, "/api/login",
		`{"username":"alice","password":"wrong"}`, "")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401; body=%s", rec.Code, rec.Body.String())
	}
}

func TestLoginHandlerBadRequest(t *testing.T) {
	handler, _ := newTestRouter(t)
	rec := doRequest(t, handler, http.MethodPost, "/api/login", `not-json`, "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
}

func TestLoginHandlerMissingFields(t *testing.T) {
	handler, _ := newTestRouter(t)
	rec := doRequest(t, handler, http.MethodPost, "/api/login", `{"username":""}`, "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
}

func TestLoginHandlerLocked(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	user := &models.User{
		BaseID:   models.BaseID{ID: 1, UUID: "11111111-1111-1111-1111-111111111111"},
		Username: "alice",
		Password: string(hash),
		Role:     "pegawai",
	}
	users := &testUserStore{
		byUsername: map[string]*models.User{"alice": user},
		byID:       map[int64]*models.User{user.ID: user},
	}
	pegawai := &testPegawaiStore{byUser: map[int64]*models.Pegawai{user.ID: {Nama: "Alice"}}}
	sessions := &testSessionStore{sessions: make(map[string]int64), ttl: time.Hour}
	throttle := &testThrottler{locked: map[string]bool{"login:alice": true}}
	auth := service.NewAuthService(users, pegawai, sessions, throttle)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewRouter(Deps{Auth: auth}, logger)

	rec := doRequest(t, handler, http.MethodPost, "/api/login",
		`{"username":"alice","password":"secret"}`, "")
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429; body=%s", rec.Code, rec.Body.String())
	}
}

func TestMeHandler(t *testing.T) {
	handler, sessions := newTestRouter(t)
	login := doRequest(t, handler, http.MethodPost, "/api/login",
		`{"username":"alice","password":"secret"}`, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login failed: %d", login.Code)
	}

	// extract token by calling auth service directly
	_ = sessions
	token := ""
	for tok := range sessions.sessions {
		token = tok
		break
	}

	rec := doRequest(t, handler, http.MethodGet, "/api/user", "", token)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestMeHandlerNoToken(t *testing.T) {
	handler, _ := newTestRouter(t)
	rec := doRequest(t, handler, http.MethodGet, "/api/user", "", "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestLogoutHandler(t *testing.T) {
	handler, sessions := newTestRouter(t)
	doRequest(t, handler, http.MethodPost, "/api/login",
		`{"username":"alice","password":"secret"}`, "")

	token := ""
	for tok := range sessions.sessions {
		token = tok
		break
	}

	rec := doRequest(t, handler, http.MethodPost, "/api/logout", "", token)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}
