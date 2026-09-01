package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vulcanshield/backend/internal/models"
)



type EntityRepository interface {
	LoadPool(ctx context.Context) (*models.EntityPool, error)
}

type pgxEntityRepository struct {
	pool *pgxpool.Pool
}


func NewEntityRepository(pool *pgxpool.Pool) EntityRepository {
	return &pgxEntityRepository{pool: pool}
}



func (r *pgxEntityRepository) LoadPool(ctx context.Context) (*models.EntityPool, error) {
	ep := &models.EntityPool{}

	
	userRows, err := r.pool.Query(ctx,
		`SELECT user_id, typical_min_amount, typical_max_amount, trust_score,
		        challenge_threshold, block_threshold
		 FROM users ORDER BY user_id`)
	if err != nil {
		return nil, err
	}
	defer userRows.Close()
	for userRows.Next() {
		var u models.UserProfile
		if err := userRows.Scan(
			&u.UserID, &u.TypicalMinAmount, &u.TypicalMaxAmount, &u.TrustScore,
			&u.ChallengeThreshold, &u.BlockThreshold,
		); err != nil {
			return nil, err
		}
		ep.Users = append(ep.Users, u)
	}
	if err := userRows.Err(); err != nil {
		return nil, err
	}

	
	ep.UserDeviceMap = map[string]string{}
	ep.UserIPMap = map[string]string{}

	deviceRows, err := r.pool.Query(ctx, `SELECT user_id, device_id FROM user_devices ORDER BY user_id, device_id`)
	if err != nil {
		return nil, err
	}
	defer deviceRows.Close()
	for deviceRows.Next() {
		var userID, deviceID string
		if err := deviceRows.Scan(&userID, &deviceID); err != nil {
			return nil, err
		}
		ep.UserDeviceMap[userID] = deviceID
	}
	if err := deviceRows.Err(); err != nil {
		return nil, err
	}

	ipRows, err := r.pool.Query(ctx, `SELECT user_id, ip_address FROM user_ips ORDER BY user_id, ip_address`)
	if err != nil {
		return nil, err
	}
	defer ipRows.Close()
	for ipRows.Next() {
		var userID, ipAddress string
		if err := ipRows.Scan(&userID, &ipAddress); err != nil {
			return nil, err
		}
		ep.UserIPMap[userID] = ipAddress
	}
	if err := ipRows.Err(); err != nil {
		return nil, err
	}

	
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

	
	ipAddrRows, err := r.pool.Query(ctx, `SELECT ip_address FROM ips ORDER BY ip_address`)
	if err != nil {
		return nil, err
	}
	defer ipAddrRows.Close()
	for ipAddrRows.Next() {
		var ip string
		if err := ipAddrRows.Scan(&ip); err != nil {
			return nil, err
		}
		ep.IPAddresses = append(ep.IPAddresses, ip)
	}
	if err := ipAddrRows.Err(); err != nil {
		return nil, err
	}

	
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
