package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration loaded from environment variables.
// Required variables are validated; missing required vars cause Load to return an error.
type Config struct {
	// Server
	Port     string
	LogLevel string

	// PostgreSQL — required for startup
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresMaxConns int

	// Redis — non-fatal if unreachable
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// Kafka — non-fatal if unreachable
	KafkaBrokers []string

	// Downstream services (placeholders, unused in Phase 3)
	MLServiceURL string
	AIServiceURL string
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.PostgresUser, c.PostgresPassword,
		c.PostgresHost, c.PostgresPort,
		c.PostgresDB,
	)
}

// RedisAddr returns host:port for the Redis client.
func (c *Config) RedisAddr() string {
	return c.RedisHost + ":" + c.RedisPort
}

// Load reads configuration from environment variables and returns a validated Config.
// Missing required variables are returned as an error; main should log and exit.
func Load() (*Config, error) {
	cfg := &Config{
		Port:          getEnv("BACKEND_PORT", "8080"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		PostgresHost:  getEnv("POSTGRES_HOST", "postgres"),
		PostgresPort:  getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:  getEnv("POSTGRES_USER", "vulcan"),
		PostgresDB:    getEnv("POSTGRES_DB", "vulcanshield"),
		RedisHost:     getEnv("REDIS_HOST", "redis"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		MLServiceURL:  getEnv("ML_SERVICE_URL", "http://ml-service:8000"),
		AIServiceURL:  getEnv("AI_SERVICE_URL", "http://ai-service:8001"),
	}

	// Required variables
	var errs []string

	cfg.PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	if cfg.PostgresPassword == "" {
		errs = append(errs, "POSTGRES_PASSWORD is required")
	}

	// Kafka brokers: comma-separated
	brokersRaw := getEnv("KAFKA_BROKERS", "kafka:9092")
	cfg.KafkaBrokers = splitAndTrim(brokersRaw, ",")

	// PostgreSQL max connections
	maxConnsStr := getEnv("POSTGRES_MAX_CONNS", "20")
	maxConns, err := strconv.Atoi(maxConnsStr)
	if err != nil || maxConns < 1 {
		errs = append(errs, fmt.Sprintf("POSTGRES_MAX_CONNS must be a positive integer, got: %q", maxConnsStr))
	} else {
		cfg.PostgresMaxConns = maxConns
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "; "))
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
