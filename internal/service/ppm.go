package service

import (
	"context"
	"errors"

	"simon-jp-api/internal/models"
	"simon-jp-api/internal/repository"
)

type PpmService struct {
	repo *repository.PpmRepository
}

func NewPpmService(repo *repository.PpmRepository) *PpmService {
	return &PpmService{repo: repo}
}

func (s *PpmService) List(ctx context.Context, pegawaiID int64, f repository.Filters, page, perPage int) ([]*models.Ppm, int, error) {
	return s.repo.List(ctx, pegawaiID, f, page, perPage)
}

func (s *PpmService) Get(ctx context.Context, pegawaiID int64, uuid string) (*models.Ppm, error) {
	p, err := s.repo.FindByUUID(ctx, pegawaiID, uuid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return p, nil
}
