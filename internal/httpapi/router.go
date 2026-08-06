package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"simon-jp-api/internal/service"
)

func NewRouter(auth *service.AuthService, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(requestLogger(logger))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	h := NewAuthHandler(auth)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Group(func(r chi.Router) {
			r.Use(h.Middleware())
			r.Post("/logout", h.Logout)
			r.Get("/me", h.Me)
		})
	})

	return r
}
