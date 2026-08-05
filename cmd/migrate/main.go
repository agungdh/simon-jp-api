package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"simon-jp-api/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate [up|down]")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	sqlDB, err := sql.Open("pgx", cfg.DBURL)
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("set dialect: %v", err)
	}

	ctx := context.Background()
	dir := "internal/db/migrations"

	switch os.Args[1] {
	case "up":
		err = goose.UpContext(ctx, sqlDB, dir)
	case "down":
		err = goose.DownContext(ctx, sqlDB, dir)
	default:
		err = errors.New("unknown command, use up or down")
	}
	if err != nil {
		log.Fatalf("goose: %v", err)
	}
}
