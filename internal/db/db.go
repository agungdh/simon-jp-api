package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Config struct {
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLife  time.Duration
}

func Connect(dsn string, cfg Config) (*bun.DB, error) {
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLife > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLife)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return bun.NewDB(sqlDB, pgdialect.New()), nil
}

func Close(ctx context.Context, bunDB *bun.DB) error {
	return bunDB.Close()
}
