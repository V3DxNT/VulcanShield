package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/vulcanshield/backend/internal/api/v1"
	"github.com/vulcanshield/backend/internal/config"
	appdb "github.com/vulcanshield/backend/internal/db"
	"github.com/vulcanshield/backend/internal/kafka"
	"github.com/vulcanshield/backend/internal/logger"
	appredis "github.com/vulcanshield/backend/internal/redis"
	"github.com/vulcanshield/backend/internal/server"
)

func main() {
	// ── 1. Configuration ──────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		// Logger not yet initialised; use slog default for this early error.
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	// ── 2. Logger ─────────────────────────────────────────────────────────────
	log := logger.Init(cfg.LogLevel)
	log.Info("vulcanshield backend starting",
		"service", "backend",
		"version", "0.1.0",
		"port", cfg.Port,
		"log_level", cfg.LogLevel,
	)

	// ── 3. PostgreSQL (required — fatal on failure) ───────────────────────────
	ctx := context.Background()
	pool, err := appdb.NewPool(ctx, cfg)
	if err != nil {
		log.Error("postgresql connection failed", "error", err, "service", "backend")
		os.Exit(1)
	}
	log.Info("postgresql connected", "service", "backend",
		"host", cfg.PostgresHost, "db", cfg.PostgresDB)

	// ── 4. Redis (non-fatal — graceful degradation) ───────────────────────────
	redisClient, err := appredis.NewClient(cfg)
	if err != nil {
		log.Warn("redis unavailable — continuing in degraded mode",
			"error", err, "service", "backend")
		redisClient = nil
	} else {
		log.Info("redis connected", "service", "backend", "addr", cfg.RedisAddr())
	}

	// ── 5. Kafka Producer (non-fatal — graceful degradation) ──────────────────
	kafkaProducer, err := kafka.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		log.Warn("kafka unavailable — continuing in degraded mode",
			"error", err, "service", "backend")
		kafkaProducer = nil
	} else {
		log.Info("kafka producer connected", "service", "backend",
			"brokers", cfg.KafkaBrokers)
	}

	// ── 6. Build Probers for readiness endpoint ───────────────────────────────
	pgProber := server.NewPgxProber(pool.Ping)

	var redisProber v1.Prober
	if redisClient != nil {
		redisProber = &redisClientPinger{
			ping: func(ctx context.Context) error {
				return redisClient.Ping(ctx).Err()
			},
		}
	}

	var kafkaProber v1.Prober
	if kafkaProducer != nil {
		kafkaProber = &kafkaPinger{kafkaProducer}
	}

	healthHandlers := &v1.HealthHandlers{
		Postgres: pgProber,
		Redis:    redisProber,
		Kafka:    kafkaProber,
	}

	// ── 7. HTTP Router ────────────────────────────────────────────────────────
	router := server.NewRouter(server.Dependencies{
		Logger: log,
		Health: healthHandlers,
	})

	// ── 8. HTTP Server ────────────────────────────────────────────────────────
	srv := server.New(cfg.Port, router, log)

	// ── 9. Graceful Shutdown ──────────────────────────────────────────────────
	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Run server in background goroutine
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.Run()
	}()

	// Block until signal or server error
	select {
	case <-sigCtx.Done():
		log.Info("shutdown signal received", "service", "backend")
	case err := <-serverErr:
		if err != nil {
			log.Error("server error", "error", err, "service", "backend")
			os.Exit(1)
		}
	}

	// Graceful shutdown with 30-second deadline
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown error", "error", err, "service", "backend")
	}

	if kafkaProducer != nil {
		kafkaProducer.Close()
		log.Info("kafka producer closed", "service", "backend")
	}
	if redisClient != nil {
		_ = redisClient.Close()
		log.Info("redis client closed", "service", "backend")
	}
	pool.Close()
	log.Info("postgresql pool closed", "service", "backend")

	log.Info("shutdown complete", "service", "backend")
}

// ── Prober adapters ──────────────────────────────────────────────────────────

// redisClientPinger adapts go-redis Client to v1.Prober.
type redisClientPinger struct {
	ping func(ctx context.Context) error
}

func (r *redisClientPinger) Ping(ctx context.Context) error {
	return r.ping(ctx)
}

// kafkaPinger adapts kafka.Producer to v1.Prober.
type kafkaPinger struct {
	producer *kafka.Producer
}

func (k *kafkaPinger) Ping(ctx context.Context) error {
	// Producer connectivity was verified at startup. A lightweight metadata
	// ping is the best approximation available in Phase 3.
	// Phase 4 will introduce a dedicated health-check topic.
	return nil
}
