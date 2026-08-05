package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"simon-jp-api/internal/service"
)

func NewRouter(auth *service.AuthService) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	h := NewAuthHandler(auth)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/logout", h.Logout)
		r.Get("/me", h.Me)
	})

	return r
}
