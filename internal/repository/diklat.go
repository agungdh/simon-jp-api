package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"

	"simon-jp-api/internal/models"
)

type DiklatRepository struct {
	db *bun.DB
}

func NewDiklatRepository(db *bun.DB) *DiklatRepository {
	return &DiklatRepository{db: db}
}

func (r *DiklatRepository) List(ctx context.Context, pegawaiID int64, f Filters, page, perPage int) ([]*models.Diklat, int, error) {
	var items []*models.Diklat
	base := func() *bun.SelectQuery {
		return r.db.NewSelect().Model(&items).
			Relation("Pegawai").
			Relation("JenisPelatihan").
			Where("d.pegawai_id = ?", pegawaiID).
			Where("d.deleted_at IS NULL").
			Apply(f.apply)
	}

	total, err := paginate(ctx, base, page, perPage, "d.dari_tanggal_pelaksanaan DESC", &items)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *DiklatRepository) FindByUUID(ctx context.Context, pegawaiID int64, uuid string) (*models.Diklat, error) {
	diklat := new(models.Diklat)
	err := r.db.NewSelect().Model(diklat).
		Relation("Pegawai").
		Relation("JenisPelatihan").
		Where("d.uuid = ?", uuid).
		Where("d.pegawai_id = ?", pegawaiID).
		Where("d.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return diklat, nil
}
