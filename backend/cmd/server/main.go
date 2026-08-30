package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vulcanshield/backend/internal/aiclient"
	v1 "github.com/vulcanshield/backend/internal/api/v1"
	"github.com/vulcanshield/backend/internal/challenge"
	"github.com/vulcanshield/backend/internal/config"
	appdb "github.com/vulcanshield/backend/internal/db"
	"github.com/vulcanshield/backend/internal/db/repository"
	"github.com/vulcanshield/backend/internal/generator"
	"github.com/vulcanshield/backend/internal/graph"
	"github.com/vulcanshield/backend/internal/kafka"
	"github.com/vulcanshield/backend/internal/logger"
	"github.com/vulcanshield/backend/internal/mlclient"
	appredis "github.com/vulcanshield/backend/internal/redis"
	"github.com/vulcanshield/backend/internal/server"
	appws "github.com/vulcanshield/backend/internal/websocket"
)

func main() {
	// ── 1. Configuration ──────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	// ── 2. Logger ─────────────────────────────────────────────────────────────
	log := logger.Init(cfg.LogLevel)
	log.Info("vulcanshield backend starting",
		"service", "backend",
		"version", "0.2.0",
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

	// ── 6. Repositories & Engine Services ────────────────────────────────────
	txRepo := repository.NewTransactionRepository(pool)
	scRepo := repository.NewScenarioRepository(pool)
	entityRepo := repository.NewEntityRepository(pool)
	riskRepo := repository.NewRiskRepository(pool)
	policyRepo := repository.NewPolicyRepository(pool)
	challengeRepo := repository.NewChallengeRepository(pool)
	graphRepo := repository.NewGraphRepository(pool)
	userRepo := repository.NewUserRepository(pool)

	var velocityEngine *appredis.VelocityEngine
	if redisClient != nil {
		velocityEngine = appredis.NewVelocityEngine(redisClient)
	}
	mlClient := mlclient.NewClient(cfg.MLServiceURL)
	otpService := challenge.NewService(redisClient)
	graphEngine := graph.NewEngine(graphRepo)
	wsHub := appws.NewHub(log)

	// ── 7. Generator Engine ──────────────────────────────────────────────────
	engine := generator.NewEngine(
		log, txRepo, scRepo, entityRepo, riskRepo,
		policyRepo, challengeRepo, userRepo,
		kafkaProducer, velocityEngine, mlClient, otpService, wsHub,
	)

	// ── 8. Build Probers for readiness endpoint ──────────────────────────────
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

	// ── 9. HTTP Router ───────────────────────────────────────────────────────
	handlers := &v1.Handlers{
		Scenarios: &v1.ScenarioHandlers{
			Engine: engine,
		},
		Transactions: &v1.TransactionHandlers{
			TxRepo: txRepo,
		},
		Challenges: &v1.ChallengeHandlers{
			ChallengeRepo: challengeRepo,
			TxRepo:        txRepo,
			PolicyRepo:    policyRepo,
			OTPService:    otpService,
		},
		Graph: &v1.GraphHandlers{
			GraphRepo: graphRepo,
			Engine:    graphEngine,
		},
		Investigations: &v1.InvestigationHandlers{
			TxRepo:   txRepo,
			RiskRepo: riskRepo,
			AIClient: aiclient.NewClient(cfg.AIServiceURL),
		},
		WSHub: wsHub,
	}

	router := server.NewRouter(server.Dependencies{
		Logger:   log,
		Health:   healthHandlers,
		Handlers: handlers,
	})

	// ── 10. HTTP Server ──────────────────────────────────────────────────────
	srv := server.New(cfg.Port, router, log)

	// ── 11. Graceful Shutdown ────────────────────────────────────────────────
	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.Run()
	}()

	select {
	case <-sigCtx.Done():
		log.Info("shutdown signal received", "service", "backend")
	case err := <-serverErr:
		if err != nil {
			log.Error("server error", "error", err, "service", "backend")
			os.Exit(1)
		}
	}

	// Stop any running scenario before server shutdown
	if _, err := engine.Stop(context.Background()); err != nil {
		log.Warn("failed to stop active scenario during shutdown", "error", err)
	}

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

type redisClientPinger struct {
	ping func(ctx context.Context) error
}

func (r *redisClientPinger) Ping(ctx context.Context) error {
	return r.ping(ctx)
}

type kafkaPinger struct {
	producer *kafka.Producer
}

func (k *kafkaPinger) Ping(ctx context.Context) error {
	return nil
}
