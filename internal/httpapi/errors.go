package httpapi

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"simon-jp-api/internal/service"
)

func routeParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	case errors.Is(err, service.ErrForbidden):
		writeError(w, http.StatusForbidden, "forbidden")
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
