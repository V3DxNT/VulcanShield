package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vulcanshield/backend/internal/models"
)



type ScenarioRepository interface {
	
	Create(ctx context.Context, run *models.ScenarioRun) error

	
	UpdateStatus(ctx context.Context, scenarioID string, status models.ScenarioStatus, endedAt *time.Time) error

	
	GetLatest(ctx context.Context) (*models.ScenarioRun, error)

	
	
	InsertAuditEvent(ctx context.Context, transactionID, eventType string, details any) error
}

type pgxScenarioRepository struct {
	pool *pgxpool.Pool
}


func NewScenarioRepository(pool *pgxpool.Pool) ScenarioRepository {
	return &pgxScenarioRepository{pool: pool}
}

func (r *pgxScenarioRepository) Create(ctx context.Context, run *models.ScenarioRun) error {
	const q = `
		INSERT INTO scenarios (scenario_id, scenario_type, transaction_count, seed, status, started_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, q,
		run.ScenarioID,
		string(run.ScenarioType),
		run.TransactionCount,
		run.Seed,
		string(run.Status),
		run.StartedAt,
	)
	return err
}

func (r *pgxScenarioRepository) UpdateStatus(ctx context.Context, scenarioID string, status models.ScenarioStatus, endedAt *time.Time) error {
	const q = `UPDATE scenarios SET status = $1, ended_at = $2 WHERE scenario_id = $3`
	tag, err := r.pool.Exec(ctx, q, string(status), endedAt, scenarioID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *pgxScenarioRepository) GetLatest(ctx context.Context) (*models.ScenarioRun, error) {
	const q = `
		SELECT scenario_id, scenario_type, transaction_count, seed, status, started_at, ended_at
		FROM scenarios
		ORDER BY started_at DESC
		LIMIT 1`
	row := r.pool.QueryRow(ctx, q)
	var s models.ScenarioRun
	var statusStr, typeStr string
	var endedAt *time.Time
	err := row.Scan(
		&s.ScenarioID, &typeStr, &s.TransactionCount,
		&s.Seed, &statusStr, &s.StartedAt, &endedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	s.ScenarioType = models.ScenarioType(typeStr)
	s.Status = models.ScenarioStatus(statusStr)
	s.EndedAt = endedAt
	return &s, nil
}

func (r *pgxScenarioRepository) InsertAuditEvent(ctx context.Context, transactionID, eventType string, details any) error {
	b, err := json.Marshal(details)
	if err != nil {
		return err
	}
	const q = `
		INSERT INTO audit_events (transaction_id, event_type, details, timestamp)
		VALUES ($1, $2, $3::jsonb, $4)`
	_, err = r.pool.Exec(ctx, q, transactionID, eventType, string(b), time.Now().UTC())
	return err
}
