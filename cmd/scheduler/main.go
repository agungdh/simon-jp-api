package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"simon-jp-api/internal/config"
	"simon-jp-api/internal/mq"
	"simon-jp-api/internal/scheduler"
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

	sch := scheduler.New()

	if err := sch.Register(cfg.PingSchedule, scheduler.NewPingJob(client)); err != nil {
		log.Fatalf("register jobs: %v", err)
	}

	sch.Start()
	defer sch.Stop()

	log.Printf("scheduler started (ping on %q)", cfg.PingSchedule)

	<-ctx.Done()
}
