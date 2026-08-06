package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"simon-jp-api/internal/service"
)

const maxBodyBytes = 1 << 20

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := bearerToken(r)
			if token == "" {
				writeError(w, http.StatusUnauthorized, "missing token")
				return
			}

			user, err := h.auth.Authenticate(r.Context(), token)
			if err != nil {
				if errors.Is(err, service.ErrInvalidToken) {
					writeError(w, http.StatusUnauthorized, "invalid token")
					return
				}
				writeError(w, http.StatusInternalServerError, "internal server error")
				return
			}

			pegawai, _ := h.auth.Pegawai(r.Context(), user.ID)

			next.ServeHTTP(w, withAuth(r, user, pegawai, token))
		})
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string       `json:"token"`
	ExpiresIn int64        `json:"expires_in"`
	User      userResponse `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	result, err := h.auth.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTooManyAttempts):
			writeError(w, http.StatusTooManyRequests, "too many login attempts, try again later")
			return
		case errors.Is(err, service.ErrInvalidCredentials):
			writeError(w, http.StatusUnauthorized, "invalid username or password")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	writeJSON(w, http.StatusOK, loginResponse{
		Token:     result.Token,
		ExpiresIn: int64(result.ExpiresIn.Seconds()),
		User:      newUserResponse(result.User, result.Pegawai),
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.auth.Logout(r.Context(), tokenFrom(r)); err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully."})
}

func (h *AuthHandler) User(w http.ResponseWriter, r *http.Request) {
	user, pegawai := userFrom(r), pegawaiFrom(r)
	if user == nil {
		writeError(w, http.StatusUnauthorized, "invalid token")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": newUserResponse(user, pegawai)})
}

func bearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
