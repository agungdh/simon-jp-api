package service

import (
	"context"
	"errors"

	"simon-jp-api/internal/repository"
	"simon-jp-api/internal/storage"
)

type DownloadResult struct {
	URL string
}

type DownloadService struct {
	diklat    *repository.DiklatRepository
	activity  *repository.ActivityRepository
	presigner storage.Presigner
}

func NewDownloadService(diklat *repository.DiklatRepository, activity *repository.ActivityRepository, presigner storage.Presigner) *DownloadService {
	return &DownloadService{diklat: diklat, activity: activity, presigner: presigner}
}

func (s *DownloadService) Resolve(ctx context.Context, module, uuid string, pegawaiID int64) (*DownloadResult, error) {
	var filename *string

	switch module {
	case "diklat":
		rec, err := s.diklat.FindByUUID(ctx, pegawaiID, uuid)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		filename = rec.Filename
	case "seminar":
		rec, err := s.activity.FindSeminarByUUID(ctx, pegawaiID, uuid)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		filename = rec.Filename
	case "webinar":
		rec, err := s.activity.FindWebinarByUUID(ctx, pegawaiID, uuid)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		filename = rec.Filename
	case "lc":
		rec, err := s.activity.FindLcByUUID(ctx, pegawaiID, uuid)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		filename = rec.Filename
	case "belajar-mandiri":
		rec, err := s.activity.FindBelajarMandiriByUUID(ctx, pegawaiID, uuid)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		filename = rec.Filename
	case "mentoring":
		rec, err := s.activity.FindMentoringByUUID(ctx, pegawaiID, uuid)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		filename = rec.Filename
	case "coaching":
		rec, err := s.activity.FindCoachingByUUID(ctx, pegawaiID, uuid)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		filename = rec.Filename
	case "workshop":
		rec, err := s.activity.FindWorkshopByUUID(ctx, pegawaiID, uuid)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		filename = rec.Filename
	default:
		return nil, ErrModuleNotFound
	}

	if filename == nil || *filename == "" {
		return nil, ErrNoFileAttached
	}

	url, err := s.presigner.PresignedGetURL(ctx, module+"/"+uuid, *filename)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return &DownloadResult{URL: url}, nil
}
