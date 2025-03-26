package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host string `env:"REDIS_HOST" env-default:"localhost"`
	Port uint16 `env:"REDIS_PORT" env-default:"6379"`
}

func New(cfg Config, ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: "",
		DB:       0,
	})
	pong := client.Ping(ctx)
	if pong.Err() != nil {
		return nil, pong.Err()
	}
	return client, nil
}
