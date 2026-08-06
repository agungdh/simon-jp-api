package repository

import (
	"context"

	"github.com/uptrace/bun"

	"simon-jp-api/internal/models"
)

type CapaianRepository struct {
	db *bun.DB
}

func NewCapaianRepository(db *bun.DB) *CapaianRepository {
	return &CapaianRepository{db: db}
}

func (r *CapaianRepository) SumByYear(ctx context.Context, pegawaiID int64, year int, table, jamColumn, dateColumn string) (int, error) {
	var sum int
	err := r.db.NewRaw(`
		SELECT COALESCE(SUM(`+jamColumn+`), 0)
		FROM `+table+`
		WHERE pegawai_id = ? AND EXTRACT(YEAR FROM `+dateColumn+`) = ? AND deleted_at IS NULL
	`, pegawaiID, year).Scan(ctx, &sum)
	return sum, err
}

func (r *CapaianRepository) SumByQuarter(ctx context.Context, pegawaiID int64, year, quarter int, table, jamColumn, dateColumn string) (int, error) {
	start := (quarter-1)*3 + 1
	end := quarter * 3
	var sum int
	err := r.db.NewRaw(`
		SELECT COALESCE(SUM(`+jamColumn+`), 0)
		FROM `+table+`
		WHERE pegawai_id = ?
		  AND EXTRACT(YEAR FROM `+dateColumn+`) = ?
		  AND EXTRACT(MONTH FROM `+dateColumn+`) BETWEEN ? AND ?
		  AND deleted_at IS NULL
	`, pegawaiID, year, start, end).Scan(ctx, &sum)
	return sum, err
}

func listByYear[T any](db *bun.DB, ctx context.Context, pegawaiID int64, year int, dateColumn string) ([]*T, error) {
	var items []*T
	err := db.NewSelect().Model(&items).
		Where("pegawai_id = ?", pegawaiID).
		Where("EXTRACT(YEAR FROM "+dateColumn+") = ?", year).
		Where("deleted_at IS NULL").
		Scan(ctx)
	return items, err
}

func (r *CapaianRepository) ListDiklatsByYear(ctx context.Context, pegawaiID int64, year int) ([]*models.Diklat, error) {
	return listByYear[models.Diklat](r.db, ctx, pegawaiID, year, "dari_tanggal_pelaksanaan")
}

func (r *CapaianRepository) ListPpmsByYear(ctx context.Context, pegawaiID int64, year int) ([]*models.Ppm, error) {
	return listByYear[models.Ppm](r.db, ctx, pegawaiID, year, "tanggal_pelaksanaan")
}

func (r *CapaianRepository) ListSeminarsByYear(ctx context.Context, pegawaiID int64, year int) ([]*models.Seminar, error) {
	return listByYear[models.Seminar](r.db, ctx, pegawaiID, year, "tanggal_pelaksanaan")
}

func (r *CapaianRepository) ListWebinarsByYear(ctx context.Context, pegawaiID int64, year int) ([]*models.Webinar, error) {
	return listByYear[models.Webinar](r.db, ctx, pegawaiID, year, "tanggal_pelaksanaan")
}

func (r *CapaianRepository) ListLcsByYear(ctx context.Context, pegawaiID int64, year int) ([]*models.Lc, error) {
	return listByYear[models.Lc](r.db, ctx, pegawaiID, year, "tanggal_pelaksanaan")
}
