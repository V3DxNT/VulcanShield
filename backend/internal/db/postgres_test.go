

package db

import (
	"context"
	"os"
	"testing"

	"github.com/vulcanshield/backend/internal/config"
)

func TestNewPool_Integration(t *testing.T) {
	
	
	t.Setenv("POSTGRES_PASSWORD", getEnvOrDefault("POSTGRES_PASSWORD", "vulcanpass"))
	t.Setenv("POSTGRES_HOST", getEnvOrDefault("POSTGRES_HOST", "localhost"))
	t.Setenv("POSTGRES_PORT", getEnvOrDefault("POSTGRES_PORT", "5433"))
	t.Setenv("POSTGRES_USER", getEnvOrDefault("POSTGRES_USER", "vulcan"))
	t.Setenv("POSTGRES_DB", getEnvOrDefault("POSTGRES_DB", "vulcanshield"))

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}

	ctx := context.Background()
	pool, err := NewPool(ctx, cfg)
	if err != nil {
		t.Fatalf("NewPool failed: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		t.Errorf("Ping after NewPool failed: %v", err)
	}
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
