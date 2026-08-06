package httpapi

import (
	"net/http"

	"simon-jp-api/internal/service"
)

type MateriPpmHandler struct {
	svc *service.MateriPpmService
}

func NewMateriPpmHandler(svc *service.MateriPpmService) *MateriPpmHandler {
	return &MateriPpmHandler{svc: svc}
}

func (h *MateriPpmHandler) List(w http.ResponseWriter, r *http.Request) {
	page, perPage := parsePagination(r)
	f := filtersFromRequest(r,
		[]string{"materi_pengembangan", "nama_pemateri", "nomor_surat"},
		"tanggal_pelaksanaan", "tanggal_pelaksanaan", "tanggal_pelaksanaan")

	items, total, err := h.svc.List(r.Context(), f, page, perPage)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	data := make([]materiPpmListItem, 0, len(items))
	for _, m := range items {
		data = append(data, newMateriPpmListItem(m))
	}
	writeJSON(w, http.StatusOK, newPageResponse(data, page, perPage, total))
}

func (h *MateriPpmHandler) Show(w http.ResponseWriter, r *http.Request) {
	m, err := h.svc.Get(r.Context(), routeParam(r, "uuid"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": newMateriPpmShow(m)})
}
