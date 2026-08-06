package httpapi

import (
	"net/http"
	"strconv"

	"simon-jp-api/internal/models"
	"simon-jp-api/internal/repository"
)

type pageResponse struct {
	Data     any `json:"data"`
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	Total    int `json:"total"`
	LastPage int `json:"last_page"`
	From     int `json:"from"`
	To       int `json:"to"`
}

func newPageResponse(data any, page, perPage, total int) pageResponse {
	lastPage := 0
	if perPage > 0 {
		lastPage = (total + perPage - 1) / perPage
	}
	from, to := 0, 0
	if total > 0 {
		from = (page-1)*perPage + 1
		to = page * perPage
		if to > total {
			to = total
		}
	}
	return pageResponse{
		Data:     data,
		Page:     page,
		PerPage:  perPage,
		Total:    total,
		LastPage: lastPage,
		From:     from,
		To:       to,
	}
}

func parsePagination(r *http.Request) (page, perPage int) {
	page = parseQueryInt(r, "page", 1)
	if page < 1 {
		page = 1
	}
	perPage = parseQueryInt(r, "per_page", 15)
	if perPage < 1 {
		perPage = 15
	}
	return page, perPage
}

func parseQueryInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func filtersFromRequest(r *http.Request, searchCols []string, dariCol, sampaiCol, tahunCol string) repository.Filters {
	q := r.URL.Query()
	f := repository.Filters{
		Search:        q.Get("search"),
		SearchColumns: searchCols,
		DariTanggal:   q.Get("dari_tanggal"),
		SampaiTanggal: q.Get("sampai_tanggal"),
		Tahun:         parseQueryInt(r, "tahun", 0),
		DariColumn:    dariCol,
		SampaiColumn:  sampaiCol,
		TahunColumn:   tahunCol,
	}
	return f
}

func writeActivityPage[T any](w http.ResponseWriter, r *http.Request, items []*T, total, page, perPage int, toDTO func(*T) activityListItem) {
	data := make([]activityListItem, 0, len(items))
	for _, it := range items {
		data = append(data, toDTO(it))
	}
	writeJSON(w, http.StatusOK, newPageResponse(data, page, perPage, total))
}

func newSeminarListItem(s *models.Seminar) activityListItem {
	return activityListItem{UUID: s.UUID, Pegawai: newPegawaiBrief(s.Pegawai), MateriPengembangan: s.MateriPengembangan, TanggalPelaksanaan: s.TanggalPelaksanaan.Format(dateLayout), JumlahJam: s.JumlahJam}
}

func newSeminarShow(s *models.Seminar) activityShow {
	return activityShow{activityListItem: newSeminarListItem(s), Filename: s.Filename, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt}
}

func newWebinarListItem(w *models.Webinar) activityListItem {
	return activityListItem{UUID: w.UUID, Pegawai: newPegawaiBrief(w.Pegawai), MateriPengembangan: w.MateriPengembangan, TanggalPelaksanaan: w.TanggalPelaksanaan.Format(dateLayout), JumlahJam: w.JumlahJam}
}

func newWebinarShow(w *models.Webinar) activityShow {
	return activityShow{activityListItem: newWebinarListItem(w), Filename: w.Filename, CreatedAt: w.CreatedAt, UpdatedAt: w.UpdatedAt}
}

func newLcListItem(l *models.Lc) activityListItem {
	return activityListItem{UUID: l.UUID, Pegawai: newPegawaiBrief(l.Pegawai), MateriPengembangan: l.MateriPengembangan, TanggalPelaksanaan: l.TanggalPelaksanaan.Format(dateLayout), JumlahJam: l.JumlahJam}
}

func newLcShow(l *models.Lc) activityShow {
	return activityShow{activityListItem: newLcListItem(l), Filename: l.Filename, CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt}
}

func newBelajarMandiriListItem(b *models.BelajarMandiri) activityListItem {
	return activityListItem{UUID: b.UUID, Pegawai: newPegawaiBrief(b.Pegawai), MateriPengembangan: b.MateriPengembangan, TanggalPelaksanaan: b.TanggalPelaksanaan.Format(dateLayout), JumlahJam: b.JumlahJam}
}

func newBelajarMandiriShow(b *models.BelajarMandiri) activityShow {
	return activityShow{activityListItem: newBelajarMandiriListItem(b), Filename: b.Filename, CreatedAt: b.CreatedAt, UpdatedAt: b.UpdatedAt}
}

func newMentoringListItem(m *models.Mentoring) activityListItem {
	return activityListItem{UUID: m.UUID, Pegawai: newPegawaiBrief(m.Pegawai), MateriPengembangan: m.MateriPengembangan, TanggalPelaksanaan: m.TanggalPelaksanaan.Format(dateLayout), JumlahJam: m.JumlahJam}
}

func newMentoringShow(m *models.Mentoring) activityShow {
	return activityShow{activityListItem: newMentoringListItem(m), Filename: m.Filename, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func newCoachingListItem(c *models.Coaching) activityListItem {
	return activityListItem{UUID: c.UUID, Pegawai: newPegawaiBrief(c.Pegawai), MateriPengembangan: c.MateriPengembangan, TanggalPelaksanaan: c.TanggalPelaksanaan.Format(dateLayout), JumlahJam: c.JumlahJam}
}

func newCoachingShow(c *models.Coaching) activityShow {
	return activityShow{activityListItem: newCoachingListItem(c), Filename: c.Filename, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
}

func newWorkshopListItem(w *models.Workshop) activityListItem {
	return activityListItem{UUID: w.UUID, Pegawai: newPegawaiBrief(w.Pegawai), MateriPengembangan: w.MateriPengembangan, TanggalPelaksanaan: w.TanggalPelaksanaan.Format(dateLayout), JumlahJam: w.JumlahJam}
}

func newWorkshopShow(w *models.Workshop) activityShow {
	return activityShow{activityListItem: newWorkshopListItem(w), Filename: w.Filename, CreatedAt: w.CreatedAt, UpdatedAt: w.UpdatedAt}
}
