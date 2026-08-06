package httpapi

import (
	"context"
	"net/http"

	"simon-jp-api/internal/repository"
	"simon-jp-api/internal/service"
)

type ActivityHandler struct {
	svc *service.ActivityService
}

func NewActivityHandler(svc *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{svc: svc}
}

func (h *ActivityHandler) ListSeminars(w http.ResponseWriter, r *http.Request) {
	handleActivityList(w, r, h.svc.ListSeminars, newSeminarListItem)
}

func (h *ActivityHandler) ListWebinars(w http.ResponseWriter, r *http.Request) {
	handleActivityList(w, r, h.svc.ListWebinars, newWebinarListItem)
}

func (h *ActivityHandler) ListLcs(w http.ResponseWriter, r *http.Request) {
	handleActivityList(w, r, h.svc.ListLcs, newLcListItem)
}

func (h *ActivityHandler) ListBelajarMandiris(w http.ResponseWriter, r *http.Request) {
	handleActivityList(w, r, h.svc.ListBelajarMandiris, newBelajarMandiriListItem)
}

func (h *ActivityHandler) ListMentorings(w http.ResponseWriter, r *http.Request) {
	handleActivityList(w, r, h.svc.ListMentorings, newMentoringListItem)
}

func (h *ActivityHandler) ListCoachings(w http.ResponseWriter, r *http.Request) {
	handleActivityList(w, r, h.svc.ListCoachings, newCoachingListItem)
}

func (h *ActivityHandler) ListWorkshops(w http.ResponseWriter, r *http.Request) {
	handleActivityList(w, r, h.svc.ListWorkshops, newWorkshopListItem)
}

func (h *ActivityHandler) ShowSeminar(w http.ResponseWriter, r *http.Request) {
	handleActivityShow(w, r, h.svc.GetSeminar, newSeminarShow)
}

func (h *ActivityHandler) ShowWebinar(w http.ResponseWriter, r *http.Request) {
	handleActivityShow(w, r, h.svc.GetWebinar, newWebinarShow)
}

func (h *ActivityHandler) ShowLc(w http.ResponseWriter, r *http.Request) {
	handleActivityShow(w, r, h.svc.GetLc, newLcShow)
}

func (h *ActivityHandler) ShowBelajarMandiri(w http.ResponseWriter, r *http.Request) {
	handleActivityShow(w, r, h.svc.GetBelajarMandiri, newBelajarMandiriShow)
}

func (h *ActivityHandler) ShowMentoring(w http.ResponseWriter, r *http.Request) {
	handleActivityShow(w, r, h.svc.GetMentoring, newMentoringShow)
}

func (h *ActivityHandler) ShowCoaching(w http.ResponseWriter, r *http.Request) {
	handleActivityShow(w, r, h.svc.GetCoaching, newCoachingShow)
}

func (h *ActivityHandler) ShowWorkshop(w http.ResponseWriter, r *http.Request) {
	handleActivityShow(w, r, h.svc.GetWorkshop, newWorkshopShow)
}

func handleActivityList[T any](w http.ResponseWriter, r *http.Request,
	list func(context.Context, int64, repository.Filters, int, int) ([]*T, int, error),
	toDTO func(*T) activityListItem,
) {
	pegawai := pegawaiFrom(r)
	if pegawai == nil {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	page, perPage := parsePagination(r)
	f := filtersFromRequest(r,
		[]string{"materi_pengembangan"},
		"tanggal_pelaksanaan", "tanggal_pelaksanaan", "tanggal_pelaksanaan")

	items, total, err := list(r.Context(), pegawai.ID, f, page, perPage)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeActivityPage(w, r, items, total, page, perPage, toDTO)
}

func handleActivityShow[T any](w http.ResponseWriter, r *http.Request,
	get func(context.Context, int64, string) (*T, error),
	toShow func(*T) activityShow,
) {
	pegawai := pegawaiFrom(r)
	if pegawai == nil {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	item, err := get(r.Context(), pegawai.ID, routeParam(r, "uuid"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": toShow(item)})
}
