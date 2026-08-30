package scenarios

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vulcanshield/backend/internal/models"
)

// AccountTakeoverScenario simulates an attacker using a new (high-risk) device
// and high-risk IP to make an unusually large transaction on a legitimate user's
// account — creating clear behavioural deviation from historical patterns.
type AccountTakeoverScenario struct{}

func (a *AccountTakeoverScenario) Type() models.ScenarioType { return models.ScenarioAccountTakeover }

func (a *AccountTakeoverScenario) Next(idx int, rng *rand.Rand, pool *models.EntityPool, targetUserIndex int) models.Transaction {
	// Target a specific user (default to first user with good trust score).
	userIdx := 0
	if targetUserIndex >= 0 && targetUserIndex < len(pool.Users) {
		userIdx = targetUserIndex
	}
	user := pool.Users[userIdx]

	// Use the last device (seeded as emulator/high-risk) and last IP (seeded as VPN/high-risk).
	deviceID := pool.DeviceIDs[len(pool.DeviceIDs)-1]
	ipAddr := pool.IPAddresses[len(pool.IPAddresses)-1]
	merchantID := pool.MerchantIDs[rng.Intn(len(pool.MerchantIDs))]

	// Amount is significantly above user's typical maximum — behavioural anomaly.
	multiplier := 5.0 + rng.Float64()*10.0 // 5x–15x typical max
	amount := randAmount(rng, user.TypicalMaxAmount*multiplier*0.8, user.TypicalMaxAmount*multiplier)

	now := time.Now().UTC()

	return models.Transaction{
		TransactionID: fmt.Sprintf("TX-ATO-%d-%05d", rng.Int63n(9000)+1000, idx),
		UserID:        user.UserID,
		DeviceID:      deviceID,
		IPAddress:     ipAddr,
		MerchantID:    merchantID,
		Amount:        amount,
		Currency:      "USD",
		Channel:       "WEB",
		Status:        models.StatusPending,
		Timestamp:     now,
		CreatedAt:     now,
	}
}
