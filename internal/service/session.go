package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const sessionKeyPrefix = "session:"

type SessionStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewSessionStore(client *redis.Client, ttl time.Duration) *SessionStore {
	return &SessionStore{client: client, ttl: ttl}
}

func (s *SessionStore) Create(ctx context.Context, token string, userID int64) error {
	key := sessionKeyPrefix + token
	return s.client.Set(ctx, key, userID, s.ttl).Err()
}

func (s *SessionStore) Get(ctx context.Context, token string) (int64, error) {
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

func (s *SessionStore) Delete(ctx context.Context, token string) error {
	key := sessionKeyPrefix + token
	return s.client.Del(ctx, key).Err()
}

func (s *SessionStore) TTL() time.Duration {
	return s.ttl
}
