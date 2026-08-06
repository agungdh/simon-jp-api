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

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().
		Model(&user).
		Where("username = ?", username).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().
		Model(&user).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().
		Model(&user).
		Where("uuid = ?", uuid).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
