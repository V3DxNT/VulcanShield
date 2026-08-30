package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vulcanshield/backend/internal/models"
)

// PolicyRepository defines data access for policy decisions.
type PolicyRepository interface {
	Create(ctx context.Context, pd *models.PolicyDecision) error
	GetByTransactionID(ctx context.Context, txID string) (*models.PolicyDecision, error)
}

type pgxPolicyRepository struct {
	pool *pgxpool.Pool
}

// NewPolicyRepository returns a PostgreSQL-backed PolicyRepository.
func NewPolicyRepository(pool *pgxpool.Pool) PolicyRepository {
	return &pgxPolicyRepository{pool: pool}
}

func (r *pgxPolicyRepository) Create(ctx context.Context, pd *models.PolicyDecision) error {
	rulesJSON, err := json.Marshal(pd.RulesTriggered)
	if err != nil {
		rulesJSON = []byte("[]")
	}

	const q = `
		INSERT INTO policy_decisions (decision_id, transaction_id, decision, created_at)
		VALUES ($1, $2, $3, $4)`

	_, err = r.pool.Exec(ctx, q,
		pd.DecisionID,
		pd.TransactionID,
		string(pd.Decision),
		pd.CreatedAt,
	)
	_ = rulesJSON
	return err
}

func (r *pgxPolicyRepository) GetByTransactionID(ctx context.Context, txID string) (*models.PolicyDecision, error) {
	const q = `
		SELECT decision_id, transaction_id, decision, created_at
		FROM policy_decisions
		WHERE transaction_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	var pd models.PolicyDecision
	var decStr string
	row := r.pool.QueryRow(ctx, q, txID)
	err := row.Scan(&pd.DecisionID, &pd.TransactionID, &decStr, &pd.CreatedAt)
	if err != nil {
		if isNoRows(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	pd.Decision = models.PolicyDecisionType(decStr)
	return &pd, nil
}
