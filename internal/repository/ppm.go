package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"

	"simon-jp-api/internal/models"
)

type PpmRepository struct {
	db *bun.DB
}

func NewPpmRepository(db *bun.DB) *PpmRepository {
	return &PpmRepository{db: db}
}

func (r *PpmRepository) List(ctx context.Context, pegawaiID int64, f Filters, page, perPage int) ([]*models.Ppm, int, error) {
	var items []*models.Ppm
	base := func() *bun.SelectQuery {
		return r.db.NewSelect().Model(&items).
			Relation("Pegawai").
			Where("ppm.pegawai_id = ?", pegawaiID).
			Where("ppm.deleted_at IS NULL").
			Apply(f.apply)
	}

	total, err := paginate(ctx, base, page, perPage, "ppm.tanggal_pelaksanaan DESC", &items)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *PpmRepository) FindByUUID(ctx context.Context, pegawaiID int64, uuid string) (*models.Ppm, error) {
	ppm := new(models.Ppm)
	err := r.db.NewSelect().Model(ppm).
		Relation("Pegawai").
		Where("ppm.uuid = ?", uuid).
		Where("ppm.pegawai_id = ?", pegawaiID).
		Where("ppm.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return ppm, nil
}
