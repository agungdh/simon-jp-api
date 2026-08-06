package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"simon-jp-api/internal/config"
	"simon-jp-api/internal/mq"
	"simon-jp-api/internal/worker"
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

	registry := worker.NewRegistry()
	registry.Register(pingHandler{})

	slog.Info("worker listening", "queue", cfg.MQQueue, "consumers", cfg.MQConsumers)
	if err := client.Consume(ctx, registry.Dispatch); err != nil && ctx.Err() == nil {
		slog.Error("consume", "error", err)
		os.Exit(1)
	}
}

type pingHandler struct{}

func (pingHandler) Type() string { return "ping" }

func (pingHandler) Handle(_ context.Context, data json.RawMessage) error {
	slog.Info("received ping job", "data", string(data))
	return nil
}
