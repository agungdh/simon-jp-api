package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"simon-jp-api/internal/models"
)

type fakeUserStore struct {
	byUsername map[string]*models.User
	byID       map[int64]*models.User
}

func (f *fakeUserStore) FindByUsername(_ context.Context, username string) (*models.User, error) {
	u, ok := f.byUsername[username]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (f *fakeUserStore) FindByID(_ context.Context, id int64) (*models.User, error) {
	u, ok := f.byID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

type fakeSessionStore struct {
	sessions map[string]int64
	ttl      time.Duration
}

func (f *fakeSessionStore) Create(_ context.Context, token string, userID int64) error {
	f.sessions[token] = userID
	return nil
}

func (f *fakeSessionStore) Get(_ context.Context, token string) (int64, error) {
	userID, ok := f.sessions[token]
	if !ok {
		return 0, ErrSessionNotFound
	}
	return userID, nil
}

func (f *fakeSessionStore) Delete(_ context.Context, token string) error {
	delete(f.sessions, token)
	return nil
}

func (f *fakeSessionStore) TTL() time.Duration { return f.ttl }

type fakeThrottler struct {
	failures    map[string]int
	locked      map[string]bool
	maxAttempts int
}

func (f *fakeThrottler) Check(_ context.Context, key string) error {
	if f.locked[key] {
		return ErrTooManyAttempts
	}
	return nil
}

func (f *fakeThrottler) RecordFailure(_ context.Context, key string) error {
	f.failures[key]++
	if f.failures[key] >= f.maxAttempts {
		f.locked[key] = true
	}
	return nil
}

func (f *fakeThrottler) Reset(_ context.Context, key string) error {
	delete(f.failures, key)
	delete(f.locked, key)
	return nil
}

func newTestUser(username, password string) *models.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return &models.User{
		BaseID:   models.BaseID{ID: 1, UUID: "11111111-1111-1111-1111-111111111111"},
		Username: username,
		Password: string(hash),
	}
}

func newTestAuth(t *testing.T) (*AuthService, *fakeUserStore, *fakeSessionStore, *fakeThrottler) {
	t.Helper()
	user := newTestUser("alice", "secret")
	users := &fakeUserStore{
		byUsername: map[string]*models.User{"alice": user},
		byID:       map[int64]*models.User{user.ID: user},
	}
	sessions := &fakeSessionStore{sessions: make(map[string]int64), ttl: time.Hour}
	throttle := &fakeThrottler{failures: make(map[string]int), locked: make(map[string]bool), maxAttempts: 5}
	return NewAuthService(users, sessions, throttle), users, sessions, throttle
}

func TestLoginSuccess(t *testing.T) {
	auth, _, sessions, throttle := newTestAuth(t)

	token, ttl, err := auth.Login(context.Background(), "alice", "secret")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if token == "" {
		t.Fatal("Login: expected non-empty token")
	}
	if ttl != time.Hour {
		t.Fatalf("Login: ttl = %v, want %v", ttl, time.Hour)
	}
	if _, ok := sessions.sessions[token]; !ok {
		t.Fatal("Login: session not created")
	}
	if len(throttle.failures) != 0 {
		t.Fatal("Login: throttle failures not reset on success")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	auth, _, _, throttle := newTestAuth(t)

	_, _, err := auth.Login(context.Background(), "alice", "wrong")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login: err = %v, want ErrInvalidCredentials", err)
	}
	if throttle.failures["login:alice"] != 1 {
		t.Fatalf("Login: failure not recorded, got %d", throttle.failures["login:alice"])
	}
}

func TestLoginUnknownUser(t *testing.T) {
	auth, _, _, throttle := newTestAuth(t)

	_, _, err := auth.Login(context.Background(), "nobody", "secret")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login: err = %v, want ErrInvalidCredentials", err)
	}
	if throttle.failures["login:nobody"] != 1 {
		t.Fatalf("Login: failure not recorded for unknown user")
	}
}

func TestLoginLocked(t *testing.T) {
	auth, _, _, throttle := newTestAuth(t)
	throttle.locked["login:alice"] = true

	_, _, err := auth.Login(context.Background(), "alice", "secret")
	if !errors.Is(err, ErrTooManyAttempts) {
		t.Fatalf("Login: err = %v, want ErrTooManyAttempts", err)
	}
}

func TestLoginLockedAfterMaxAttempts(t *testing.T) {
	auth, _, _, throttle := newTestAuth(t)
	key := "login:alice"

	for i := 0; i < throttle.maxAttempts; i++ {
		_ = throttle.RecordFailure(context.Background(), key)
	}

	_, _, err := auth.Login(context.Background(), "alice", "secret")
	if !errors.Is(err, ErrTooManyAttempts) {
		t.Fatalf("Login: err = %v, want ErrTooManyAttempts after max failures", err)
	}
}

func TestAuthenticate(t *testing.T) {
	auth, _, _, _ := newTestAuth(t)
	ctx := context.Background()
	token, _, _ := auth.Login(ctx, "alice", "secret")

	user, err := auth.Authenticate(ctx, token)
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if user.Username != "alice" {
		t.Fatalf("Authenticate: username = %q, want alice", user.Username)
	}
}

func TestAuthenticateInvalidToken(t *testing.T) {
	auth, _, _, _ := newTestAuth(t)

	_, err := auth.Authenticate(context.Background(), "bogus")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("Authenticate: err = %v, want ErrInvalidToken", err)
	}
}

func TestLogout(t *testing.T) {
	auth, _, sessions, _ := newTestAuth(t)
	ctx := context.Background()
	token, _, _ := auth.Login(ctx, "alice", "secret")

	if err := auth.Logout(ctx, token); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	if _, ok := sessions.sessions[token]; ok {
		t.Fatal("Logout: session still present")
	}
}

func TestLogoutInvalidToken(t *testing.T) {
	auth, _, _, _ := newTestAuth(t)

	err := auth.Logout(context.Background(), "bogus")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("Logout: err = %v, want ErrInvalidToken", err)
	}
}
