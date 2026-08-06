package httpapi

import (
	"net/http"

	"simon-jp-api/internal/service"
)

type DiklatHandler struct {
	svc *service.DiklatService
}

func NewDiklatHandler(svc *service.DiklatService) *DiklatHandler {
	return &DiklatHandler{svc: svc}
}

func (h *DiklatHandler) List(w http.ResponseWriter, r *http.Request) {
	pegawai := pegawaiFrom(r)
	if pegawai == nil {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	page, perPage := parsePagination(r)
	f := filtersFromRequest(r,
		[]string{"nomor_surat", "materi_pengembangan"},
		"dari_tanggal_pelaksanaan", "sampai_tanggal_pelaksanaan", "dari_tanggal_pelaksanaan")

	items, total, err := h.svc.List(r.Context(), pegawai.ID, f, page, perPage)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	data := make([]diklatListItem, 0, len(items))
	for _, d := range items {
		data = append(data, newDiklatListItem(d))
	}
	writeJSON(w, http.StatusOK, newPageResponse(data, page, perPage, total))
}

func (h *DiklatHandler) Show(w http.ResponseWriter, r *http.Request) {
	pegawai := pegawaiFrom(r)
	if pegawai == nil {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	d, err := h.svc.Get(r.Context(), pegawai.ID, routeParam(r, "uuid"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": newDiklatShow(d)})
}
