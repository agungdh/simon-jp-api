package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DBURL           string        `env:"DB_URL"`
	DBMaxOpenConns  int           `env:"DB_MAX_OPEN_CONNS" envDefault:"20"`
	DBMaxIdleConns  int           `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
	DBConnMaxLife   time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"30m"`
	RedisAddr       string        `env:"REDIS_ADDR"`
	RedisPass       string        `env:"REDIS_PASSWORD"`
	RedisDB         int           `env:"REDIS_DB" envDefault:"0"`
	RedisPoolSize   int           `env:"REDIS_POOL_SIZE" envDefault:"10"`
	RedisTimeout    time.Duration `env:"REDIS_TIMEOUT" envDefault:"5s"`
	Port            string        `env:"PORT" envDefault:"8080"`
	SessionTTL      time.Duration `env:"SESSION_TTL" envDefault:"24h"`
	LoginMaxAttempt int           `env:"LOGIN_MAX_ATTEMPTS" envDefault:"5"`
	LoginLockoutTTL time.Duration `env:"LOGIN_LOCKOUT_TTL" envDefault:"15m"`
	MQURL           string        `env:"MQ_URL"`
	MQQueue         string        `env:"MQ_QUEUE" envDefault:"jobs"`
	MQPrefetch      int           `env:"MQ_PREFETCH" envDefault:"10"`
	MQConsumers     int           `env:"MQ_CONSUMERS" envDefault:"5"`
	MQMaxRetries    int           `env:"MQ_MAX_RETRIES" envDefault:"5"`
	PingSchedule    string        `env:"PING_SCHEDULE" envDefault:"@every 1m"`
	JobTimeout      time.Duration `env:"JOB_TIMEOUT" envDefault:"30s"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DBURL:     "postgres://admin:admin@localhost:5432/simonjp?sslmode=disable",
		RedisAddr: "localhost:6379",
		MQURL:     "amqp://guest:guest@localhost:5672/",
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	return cfg, nil
}
