package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"simon-jp-api/internal/repository"
)

type DownloadResult struct {
	OriginalName string
	ContentType  string
	Path         string
}

type DownloadService struct {
	diklat     *repository.DiklatRepository
	activity   *repository.ActivityRepository
	storageDir string
}

func NewDownloadService(diklat *repository.DiklatRepository, activity *repository.ActivityRepository, storageDir string) *DownloadService {
	return &DownloadService{diklat: diklat, activity: activity, storageDir: storageDir}
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

	path := filepath.Join(s.storageDir, module, uuid)
	if _, err := os.Stat(path); err != nil {
		return nil, ErrFileNotFound
	}

	return &DownloadResult{
		OriginalName: *filename,
		ContentType:  mimeFor(*filename),
		Path:         path,
	}, nil
}

func mimeFor(filename string) string {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".txt":
		return "text/plain"
	case ".zip":
		return "application/zip"
	case ".rar":
		return "application/x-rar-compressed"
	default:
		return "application/octet-stream"
	}
}
