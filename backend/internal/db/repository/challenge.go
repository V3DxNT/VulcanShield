package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vulcanshield/backend/internal/models"
)

type ChallengeRepository interface {
	Create(ctx context.Context, c *models.OTPChallenge) error
	GetByID(ctx context.Context, challengeID string) (*models.OTPChallenge, error)
	GetByTransactionID(ctx context.Context, txID string) (*models.OTPChallenge, error)
	Update(ctx context.Context, c *models.OTPChallenge) error
}

type pgxChallengeRepository struct {
	pool *pgxpool.Pool
}

func NewChallengeRepository(pool *pgxpool.Pool) ChallengeRepository {
	return &pgxChallengeRepository{pool: pool}
}

func (r *pgxChallengeRepository) Create(ctx context.Context, c *models.OTPChallenge) error {
	const q = `
		INSERT INTO otp_challenges
		  (challenge_id, transaction_id, otp_code_hash, status, attempts, max_attempts, created_at, expires_at, verified_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.pool.Exec(ctx, q,
		c.ChallengeID,
		c.TransactionID,
		c.OTPCodeHash,
		string(c.Status),
		c.Attempts,
		c.MaxAttempts,
		c.CreatedAt,
		c.ExpiresAt,
		c.VerifiedAt,
	)
	return err
}

func (r *pgxChallengeRepository) GetByID(ctx context.Context, challengeID string) (*models.OTPChallenge, error) {
	const q = `
		SELECT challenge_id, transaction_id, otp_code_hash, status, attempts, max_attempts, created_at, expires_at, verified_at
		FROM otp_challenges
		WHERE challenge_id = $1`

	var c models.OTPChallenge
	var statusStr string
	var verifiedAt *time.Time

	row := r.pool.QueryRow(ctx, q, challengeID)
	err := row.Scan(
		&c.ChallengeID, &c.TransactionID, &c.OTPCodeHash, &statusStr,
		&c.Attempts, &c.MaxAttempts, &c.CreatedAt, &c.ExpiresAt, &verifiedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	c.Status = models.ChallengeStatus(statusStr)
	c.VerifiedAt = verifiedAt
	return &c, nil
}

func (r *pgxChallengeRepository) GetByTransactionID(ctx context.Context, txID string) (*models.OTPChallenge, error) {
	const q = `
		SELECT challenge_id, transaction_id, otp_code_hash, status, attempts, max_attempts, created_at, expires_at, verified_at
		FROM otp_challenges
		WHERE transaction_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	var c models.OTPChallenge
	var statusStr string
	var verifiedAt *time.Time

	row := r.pool.QueryRow(ctx, q, txID)
	err := row.Scan(
		&c.ChallengeID, &c.TransactionID, &c.OTPCodeHash, &statusStr,
		&c.Attempts, &c.MaxAttempts, &c.CreatedAt, &c.ExpiresAt, &verifiedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	c.Status = models.ChallengeStatus(statusStr)
	c.VerifiedAt = verifiedAt
	return &c, nil
}

func (r *pgxChallengeRepository) Update(ctx context.Context, c *models.OTPChallenge) error {
	const q = `
		UPDATE otp_challenges
		SET status = $1, attempts = $2, verified_at = $3
		WHERE challenge_id = $4`

	_, err := r.pool.Exec(ctx, q, string(c.Status), c.Attempts, c.VerifiedAt, c.ChallengeID)
	return err
}
