package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"simon-jp-api/internal/mq"
)

type PingJob struct {
	client *mq.Client
}

func NewPingJob(client *mq.Client) *PingJob {
	return &PingJob{client: client}
}

func (j *PingJob) Run(ctx context.Context) error {
	data, err := json.Marshal(map[string]any{
		"at": time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	return j.client.Publish(ctx, mq.Message{
		Type: "ping",
		Data: data,
	})
}
