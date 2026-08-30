package scenarios

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vulcanshield/backend/internal/models"
)

// NormalScenario generates realistic distributed payment traffic across all
// seeded users with amounts within each user's typical range.
type NormalScenario struct{}

func (n *NormalScenario) Type() models.ScenarioType { return models.ScenarioNormal }

func (n *NormalScenario) Next(idx int, rng *rand.Rand, pool *models.EntityPool, targetUserIndex int) models.Transaction {
	userIdx := idx % len(pool.Users)
	if targetUserIndex >= 0 && targetUserIndex < len(pool.Users) {
		userIdx = targetUserIndex
	}
	user := pool.Users[userIdx]

	amount := randAmount(rng, user.TypicalMinAmount, user.TypicalMaxAmount)
	deviceIdx := rng.Intn(len(pool.DeviceIDs))
	ipIdx := rng.Intn(len(pool.IPAddresses))
	merchantIdx := rng.Intn(len(pool.MerchantIDs))
	now := time.Now().UTC()

	return models.Transaction{
		TransactionID: fmt.Sprintf("TX-%d-%05d", rng.Int63n(9000)+1000, idx),
		UserID:        user.UserID,
		DeviceID:      pool.DeviceIDs[deviceIdx],
		IPAddress:     pool.IPAddresses[ipIdx],
		MerchantID:    pool.MerchantIDs[merchantIdx],
		Amount:        amount,
		Currency:      "USD",
		Channel:       "WEB",
		Status:        models.StatusPending,
		Timestamp:     now,
		CreatedAt:     now,
	}
}

// randAmount returns a random float64 in [min, max] rounded to 2 decimal places.
func randAmount(rng *rand.Rand, min, max float64) float64 {
	if min >= max {
		return min
	}
	raw := min + rng.Float64()*(max-min)
	return float64(int(raw*100)) / 100
}
