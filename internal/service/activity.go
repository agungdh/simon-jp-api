package service

import (
	"context"
	"errors"

	"simon-jp-api/internal/models"
	"simon-jp-api/internal/repository"
)

type ActivityService struct {
	repo *repository.ActivityRepository
}

func NewActivityService(repo *repository.ActivityRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

func (s *ActivityService) ListSeminars(ctx context.Context, pegawaiID int64, f repository.Filters, page, perPage int) ([]*models.Seminar, int, error) {
	return s.repo.ListSeminars(ctx, pegawaiID, f, page, perPage)
}

func (s *ActivityService) ListWebinars(ctx context.Context, pegawaiID int64, f repository.Filters, page, perPage int) ([]*models.Webinar, int, error) {
	return s.repo.ListWebinars(ctx, pegawaiID, f, page, perPage)
}

func (s *ActivityService) ListLcs(ctx context.Context, pegawaiID int64, f repository.Filters, page, perPage int) ([]*models.Lc, int, error) {
	return s.repo.ListLcs(ctx, pegawaiID, f, page, perPage)
}

func (s *ActivityService) ListBelajarMandiris(ctx context.Context, pegawaiID int64, f repository.Filters, page, perPage int) ([]*models.BelajarMandiri, int, error) {
	return s.repo.ListBelajarMandiris(ctx, pegawaiID, f, page, perPage)
}

func (s *ActivityService) ListMentorings(ctx context.Context, pegawaiID int64, f repository.Filters, page, perPage int) ([]*models.Mentoring, int, error) {
	return s.repo.ListMentorings(ctx, pegawaiID, f, page, perPage)
}

func (s *ActivityService) ListCoachings(ctx context.Context, pegawaiID int64, f repository.Filters, page, perPage int) ([]*models.Coaching, int, error) {
	return s.repo.ListCoachings(ctx, pegawaiID, f, page, perPage)
}

func (s *ActivityService) ListWorkshops(ctx context.Context, pegawaiID int64, f repository.Filters, page, perPage int) ([]*models.Workshop, int, error) {
	return s.repo.ListWorkshops(ctx, pegawaiID, f, page, perPage)
}

func (s *ActivityService) GetSeminar(ctx context.Context, pegawaiID int64, uuid string) (*models.Seminar, error) {
	return getActivity(ctx, pegawaiID, uuid, s.repo.FindSeminarByUUID)
}

func (s *ActivityService) GetWebinar(ctx context.Context, pegawaiID int64, uuid string) (*models.Webinar, error) {
	return getActivity(ctx, pegawaiID, uuid, s.repo.FindWebinarByUUID)
}

func (s *ActivityService) GetLc(ctx context.Context, pegawaiID int64, uuid string) (*models.Lc, error) {
	return getActivity(ctx, pegawaiID, uuid, s.repo.FindLcByUUID)
}

func (s *ActivityService) GetBelajarMandiri(ctx context.Context, pegawaiID int64, uuid string) (*models.BelajarMandiri, error) {
	return getActivity(ctx, pegawaiID, uuid, s.repo.FindBelajarMandiriByUUID)
}

func (s *ActivityService) GetMentoring(ctx context.Context, pegawaiID int64, uuid string) (*models.Mentoring, error) {
	return getActivity(ctx, pegawaiID, uuid, s.repo.FindMentoringByUUID)
}

func (s *ActivityService) GetCoaching(ctx context.Context, pegawaiID int64, uuid string) (*models.Coaching, error) {
	return getActivity(ctx, pegawaiID, uuid, s.repo.FindCoachingByUUID)
}

func (s *ActivityService) GetWorkshop(ctx context.Context, pegawaiID int64, uuid string) (*models.Workshop, error) {
	return getActivity(ctx, pegawaiID, uuid, s.repo.FindWorkshopByUUID)
}

func getActivity[T any](ctx context.Context, pegawaiID int64, uuid string, find func(context.Context, int64, string) (*T, error)) (*T, error) {
	item, err := find(ctx, pegawaiID, uuid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item, nil
}
