package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vulcanshield/backend/internal/models"
)

// ErrNotFound is returned when a requested record does not exist.
var ErrNotFound = errors.New("record not found")

// TransactionRepository defines the data access interface for transactions.
// Phase 4 (Transaction Generator) will provide a full implementation.
type TransactionRepository interface {
	// GetByID retrieves a single transaction by its ID.
	GetByID(ctx context.Context, id string) (*models.Transaction, error)

	// ListByUser retrieves recent transactions for a user, ordered by timestamp desc.
	ListByUser(ctx context.Context, userID string, limit int) ([]*models.Transaction, error)

	// Create persists a new transaction record.
	Create(ctx context.Context, tx *models.Transaction) error

	// UpdateStatus updates the lifecycle status of a transaction.
	UpdateStatus(ctx context.Context, id string, status models.TransactionStatus) error
}

// pgxTransactionRepository is the PostgreSQL-backed implementation.
type pgxTransactionRepository struct {
	pool *pgxpool.Pool
}

// NewTransactionRepository returns a PostgreSQL-backed TransactionRepository.
func NewTransactionRepository(pool *pgxpool.Pool) TransactionRepository {
	return &pgxTransactionRepository{pool: pool}
}

// GetByID retrieves a transaction by primary key.
// Returns ErrNotFound if no matching record exists.
func (r *pgxTransactionRepository) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	const q = `
		SELECT transaction_id, user_id, device_id, ip_address, merchant_id,
		       amount, currency, channel, status, timestamp, created_at
		FROM transactions
		WHERE transaction_id = $1`

	row := r.pool.QueryRow(ctx, q, id)
	var t models.Transaction
	err := row.Scan(
		&t.TransactionID, &t.UserID, &t.DeviceID, &t.IPAddress, &t.MerchantID,
		&t.Amount, &t.Currency, &t.Channel, &t.Status, &t.Timestamp, &t.CreatedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

// ListByUser returns up to limit transactions for userID, newest first.
func (r *pgxTransactionRepository) ListByUser(ctx context.Context, userID string, limit int) ([]*models.Transaction, error) {
	const q = `
		SELECT transaction_id, user_id, device_id, ip_address, merchant_id,
		       amount, currency, channel, status, timestamp, created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2`

	rows, err := r.pool.Query(ctx, q, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []*models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(
			&t.TransactionID, &t.UserID, &t.DeviceID, &t.IPAddress, &t.MerchantID,
			&t.Amount, &t.Currency, &t.Channel, &t.Status, &t.Timestamp, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		txs = append(txs, &t)
	}
	return txs, rows.Err()
}

// Create inserts a new transaction into PostgreSQL.
func (r *pgxTransactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
	const q = `
		INSERT INTO transactions
		  (transaction_id, user_id, device_id, ip_address, merchant_id,
		   amount, currency, channel, status, timestamp, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

	_, err := r.pool.Exec(ctx, q,
		tx.TransactionID, tx.UserID, tx.DeviceID, tx.IPAddress, tx.MerchantID,
		tx.Amount, tx.Currency, tx.Channel, tx.Status, tx.Timestamp, tx.CreatedAt,
	)
	return err
}

// UpdateStatus updates the lifecycle status of an existing transaction.
func (r *pgxTransactionRepository) UpdateStatus(ctx context.Context, id string, status models.TransactionStatus) error {
	const q = `UPDATE transactions SET status = $1 WHERE transaction_id = $2`
	tag, err := r.pool.Exec(ctx, q, status, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// isNoRows returns true for pgx "no rows" errors.
func isNoRows(err error) bool {
	return err != nil && err.Error() == "no rows in result set"
}
