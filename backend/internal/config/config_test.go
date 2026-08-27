package config

import (
	"os"
	"testing"
)

func TestLoad_MissingPassword(t *testing.T) {
	os.Unsetenv("POSTGRES_PASSWORD")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error when POSTGRES_PASSWORD is missing, got nil")
	}
}

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("POSTGRES_PASSWORD", "secret")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}
	if cfg.PostgresMaxConns != 20 {
		t.Errorf("expected default max conns 20, got %d", cfg.PostgresMaxConns)
	}
	if len(cfg.KafkaBrokers) == 0 {
		t.Error("expected at least one kafka broker")
	}
}

func TestLoad_CustomPort(t *testing.T) {
	t.Setenv("POSTGRES_PASSWORD", "secret")
	t.Setenv("BACKEND_PORT", "9090")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Port)
	}
}

func TestLoad_InvalidMaxConns(t *testing.T) {
	t.Setenv("POSTGRES_PASSWORD", "secret")
	t.Setenv("POSTGRES_MAX_CONNS", "notanumber")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid POSTGRES_MAX_CONNS, got nil")
	}
}

func TestConfig_DSN(t *testing.T) {
	t.Setenv("POSTGRES_PASSWORD", "testpass")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dsn := cfg.DSN()
	if dsn == "" {
		t.Error("DSN should not be empty")
	}
}
