package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"simon-jp-api/internal/models"
)

type contextKey int

const (
	userContextKey contextKey = iota
	tokenContextKey
)

func withAuth(r *http.Request, user *models.User, token string) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	ctx = context.WithValue(ctx, tokenContextKey, token)
	return r.WithContext(ctx)
}

func userFrom(r *http.Request) *models.User {
	user, _ := r.Context().Value(userContextKey).(*models.User)
	return user
}

func tokenFrom(r *http.Request) string {
	token, _ := r.Context().Value(tokenContextKey).(string)
	return token
}

func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			logger.Info("http",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", middleware.GetReqID(r.Context()),
			)
		})
	}
}
