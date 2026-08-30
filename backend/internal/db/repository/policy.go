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
	// Keep the deterministic policy explanation with the decision.  The original
	// schema calls this JSONB column `reasons`; older rows contain just the rule
	// list, while new rows store both the human-readable reason and its rule IDs.
	reasonsJSON, err := json.Marshal(struct {
		Reason string   `json:"reason"`
		Rules  []string `json:"rules_triggered"`
	}{Reason: pd.Reason, Rules: pd.RulesTriggered})
	if err != nil {
		reasonsJSON = []byte(`{"reason":"","rules_triggered":[]}`)
	}

	version := pd.PolicyVersion
	if version == "" {
		version = "v1.0"
	}

	const q = `
		INSERT INTO policy_decisions (
			decision_id, transaction_id, decision,
			risk_score, challenge_threshold, block_threshold,
			policy_version, reasons, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err = r.pool.Exec(ctx, q,
		pd.DecisionID,
		pd.TransactionID,
		string(pd.Decision),
		pd.RiskScore,
		pd.ChallengeThreshold,
		pd.BlockThreshold,
		version,
		reasonsJSON,
		pd.CreatedAt,
	)
	return err
}

func (r *pgxPolicyRepository) GetByTransactionID(ctx context.Context, txID string) (*models.PolicyDecision, error) {
	const q = `
		SELECT decision_id, transaction_id, decision,
		       risk_score, challenge_threshold, block_threshold,
		       policy_version, reasons, created_at
		FROM policy_decisions
		WHERE transaction_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	var pd models.PolicyDecision
	var decStr string
	var reasons []byte
	row := r.pool.QueryRow(ctx, q, txID)
	err := row.Scan(
		&pd.DecisionID, &pd.TransactionID, &decStr,
		&pd.RiskScore, &pd.ChallengeThreshold, &pd.BlockThreshold,
		&pd.PolicyVersion, &reasons, &pd.CreatedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	pd.Decision = models.PolicyDecisionType(decStr)
	if len(reasons) > 0 {
		var stored struct {
			Reason string   `json:"reason"`
			Rules  []string `json:"rules_triggered"`
		}
		if err := json.Unmarshal(reasons, &stored); err == nil && (stored.Reason != "" || stored.Rules != nil) {
			pd.Reason = stored.Reason
			pd.RulesTriggered = stored.Rules
		} else {
			// Backward compatibility for pre-explanation rows.
			_ = json.Unmarshal(reasons, &pd.RulesTriggered)
		}
	}
	return &pd, nil
}
