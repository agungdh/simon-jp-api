package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	PoolSize int
	Timeout  time.Duration
}

func ConnectRedis(ctx context.Context, addr, password string, dbIdx int, cfg RedisConfig) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbIdx,
	}
	if cfg.PoolSize > 0 {
		opts.PoolSize = cfg.PoolSize
	}
	if cfg.Timeout > 0 {
		opts.DialTimeout = cfg.Timeout
		opts.ReadTimeout = cfg.Timeout
		opts.WriteTimeout = cfg.Timeout
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
