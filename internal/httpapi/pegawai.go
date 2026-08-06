package httpapi

import (
	"net/http"
	"strconv"

	"simon-jp-api/internal/service"
)

type PegawaiHandler struct {
	svc *service.PegawaiService
}

func NewPegawaiHandler(svc *service.PegawaiService) *PegawaiHandler {
	return &PegawaiHandler{svc: svc}
}

func (h *PegawaiHandler) List(w http.ResponseWriter, r *http.Request) {
	pegawai := pegawaiFrom(r)
	if pegawai == nil {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	page, perPage := parsePagination(r)
	q := r.URL.Query()
	bidangID, _ := strconv.ParseInt(q.Get("bidang_id"), 10, 64)

	items, total, err := h.svc.ListOwn(r.Context(), pegawai.ID, q.Get("search"), q.Get("status"), bidangID, page, perPage)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	data := make([]pegawaiListItem, 0, len(items))
	for _, p := range items {
		data = append(data, newPegawaiListItem(p))
	}
	writeJSON(w, http.StatusOK, newPageResponse(data, page, perPage, total))
}

func (h *PegawaiHandler) Show(w http.ResponseWriter, r *http.Request) {
	pegawai := pegawaiFrom(r)
	if pegawai == nil {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	p, err := h.svc.GetOwn(r.Context(), pegawai.ID, routeParam(r, "uuid"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": newPegawaiShow(p)})
}
