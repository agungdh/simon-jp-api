package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"simon-jp-api/internal/config"
	"simon-jp-api/internal/mq"
	"simon-jp-api/internal/scheduler"
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

	client, err := mq.Connect(ctx, cfg.MQURL, cfg.MQQueue, mq.Options{
		Prefetch:   cfg.MQPrefetch,
		Consumers:  cfg.MQConsumers,
		MaxRetries: cfg.MQMaxRetries,
	})
	if err != nil {
		slog.Error("connect mq", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	sch := scheduler.New(cfg.JobTimeout)

	if err := sch.Register(cfg.PingSchedule, scheduler.NewPingJob(client)); err != nil {
		slog.Error("register jobs", "error", err)
		os.Exit(1)
	}

	sch.Start()
	defer sch.Stop()

	slog.Info("scheduler started", "ping_schedule", cfg.PingSchedule)

	<-ctx.Done()
}
