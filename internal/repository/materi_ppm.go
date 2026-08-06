package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"

	"simon-jp-api/internal/models"
)

type MateriPpmRepository struct {
	db *bun.DB
}

func NewMateriPpmRepository(db *bun.DB) *MateriPpmRepository {
	return &MateriPpmRepository{db: db}
}

func (r *MateriPpmRepository) List(ctx context.Context, f Filters, page, perPage int) ([]*models.MateriPpm, int, error) {
	var items []*models.MateriPpm
	base := func() *bun.SelectQuery {
		return r.db.NewSelect().Model(&items).
			Where("deleted_at IS NULL").
			Apply(f.apply)
	}

	total, err := paginate(ctx, base, page, perPage, "tanggal_pelaksanaan DESC", &items)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *MateriPpmRepository) FindByUUID(ctx context.Context, uuid string) (*models.MateriPpm, error) {
	item := new(models.MateriPpm)
	err := r.db.NewSelect().Model(item).
		Where("uuid = ?", uuid).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item, nil
}
