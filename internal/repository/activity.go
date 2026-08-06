package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"

	"simon-jp-api/internal/models"
)

type ActivityRepository struct {
	db *bun.DB
}

func NewActivityRepository(db *bun.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func listActivity[T any](db *bun.DB, ctx context.Context, table string, pegawaiID int64, f Filters, page, perPage int, order string) ([]*T, int, error) {
	var items []*T
	base := func() *bun.SelectQuery {
		return db.NewSelect().Model(&items).
			Relation("Pegawai").
			Where(table+".pegawai_id = ?", pegawaiID).
			Where(table + ".deleted_at IS NULL").
			Apply(f.apply)
	}

	total, err := paginate(ctx, base, page, perPage, order, &items)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func findActivityByUUID[T any](db *bun.DB, ctx context.Context, table string, pegawaiID int64, uuid string) (*T, error) {
	item := new(T)
	err := db.NewSelect().Model(item).
		Relation("Pegawai").
		Where(table+".uuid = ?", uuid).
		Where(table+".pegawai_id = ?", pegawaiID).
		Where(table + ".deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item, nil
}

func (r *ActivityRepository) ListSeminars(ctx context.Context, pegawaiID int64, f Filters, page, perPage int) ([]*models.Seminar, int, error) {
	return listActivity[models.Seminar](r.db, ctx, "s", pegawaiID, f, page, perPage, "tanggal_pelaksanaan DESC")
}

func (r *ActivityRepository) ListWebinars(ctx context.Context, pegawaiID int64, f Filters, page, perPage int) ([]*models.Webinar, int, error) {
	return listActivity[models.Webinar](r.db, ctx, "wb", pegawaiID, f, page, perPage, "tanggal_pelaksanaan DESC")
}

func (r *ActivityRepository) ListLcs(ctx context.Context, pegawaiID int64, f Filters, page, perPage int) ([]*models.Lc, int, error) {
	return listActivity[models.Lc](r.db, ctx, "lc", pegawaiID, f, page, perPage, "tanggal_pelaksanaan DESC")
}

func (r *ActivityRepository) ListBelajarMandiris(ctx context.Context, pegawaiID int64, f Filters, page, perPage int) ([]*models.BelajarMandiri, int, error) {
	return listActivity[models.BelajarMandiri](r.db, ctx, "bm", pegawaiID, f, page, perPage, "tanggal_pelaksanaan DESC")
}

func (r *ActivityRepository) ListMentorings(ctx context.Context, pegawaiID int64, f Filters, page, perPage int) ([]*models.Mentoring, int, error) {
	return listActivity[models.Mentoring](r.db, ctx, "mt", pegawaiID, f, page, perPage, "tanggal_pelaksanaan DESC")
}

func (r *ActivityRepository) ListCoachings(ctx context.Context, pegawaiID int64, f Filters, page, perPage int) ([]*models.Coaching, int, error) {
	return listActivity[models.Coaching](r.db, ctx, "co", pegawaiID, f, page, perPage, "tanggal_pelaksanaan DESC")
}

func (r *ActivityRepository) ListWorkshops(ctx context.Context, pegawaiID int64, f Filters, page, perPage int) ([]*models.Workshop, int, error) {
	return listActivity[models.Workshop](r.db, ctx, "ws", pegawaiID, f, page, perPage, "tanggal_pelaksanaan DESC")
}

func (r *ActivityRepository) FindSeminarByUUID(ctx context.Context, pegawaiID int64, uuid string) (*models.Seminar, error) {
	return findActivityByUUID[models.Seminar](r.db, ctx, "s", pegawaiID, uuid)
}

func (r *ActivityRepository) FindWebinarByUUID(ctx context.Context, pegawaiID int64, uuid string) (*models.Webinar, error) {
	return findActivityByUUID[models.Webinar](r.db, ctx, "wb", pegawaiID, uuid)
}

func (r *ActivityRepository) FindLcByUUID(ctx context.Context, pegawaiID int64, uuid string) (*models.Lc, error) {
	return findActivityByUUID[models.Lc](r.db, ctx, "lc", pegawaiID, uuid)
}

func (r *ActivityRepository) FindBelajarMandiriByUUID(ctx context.Context, pegawaiID int64, uuid string) (*models.BelajarMandiri, error) {
	return findActivityByUUID[models.BelajarMandiri](r.db, ctx, "bm", pegawaiID, uuid)
}

func (r *ActivityRepository) FindMentoringByUUID(ctx context.Context, pegawaiID int64, uuid string) (*models.Mentoring, error) {
	return findActivityByUUID[models.Mentoring](r.db, ctx, "mt", pegawaiID, uuid)
}

func (r *ActivityRepository) FindCoachingByUUID(ctx context.Context, pegawaiID int64, uuid string) (*models.Coaching, error) {
	return findActivityByUUID[models.Coaching](r.db, ctx, "co", pegawaiID, uuid)
}

func (r *ActivityRepository) FindWorkshopByUUID(ctx context.Context, pegawaiID int64, uuid string) (*models.Workshop, error) {
	return findActivityByUUID[models.Workshop](r.db, ctx, "ws", pegawaiID, uuid)
}
