package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/vulcanshield/backend/internal/config"
)

// NewClient creates a go-redis/v9 client and pings Redis to verify connectivity.
// Redis connectivity is non-fatal: errors are returned to the caller which should
// log a warning and continue (graceful degradation per PROJECT_SPEC.md §100).
func NewClient(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return client, nil
}
