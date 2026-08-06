package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const sessionKeyPrefix = "session:"

type SessionStore interface {
	Create(ctx context.Context, token string, userID int64) error
	Get(ctx context.Context, token string) (int64, error)
	Delete(ctx context.Context, token string) error
	TTL() time.Duration
}

type RedisSessionStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewSessionStore(client *redis.Client, ttl time.Duration) *RedisSessionStore {
	return &RedisSessionStore{client: client, ttl: ttl}
}

func (s *RedisSessionStore) Create(ctx context.Context, token string, userID int64) error {
	key := sessionKeyPrefix + token
	return s.client.Set(ctx, key, userID, s.ttl).Err()
}

func (s *RedisSessionStore) Get(ctx context.Context, token string) (int64, error) {
	key := sessionKeyPrefix + token
	val, err := s.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, ErrSessionNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("get session: %w", err)
	}
	return val, nil
}

func (s *RedisSessionStore) Delete(ctx context.Context, token string) error {
	key := sessionKeyPrefix + token
	return s.client.Del(ctx, key).Err()
}

func (s *RedisSessionStore) TTL() time.Duration {
	return s.ttl
}
