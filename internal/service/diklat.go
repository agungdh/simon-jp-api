package service

import (
	"context"
	"errors"

	"simon-jp-api/internal/models"
	"simon-jp-api/internal/repository"
)

type DiklatService struct {
	repo *repository.DiklatRepository
}

func NewDiklatService(repo *repository.DiklatRepository) *DiklatService {
	return &DiklatService{repo: repo}
}

func (s *DiklatService) List(ctx context.Context, pegawaiID int64, f repository.Filters, page, perPage int) ([]*models.Diklat, int, error) {
	return s.repo.List(ctx, pegawaiID, f, page, perPage)
}

func (s *DiklatService) Get(ctx context.Context, pegawaiID int64, uuid string) (*models.Diklat, error) {
	d, err := s.repo.FindByUUID(ctx, pegawaiID, uuid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return d, nil
}
