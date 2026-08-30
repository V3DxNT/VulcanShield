package generator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/vulcanshield/backend/internal/db/repository"
	"github.com/vulcanshield/backend/internal/kafka"
	"github.com/vulcanshield/backend/internal/models"
)

// ErrScenarioRunning is returned when Start is called while a scenario is already running.
var ErrScenarioRunning = errors.New("a scenario is already running")

// Engine manages the lifecycle of scenario generation runs.
// Constraints enforced:
//   - Only one scenario may run at a time (returns ErrScenarioRunning on conflict).
//   - Stop() blocks until the generator goroutine has fully exited (guaranteed).
//   - Kafka failure is logged as a warning but does NOT suppress the transaction record.
//   - PostgreSQL failure on transaction save is a hard error that halts generation.
type Engine struct {
	mu         sync.Mutex
	active     *models.ScenarioRun
	cancelFn   context.CancelFunc
	done       chan struct{}
	log        *slog.Logger
	txRepo     repository.TransactionRepository
	scRepo     repository.ScenarioRepository
	entityRepo repository.EntityRepository
	kafka      *kafka.Producer // may be nil if Kafka unavailable
}

// NewEngine creates an Engine ready to accept Start calls.
func NewEngine(
	log *slog.Logger,
	txRepo repository.TransactionRepository,
	scRepo repository.ScenarioRepository,
	entityRepo repository.EntityRepository,
	kafkaProducer *kafka.Producer,
) *Engine {
	return &Engine{
		log:        log,
		txRepo:     txRepo,
		scRepo:     scRepo,
		entityRepo: entityRepo,
		kafka:      kafkaProducer,
	}
}

// Start begins a new scenario generation run. Returns ErrScenarioRunning (HTTP 409)
// if a scenario is already active. The run executes in a background goroutine.
func (e *Engine) Start(ctx context.Context, req models.ScenarioStartRequest) (*models.ScenarioRun, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.active != nil && e.active.Status == models.ScenarioRunning {
		return nil, ErrScenarioRunning
	}

	// Apply defaults
	if req.Transactions <= 0 {
		req.Transactions = 100
	}
	if req.IntervalMS <= 0 {
		req.IntervalMS = 1000
	}
	if req.Seed == 0 {
		req.Seed = 42
	}

	// Validate scenario type
	if !validScenario(req.Scenario) {
		req.Scenario = models.ScenarioNormal
	}

	run := &models.ScenarioRun{
		ScenarioID:       fmt.Sprintf("SC-%d", time.Now().UnixMilli()),
		ScenarioType:     req.Scenario,
		TransactionCount: req.Transactions,
		GeneratedCount:   0,
		Seed:             req.Seed,
		Status:           models.ScenarioRunning,
		StartedAt:        time.Now().UTC(),
	}

	if err := e.scRepo.Create(ctx, run); err != nil {
		return nil, fmt.Errorf("persisting scenario run: %w", err)
	}

	runCtx, cancel := context.WithCancel(context.Background())
	e.active = run
	e.cancelFn = cancel
	e.done = make(chan struct{})

	go e.runLoop(runCtx, run, req)
	return run, nil
}

// Stop halts the active scenario and blocks until the goroutine exits.
// No further transactions are produced after Stop returns.
func (e *Engine) Stop(ctx context.Context) (*models.ScenarioRun, error) {
	e.mu.Lock()
	if e.active == nil || e.active.Status != models.ScenarioRunning {
		e.mu.Unlock()
		return e.active, nil
	}
	cancel := e.cancelFn
	done := e.done
	e.mu.Unlock()

	cancel()
	<-done // wait until goroutine exits (guaranteed stop)

	e.mu.Lock()
	defer e.mu.Unlock()
	return e.active, nil
}

// GetStatus returns the current scenario run, or nil if no run has occurred.
func (e *Engine) GetStatus() *models.ScenarioRun {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.active
}

