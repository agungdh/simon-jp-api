package service

import (
	"context"
	"errors"

	"simon-jp-api/internal/models"
	"simon-jp-api/internal/repository"
)

type MateriPpmService struct {
	repo *repository.MateriPpmRepository
}

func NewMateriPpmService(repo *repository.MateriPpmRepository) *MateriPpmService {
	return &MateriPpmService{repo: repo}
}

func (s *MateriPpmService) List(ctx context.Context, f repository.Filters, page, perPage int) ([]*models.MateriPpm, int, error) {
	return s.repo.List(ctx, f, page, perPage)
}

func (s *MateriPpmService) Get(ctx context.Context, uuid string) (*models.MateriPpm, error) {
	m, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return m, nil
}
