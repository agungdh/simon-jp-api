package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DBURL      string        `env:"DB_URL"`
	RedisAddr  string        `env:"REDIS_ADDR"`
	RedisPass  string        `env:"REDIS_PASSWORD"`
	RedisDB    int           `env:"REDIS_DB" envDefault:"0"`
	Port       string        `env:"PORT" envDefault:"8080"`
	SessionTTL time.Duration `env:"SESSION_TTL" envDefault:"24h"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DBURL:     "postgres://admin:admin@localhost:5432/simonjp?sslmode=disable",
		RedisAddr: "localhost:6379",
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	return cfg, nil
}
