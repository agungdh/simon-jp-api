package httpapi

import (
	"net/http"

	"simon-jp-api/internal/service"
)

type PpmHandler struct {
	svc *service.PpmService
}

func NewPpmHandler(svc *service.PpmService) *PpmHandler {
	return &PpmHandler{svc: svc}
}

func (h *PpmHandler) List(w http.ResponseWriter, r *http.Request) {
	pegawai := pegawaiFrom(r)
	if pegawai == nil {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	page, perPage := parsePagination(r)
	f := filtersFromRequest(r,
		[]string{"nomor_surat", "materi_pengembangan"},
		"tanggal_pelaksanaan", "tanggal_pelaksanaan", "tanggal_pelaksanaan")

	items, total, err := h.svc.List(r.Context(), pegawai.ID, f, page, perPage)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	data := make([]ppmListItem, 0, len(items))
	for _, p := range items {
		data = append(data, newPpmListItem(p))
	}
	writeJSON(w, http.StatusOK, newPageResponse(data, page, perPage, total))
}

func (h *PpmHandler) Show(w http.ResponseWriter, r *http.Request) {
	pegawai := pegawaiFrom(r)
	if pegawai == nil {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	p, err := h.svc.Get(r.Context(), pegawai.ID, routeParam(r, "uuid"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": newPpmShow(p)})
}
