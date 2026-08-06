package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"

	"simon-jp-api/internal/models"
)

type PegawaiRepository struct {
	db *bun.DB
}

func NewPegawaiRepository(db *bun.DB) *PegawaiRepository {
	return &PegawaiRepository{db: db}
}

func (r *PegawaiRepository) FindByUserID(ctx context.Context, userID int64) (*models.Pegawai, error) {
	pegawai := new(models.Pegawai)
	err := r.db.NewSelect().Model(pegawai).
		Relation("Bidang").
		Relation("PangkatGolongan").
		Where("p.user_id = ?", userID).
		Where("p.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return pegawai, nil
}

func (r *PegawaiRepository) FindByUUID(ctx context.Context, uuid string) (*models.Pegawai, error) {
	pegawai := new(models.Pegawai)
	err := r.db.NewSelect().Model(pegawai).
		Relation("Bidang").
		Relation("PangkatGolongan").
		Where("p.uuid = ?", uuid).
		Where("p.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return pegawai, nil
}

func (r *PegawaiRepository) ListOwn(ctx context.Context, pegawaiID int64, f Filters, page, perPage int) ([]*models.Pegawai, int, error) {
	var items []*models.Pegawai
	base := func() *bun.SelectQuery {
		return r.db.NewSelect().Model(&items).
			Relation("Bidang").
			Relation("PangkatGolongan").
			Where("p.id = ?", pegawaiID).
			Where("p.deleted_at IS NULL").
			Apply(f.apply)
	}

	total, err := paginate(ctx, base, page, perPage, "p.nama ASC", &items)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
