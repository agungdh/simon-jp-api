package httpapi

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"simon-jp-api/internal/service"
)

type DownloadHandler struct {
	svc *service.DownloadService
}

func NewDownloadHandler(svc *service.DownloadService) *DownloadHandler {
	return &DownloadHandler{svc: svc}
}

func (h *DownloadHandler) Download(w http.ResponseWriter, r *http.Request) {
	pegawai := pegawaiFrom(r)
	if pegawai == nil {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}

	module := chi.URLParam(r, "module")
	uuid := chi.URLParam(r, "uuid")

	result, err := h.svc.Resolve(r.Context(), module, uuid, pegawai.ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrModuleNotFound):
			writeError(w, http.StatusNotFound, "module not found")
		case errors.Is(err, service.ErrNotFound):
			writeError(w, http.StatusNotFound, "record not found")
		case errors.Is(err, service.ErrNoFileAttached):
			writeError(w, http.StatusNotFound, "no file attached")
		case errors.Is(err, service.ErrFileNotFound):
			writeError(w, http.StatusNotFound, "file not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	http.Redirect(w, r, result.URL, http.StatusFound)
}
