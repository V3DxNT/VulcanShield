package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vulcanshield/backend/internal/models"
)


type RiskRepository interface {
	Create(ctx context.Context, ra *models.RiskAssessment) error
	GetByTransactionID(ctx context.Context, txID string) (*models.RiskAssessment, error)
}

type pgxRiskRepository struct {
	pool *pgxpool.Pool
}


func NewRiskRepository(pool *pgxpool.Pool) RiskRepository {
	return &pgxRiskRepository{pool: pool}
}

func (r *pgxRiskRepository) Create(ctx context.Context, ra *models.RiskAssessment) error {
	snapshotJSON, err := json.Marshal(ra.FeatureSnapshot)
	if err != nil {
		snapshotJSON = []byte("{}")
	}

	const q = `
		INSERT INTO risk_assessments
		  (assessment_id, transaction_id, fraud_probability, anomaly_score,
		   fraud_model_version, anomaly_model_version, risk_score, feature_snapshot, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb, $9)`

	_, err = r.pool.Exec(ctx, q,
		ra.AssessmentID,
		ra.TransactionID,
		ra.FraudProbability,
		ra.AnomalyScore,
		ra.FraudModelVersion,
		ra.AnomalyModelVersion,
		ra.RiskScore,
		string(snapshotJSON),
		ra.CreatedAt,
	)
	return err
}

func (r *pgxRiskRepository) GetByTransactionID(ctx context.Context, txID string) (*models.RiskAssessment, error) {
	const q = `
		SELECT assessment_id, transaction_id, fraud_probability, anomaly_score,
		       fraud_model_version, anomaly_model_version, risk_score, feature_snapshot, created_at
		FROM risk_assessments
		WHERE transaction_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	var ra models.RiskAssessment
	var snapshotBytes []byte
	row := r.pool.QueryRow(ctx, q, txID)
	err := row.Scan(
		&ra.AssessmentID,
		&ra.TransactionID,
		&ra.FraudProbability,
		&ra.AnomalyScore,
		&ra.FraudModelVersion,
		&ra.AnomalyModelVersion,
		&ra.RiskScore,
		&snapshotBytes,
		&ra.CreatedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	_ = json.Unmarshal(snapshotBytes, &ra.FeatureSnapshot)
	return &ra, nil
}