// runLoop is the background goroutine executing the scenario.
func (e *Engine) runLoop(ctx context.Context, run *models.ScenarioRun, req models.ScenarioStartRequest) {
	defer close(e.done)

	// Load entity pool
	pool, err := e.entityRepo.LoadPool(ctx)
	if err != nil {
		e.log.Error("generator: failed to load entity pool", "error", err, "scenario_id", run.ScenarioID)
		e.finalize(run, models.ScenarioStopped)
		return
	}
	if len(pool.Users) == 0 || len(pool.DeviceIDs) == 0 {
		e.log.Error("generator: entity pool is empty — run seed data", "scenario_id", run.ScenarioID)
		e.finalize(run, models.ScenarioStopped)
		return
	}

	scene := ScenarioFor(req.Scenario)
	gen := NewBaseGenerator(req.Seed, pool, scene)

	targetUserIndex := -1
	if req.CustomerID != "" {
		for i, u := range pool.Users {
			if u.UserID == req.CustomerID {
				targetUserIndex = i
				break
			}
		}
	}

	ticker := time.NewTicker(time.Duration(req.IntervalMS) * time.Millisecond)
	defer ticker.Stop()

	generated := 0
	for generated < req.Transactions {
		select {
		case <-ctx.Done():
			e.log.Info("generator: stopped by request", "scenario_id", run.ScenarioID, "generated", generated)
			e.finalizeWithCount(run, models.ScenarioStopped, generated)
			return
		case <-ticker.C:
			tx := gen.Next(generated, targetUserIndex)

			// 1. Persist to PostgreSQL (authoritative store — hard failure halts run)
			if err := e.txRepo.Create(ctx, &tx); err != nil {
				if ctx.Err() != nil {
					// Shutdown in progress — tolerate the error
					e.finalizeWithCount(run, models.ScenarioStopped, generated)
					return
				}
				e.log.Error("generator: failed to persist transaction", "error", err,
					"transaction_id", tx.TransactionID, "scenario_id", run.ScenarioID)
				e.finalizeWithCount(run, models.ScenarioStopped, generated)
				return
			}

			// 2. Audit event (non-fatal on failure)
			auditDetails := map[string]any{
				"scenario_id":   run.ScenarioID,
				"scenario_type": string(run.ScenarioType),
				"amount":        tx.Amount,
				"user_id":       tx.UserID,
			}
			if err := e.scRepo.InsertAuditEvent(ctx, tx.TransactionID, "TRANSACTION_CREATED", auditDetails); err != nil {
				e.log.Warn("generator: audit event insert failed", "error", err, "transaction_id", tx.TransactionID)
			}

			// 3. Kafka event publication (non-fatal — Kafka failure must not suppress created record)
			if e.kafka != nil {
				payload, _ := json.Marshal(tx)
				if err := e.kafka.Produce(ctx, kafka.TopicTransactionCreated, tx.TransactionID, payload); err != nil {
					e.log.Warn("generator: kafka produce failed", "error", err, "transaction_id", tx.TransactionID)
					// Do NOT halt — transaction is already persisted in PostgreSQL
				}
			}

			generated++
			e.mu.Lock()
			run.GeneratedCount = generated
			e.mu.Unlock()

			e.log.Debug("generator: transaction emitted",
				"transaction_id", tx.TransactionID,
				"amount", tx.Amount,
				"user_id", tx.UserID,
				"scenario_id", run.ScenarioID,
			)
		}
	}

	e.log.Info("generator: scenario completed", "scenario_id", run.ScenarioID, "generated", generated)
	e.finalizeWithCount(run, models.ScenarioCompleted, generated)
}

// finalize updates the scenario status in memory and PostgreSQL.
func (e *Engine) finalize(run *models.ScenarioRun, status models.ScenarioStatus) {
	e.finalizeWithCount(run, status, run.GeneratedCount)
}

func (e *Engine) finalizeWithCount(run *models.ScenarioRun, status models.ScenarioStatus, count int) {
	now := time.Now().UTC()
	e.mu.Lock()
	run.Status = status
	run.GeneratedCount = count
	run.EndedAt = &now
	e.mu.Unlock()

	// Best-effort DB update — use a fresh background context as the run ctx may be cancelled
	updateCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.scRepo.UpdateStatus(updateCtx, run.ScenarioID, status, &now); err != nil {
		e.log.Warn("generator: failed to update scenario status", "error", err, "scenario_id", run.ScenarioID)
	}
}

// validScenario checks whether a scenario type string is one of the canonical names.
func validScenario(t models.ScenarioType) bool {
	switch t {
	case models.ScenarioNormal,
		models.ScenarioVelocityAttack,
		models.ScenarioAccountTakeover,
		models.ScenarioDeviceFarm,
		models.ScenarioIPAbuse,
		models.ScenarioAmountAnomaly:
		return true
	}
	return false
}
