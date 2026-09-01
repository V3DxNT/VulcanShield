package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vulcanshield/backend/internal/models"
)


var ErrNotFound = errors.New("record not found")


type TransactionRepository interface {
	
	GetByID(ctx context.Context, id string) (*models.Transaction, error)

	
	ListByUser(ctx context.Context, userID string, limit int) ([]*models.Transaction, error)

	
	ListRecent(ctx context.Context, limit int) ([]*models.Transaction, error)

	
	Create(ctx context.Context, tx *models.Transaction) error

	
	UpdateStatus(ctx context.Context, id string, status models.TransactionStatus) error
}


type pgxTransactionRepository struct {
	pool *pgxpool.Pool
}


func NewTransactionRepository(pool *pgxpool.Pool) TransactionRepository {
	return &pgxTransactionRepository{pool: pool}
}



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


func (r *pgxTransactionRepository) ListRecent(ctx context.Context, limit int) ([]*models.Transaction, error) {
	const q = `
		SELECT transaction_id, user_id, device_id, ip_address, merchant_id,
		       amount, currency, channel, status, timestamp, created_at
		FROM transactions
		ORDER BY timestamp DESC
		LIMIT $1`

	rows, err := r.pool.Query(ctx, q, limit)
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


func isNoRows(err error) bool {
	return err != nil && err.Error() == "no rows in result set"
}
