package scheduler

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
)

type Runner interface {
	Run(ctx context.Context) error
}

type Scheduler struct {
	cron *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

func (s *Scheduler) Register(schedule string, job Runner) error {
	if _, err := s.cron.AddFunc(schedule, func() {
		_ = job.Run(context.Background())
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
