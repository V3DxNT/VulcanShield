package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vulcanshield/backend/internal/models"
)


type UserRepository interface {
	
	
	GetByID(ctx context.Context, userID string) (*models.User, error)
}


type pgxUserRepository struct {
	pool *pgxpool.Pool
}


func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &pgxUserRepository{pool: pool}
}



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
