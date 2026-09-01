

package redis

import (
	"testing"

	"github.com/vulcanshield/backend/internal/config"
)

func TestNewClient_Integration(t *testing.T) {
	
	
	t.Setenv("POSTGRES_PASSWORD", "vulcanpass")
	t.Setenv("REDIS_HOST", "localhost")
	t.Setenv("REDIS_PORT", "6379")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()
}
