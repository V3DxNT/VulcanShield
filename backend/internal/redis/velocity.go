package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vulcanshield/backend/internal/models"
)


type VelocitySignals struct {
	UserTxCount60s   int64   `json:"user_tx_count_60s"`
	IPTxCount60s     int64   `json:"ip_tx_count_60s"`
	DeviceTxCount60s int64   `json:"device_tx_count_60s"`
	UserAmountSum60s float64 `json:"user_amount_sum_60s"`
}


type VelocityEngine struct {
	client *redis.Client
}


func NewVelocityEngine(client *redis.Client) *VelocityEngine {
	return &VelocityEngine{client: client}
}



func (v *VelocityEngine) RecordTransaction(ctx context.Context, tx *models.Transaction) error {
	if v.client == nil {
		return nil 
	}

	score := float64(tx.Timestamp.UnixMilli())
	member := fmt.Sprintf("%s:%f", tx.TransactionID, tx.Amount)

	pipe := v.client.Pipeline()

	userKey := fmt.Sprintf("velocity:user:%s", tx.UserID)
	ipKey := fmt.Sprintf("velocity:ip:%s", tx.IPAddress)
	deviceKey := fmt.Sprintf("velocity:device:%s", tx.DeviceID)

	pipe.ZAdd(ctx, userKey, redis.Z{Score: score, Member: member})
	pipe.Expire(ctx, userKey, 120*time.Second)

	pipe.ZAdd(ctx, ipKey, redis.Z{Score: score, Member: member})
	pipe.Expire(ctx, ipKey, 120*time.Second)

	pipe.ZAdd(ctx, deviceKey, redis.Z{Score: score, Member: member})
	pipe.Expire(ctx, deviceKey, 120*time.Second)

	_, err := pipe.Exec(ctx)
	return err
}


func (v *VelocityEngine) GetVelocitySignals(ctx context.Context, tx *models.Transaction) (*VelocitySignals, error) {
	if v.client == nil {
		return &VelocitySignals{}, nil 
	}

	nowMs := tx.Timestamp.UnixMilli()
	if nowMs == 0 {
		nowMs = time.Now().UnixMilli()
	}
	minScore := strconv.FormatInt(nowMs-60000, 10)
	maxScore := strconv.FormatInt(nowMs, 10)

	userKey := fmt.Sprintf("velocity:user:%s", tx.UserID)
	ipKey := fmt.Sprintf("velocity:ip:%s", tx.IPAddress)
	deviceKey := fmt.Sprintf("velocity:device:%s", tx.DeviceID)

	pipe := v.client.Pipeline()

	
	pipe.ZRemRangeByScore(ctx, userKey, "-inf", "("+minScore)
	pipe.ZRemRangeByScore(ctx, ipKey, "-inf", "("+minScore)
	pipe.ZRemRangeByScore(ctx, deviceKey, "-inf", "("+minScore)

	userCountCmd := pipe.ZCount(ctx, userKey, minScore, maxScore)
	ipCountCmd := pipe.ZCount(ctx, ipKey, minScore, maxScore)
	deviceCountCmd := pipe.ZCount(ctx, deviceKey, minScore, maxScore)
	userMembersCmd := pipe.ZRangeByScore(ctx, userKey, &redis.ZRangeBy{Min: minScore, Max: maxScore})

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return &VelocitySignals{}, nil 
	}

	var amountSum float64
	for _, m := range userMembersCmd.Val() {
		
		var amount float64
		var txID string
		if n, _ := fmt.Sscanf(m, "%s:%f", &txID, &amount); n == 2 {
			amountSum += amount
		}
	}

	return &VelocitySignals{
		UserTxCount60s:   userCountCmd.Val(),
		IPTxCount60s:     ipCountCmd.Val(),
		DeviceTxCount60s: deviceCountCmd.Val(),
		UserAmountSum60s: amountSum,
	}, nil
}
