package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

type Runner interface {
	Run(ctx context.Context) error
}

type Scheduler struct {
	cron    *cron.Cron
	timeout time.Duration
}

func New(timeout time.Duration) *Scheduler {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &Scheduler{
		cron:    cron.New(),
		timeout: timeout,
	}
}

func (s *Scheduler) Register(schedule string, job Runner) error {
	if _, err := s.cron.AddFunc(schedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
		defer cancel()
		if err := job.Run(ctx); err != nil {
			slog.Error("scheduled job failed", "error", err)
		}
	}); err != nil {
		return fmt.Errorf("register job %q: %w", schedule, err)
	}
	return nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() context.Context {
	return s.cron.Stop()
}
