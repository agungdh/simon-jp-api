package httpapi

import (
	"errors"
	"net/http"

	"simon-jp-api/internal/service"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

type detailCapaianDTO struct {
	Diklat  int `json:"diklat"`
	Ppm     int `json:"ppm"`
	Seminar int `json:"seminar"`
	Webinar int `json:"webinar"`
	Lc      int `json:"lc"`
}

type rekapitulasiDTO struct {
	Triwulan1 int `json:"triwulan_1"`
	Triwulan2 int `json:"triwulan_2"`
	Triwulan3 int `json:"triwulan_3"`
	Triwulan4 int `json:"triwulan_4"`
}

type detailItemDTO struct {
	UUID      string `json:"uuid"`
	Tipe      string `json:"tipe"`
	Tanggal   string `json:"tanggal"`
	Materi    string `json:"materi"`
	JumlahJam int    `json:"jumlah_jam"`
}

type dashboardData struct {
	Pegawai       pegawaiBrief     `json:"pegawai"`
	Tahun         int              `json:"tahun"`
	Kategori      *string          `json:"kategori"`
	JumlahMinimal int              `json:"jumlah_minimal"`
	JumlahCapaian int              `json:"jumlah_capaian"`
	DetailCapaian detailCapaianDTO `json:"detail_capaian"`
	Rekapitulasi  rekapitulasiDTO  `json:"rekapitulasi"`
	DetailItems   []detailItemDTO  `json:"detail_items"`
}

func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	user := userFrom(r)
	if user == nil {
		writeError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	res, err := h.svc.Get(r.Context(), user.ID, parseQueryInt(r, "tahun", 0))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	data := dashboardData{
		Pegawai:       newPegawaiBrief(res.Pegawai),
		Tahun:         res.Tahun,
		Kategori:      res.Kategori,
		JumlahMinimal: res.JumlahMinimal,
		JumlahCapaian: res.JumlahCapaian,
		DetailCapaian: detailCapaianDTO{
			Diklat:  res.DetailCapaian.Diklat,
			Ppm:     res.DetailCapaian.Ppm,
			Seminar: res.DetailCapaian.Seminar,
			Webinar: res.DetailCapaian.Webinar,
			Lc:      res.DetailCapaian.Lc,
		},
		Rekapitulasi: rekapitulasiDTO{
			Triwulan1: res.Rekapitulasi.Triwulan1,
			Triwulan2: res.Rekapitulasi.Triwulan2,
			Triwulan3: res.Rekapitulasi.Triwulan3,
			Triwulan4: res.Rekapitulasi.Triwulan4,
		},
		DetailItems: make([]detailItemDTO, 0, len(res.DetailItems)),
	}
	for _, it := range res.DetailItems {
		data.DetailItems = append(data.DetailItems, detailItemDTO{
			UUID:      it.UUID,
			Tipe:      it.Tipe,
			Tanggal:   it.Tanggal.Format(dateLayout),
			Materi:    it.Materi,
			JumlahJam: it.JumlahJam,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": data})
}
