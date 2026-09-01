package scenarios

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vulcanshield/backend/internal/models"
)



type NormalScenario struct{}

func (n *NormalScenario) Type() models.ScenarioType { return models.ScenarioNormal }

func (n *NormalScenario) Next(idx int, rng *rand.Rand, pool *models.EntityPool, targetUserIndex int) models.Transaction {
	userIdx := idx % len(pool.Users)
	if targetUserIndex >= 0 && targetUserIndex < len(pool.Users) {
		userIdx = targetUserIndex
	}
	user := pool.Users[userIdx]

	amount := randAmount(rng, user.TypicalMinAmount, user.TypicalMaxAmount)
	deviceID := userDeviceID(pool, user)
	ipAddr := userIP(pool, user)
	merchantIdx := (idx + 2 + rng.Intn(len(pool.MerchantIDs))) % len(pool.MerchantIDs)
	now := progressiveTimestamp(idx, rng)

	return models.Transaction{
		TransactionID: fmt.Sprintf("TX-%d-%05d", rng.Int63n(9000)+1000, idx),
		UserID:        user.UserID,
		DeviceID:      deviceID,
		IPAddress:     ipAddr,
		MerchantID:    pool.MerchantIDs[merchantIdx],
		Amount:        amount,
		Currency:      "INR",
		Channel:       "WEB",
		Status:        models.StatusPending,
		Timestamp:     now,
		CreatedAt:     now,
	}
}

func progressiveTimestamp(idx int, rng *rand.Rand) time.Time {
	if idx < 0 {
		idx = 0
	}
	return time.Now().UTC().Add(time.Duration(idx)*time.Second + time.Duration(rng.Intn(250))*time.Millisecond)
}


func randAmount(rng *rand.Rand, min, max float64) float64 {
	if min >= max {
		return min
	}
	raw := min + rng.Float64()*(max-min)
	return float64(int(raw*100)) / 100
}
