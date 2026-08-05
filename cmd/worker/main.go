package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"simon-jp-api/internal/config"
	"simon-jp-api/internal/mq"
	"simon-jp-api/internal/worker"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	client, err := mq.Connect(ctx, cfg.MQURL, cfg.MQQueue)
	if err != nil {
		log.Fatalf("connect mq: %v", err)
	}
	defer client.Close()

	registry := worker.NewRegistry()
	registry.Register(pingHandler{})

	log.Printf("worker listening on queue %q", cfg.MQQueue)
	if err := client.Consume(ctx, registry.Dispatch); err != nil {
		log.Fatalf("consume: %v", err)
	}
}

type pingHandler struct{}

func (pingHandler) Type() string { return "ping" }

func (pingHandler) Handle(_ context.Context, data json.RawMessage) error {
	log.Printf("worker: received ping job: %s", data)
	return nil
}
