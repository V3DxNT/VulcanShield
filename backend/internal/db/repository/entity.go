package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vulcanshield/backend/internal/models"
)

// EntityRepository loads synthetic seed entities from PostgreSQL.
// The generator uses this pool to build realistic transaction records.
type EntityRepository interface {
	LoadPool(ctx context.Context) (*models.EntityPool, error)
}

type pgxEntityRepository struct {
	pool *pgxpool.Pool
}

// NewEntityRepository returns a PostgreSQL-backed EntityRepository.
func NewEntityRepository(pool *pgxpool.Pool) EntityRepository {
	return &pgxEntityRepository{pool: pool}
}

// LoadPool fetches all seeded users, devices, IPs, and merchants from PostgreSQL.
// Called once at generator startup. Fields match the Phase 2 schema exactly.
func (r *pgxEntityRepository) LoadPool(ctx context.Context) (*models.EntityPool, error) {
	ep := &models.EntityPool{}

	// Load user profiles (minimal projection needed for generation)
	userRows, err := r.pool.Query(ctx,
		`SELECT user_id, typical_min_amount, typical_max_amount, trust_score FROM users ORDER BY user_id`)
	if err != nil {
		return nil, err
	}
	defer userRows.Close()
	for userRows.Next() {
		var u models.UserProfile
		if err := userRows.Scan(&u.UserID, &u.TypicalMinAmount, &u.TypicalMaxAmount, &u.TrustScore); err != nil {
			return nil, err
		}
		ep.Users = append(ep.Users, u)
	}
	if err := userRows.Err(); err != nil {
		return nil, err
	}

	// Load device IDs
	devRows, err := r.pool.Query(ctx, `SELECT device_id FROM devices ORDER BY device_id`)
	if err != nil {
		return nil, err
	}
	defer devRows.Close()
	for devRows.Next() {
		var id string
		if err := devRows.Scan(&id); err != nil {
			return nil, err
		}
		ep.DeviceIDs = append(ep.DeviceIDs, id)
	}
	if err := devRows.Err(); err != nil {
		return nil, err
	}

	// Load IP addresses
	ipRows, err := r.pool.Query(ctx, `SELECT ip_address FROM ips ORDER BY ip_address`)
	if err != nil {
		return nil, err
	}
	defer ipRows.Close()
	for ipRows.Next() {
		var ip string
		if err := ipRows.Scan(&ip); err != nil {
			return nil, err
		}
		ep.IPAddresses = append(ep.IPAddresses, ip)
	}
	if err := ipRows.Err(); err != nil {
		return nil, err
	}

	// Load merchant IDs
	mRows, err := r.pool.Query(ctx, `SELECT merchant_id FROM merchants ORDER BY merchant_id`)
	if err != nil {
		return nil, err
	}
	defer mRows.Close()
	for mRows.Next() {
		var id string
		if err := mRows.Scan(&id); err != nil {
			return nil, err
		}
		ep.MerchantIDs = append(ep.MerchantIDs, id)
	}
	if err := mRows.Err(); err != nil {
		return nil, err
	}

	return ep, nil
}
