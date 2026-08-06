package service

import (
	"context"
	"errors"

	"simon-jp-api/internal/models"
	"simon-jp-api/internal/repository"
)

type PegawaiService struct {
	repo *repository.PegawaiRepository
}

func NewPegawaiService(repo *repository.PegawaiRepository) *PegawaiService {
	return &PegawaiService{repo: repo}
}

func (s *PegawaiService) ListOwn(ctx context.Context, pegawaiID int64, search, status string, bidangID int64, page, perPage int) ([]*models.Pegawai, int, error) {
	f := repository.Filters{
		Search:        search,
		SearchColumns: []string{"nama", "nip"},
		Status:        status,
		StatusColumn:  "status",
		ExactID:       bidangID,
		ExactColumn:   "bidang_id",
	}
	return s.repo.ListOwn(ctx, pegawaiID, f, page, perPage)
}

func (s *PegawaiService) GetOwn(ctx context.Context, pegawaiID int64, uuid string) (*models.Pegawai, error) {
	p, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if p.ID != pegawaiID {
		return nil, ErrForbidden
	}
	return p, nil
}
