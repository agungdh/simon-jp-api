package main

import (
	"context"
	"log"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	bunDB, err := db.Connect(cfg.DBURL)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer db.Close(context.Background(), bunDB)

	if err := db.Migrate(ctx, bunDB); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	if err := db.Seed(ctx, bunDB); err != nil {
		log.Fatalf("seed: %v", err)
	}

	redisClient, err := db.ConnectRedis(ctx, cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)
	if err != nil {
		log.Fatalf("connect redis: %v", err)
	}
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(bunDB)
	sessionStore := service.NewSessionStore(redisClient, cfg.SessionTTL)
	authService := service.NewAuthService(userRepo, sessionStore)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: httpapi.NewRouter(authService),
	}

	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown: %v", err)
	}
}
