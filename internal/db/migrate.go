package db

import (
	"context"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/uptrace/bun"

	_ "simon-jp-api/internal/db/migrations"
)

//go:embed migrations/*
var migrationsFS embed.FS

func Migrate(ctx context.Context, bunDB *bun.DB) error {
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.UpContext(ctx, bunDB.DB, "migrations"); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
