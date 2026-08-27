package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vulcanshield/backend/internal/models"
)

// UserRepository defines data access for user profiles.
type UserRepository interface {
	// GetByID retrieves a user by their user_id.
	// Returns ErrNotFound if no user exists with that ID.
	GetByID(ctx context.Context, userID string) (*models.User, error)
}

// pgxUserRepository is the PostgreSQL-backed implementation.
type pgxUserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository returns a PostgreSQL-backed UserRepository.
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &pgxUserRepository{pool: pool}
}

// GetByID retrieves a user by primary key.
// challenge_threshold and block_threshold are the authoritative per-user policy params.
func (r *pgxUserRepository) GetByID(ctx context.Context, userID string) (*models.User, error) {
	const q = `
		SELECT user_id, name, email, risk_tolerance,
		       challenge_threshold, block_threshold,
		       account_age_days, trust_score,
		       typical_min_amount, typical_max_amount,
		       created_at, updated_at
		FROM users
		WHERE user_id = $1`

	row := r.pool.QueryRow(ctx, q, userID)
	var u models.User
	err := row.Scan(
		&u.UserID, &u.Name, &u.Email, &u.RiskTolerance,
		&u.ChallengeThreshold, &u.BlockThreshold,
		&u.AccountAgeDays, &u.TrustScore,
		&u.TypicalMinAmount, &u.TypicalMaxAmount,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
