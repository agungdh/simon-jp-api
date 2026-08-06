package db

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"

	"simon-jp-api/internal/models"
)

func Seed(ctx context.Context, bunDB *bun.DB) error {
	exists, err := bunDB.NewSelect().
		Model((*models.User)(nil)).
		Where("username = ?", "admin").
		Exists(ctx)
	if err != nil {
		return fmt.Errorf("check admin user: %w", err)
	}
	if exists {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	_, err = bunDB.NewInsert().
		Model(&models.User{
			Username: "admin",
			Password: string(hash),
		}).
		Column("username", "password").
		On("CONFLICT (username) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("seed admin user: %w", err)
	}

	return nil
}
