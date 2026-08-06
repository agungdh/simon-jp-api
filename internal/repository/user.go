package repository

import (
	"context"

	"github.com/uptrace/bun"

	"simon-jp-api/internal/models"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) scanUser(ctx context.Context, q *bun.SelectQuery) (*models.User, error) {
	user := new(models.User)
	if err := q.Model(user).Where("deleted_at IS NULL").Scan(ctx); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	return r.scanUser(ctx, r.db.NewSelect().Where("username = ?", username))
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*models.User, error) {
	return r.scanUser(ctx, r.db.NewSelect().Where("id = ?", id))
}

func (r *UserRepository) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	return r.scanUser(ctx, r.db.NewSelect().Where("uuid = ?", uuid))
}
