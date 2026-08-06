package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrTooManyAttempts = errors.New("too many attempts")

type Throttler interface {
	Check(ctx context.Context, key string) error
	RecordFailure(ctx context.Context, key string) error
	Reset(ctx context.Context, key string) error
}

const throttleKeyPrefix = "throttle:"

type RedisThrottle struct {
	client      *redis.Client
	maxAttempts int
	ttl         time.Duration
}

func NewThrottle(client *redis.Client, maxAttempts int, ttl time.Duration) *RedisThrottle {
	return &RedisThrottle{client: client, maxAttempts: maxAttempts, ttl: ttl}
}

func (t *RedisThrottle) Check(ctx context.Context, key string) error {
	n, err := t.client.Get(ctx, throttleKeyPrefix+key).Int()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("check throttle: %w", err)
	}
	if n >= t.maxAttempts {
		return ErrTooManyAttempts
	}
	return nil
}

func (t *RedisThrottle) RecordFailure(ctx context.Context, key string) error {
	pipe := t.client.Pipeline()
	incr := pipe.Incr(ctx, throttleKeyPrefix+key)
	pipe.Expire(ctx, throttleKeyPrefix+key, t.ttl)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("record throttle: %w", err)
	}
	if incr.Val() >= int64(t.maxAttempts) {
		return ErrTooManyAttempts
	}
	return nil
}

func (t *RedisThrottle) Reset(ctx context.Context, key string) error {
	return t.client.Del(ctx, throttleKeyPrefix+key).Err()
}
