package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"simon-jp-api/internal/config"
	"simon-jp-api/internal/db"
	"simon-jp-api/internal/httpapi"
	"simon-jp-api/internal/repository"
	"simon-jp-api/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	bunDB, err := db.Connect(cfg.DBURL, db.Config{
		MaxOpenConns: cfg.DBMaxOpenConns,
		MaxIdleConns: cfg.DBMaxIdleConns,
		ConnMaxLife:  cfg.DBConnMaxLife,
	})
	if err != nil {
		slog.Error("connect postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close(context.Background(), bunDB)

	if err := db.Migrate(ctx, bunDB); err != nil {
		slog.Error("migrate", "error", err)
		os.Exit(1)
	}

	redisClient, err := db.ConnectRedis(ctx, cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB, db.RedisConfig{
		PoolSize: cfg.RedisPoolSize,
		Timeout:  cfg.RedisTimeout,
	})
	if err != nil {
		slog.Error("connect redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(bunDB)
	sessionStore := service.NewSessionStore(redisClient, cfg.SessionTTL)
	throttle := service.NewThrottle(redisClient, cfg.LoginMaxAttempt, cfg.LoginLockoutTTL)
	authService := service.NewAuthService(userRepo, sessionStore, throttle)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpapi.NewRouter(authService, logger),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		slog.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen and serve", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown", "error", err)
	}
}
