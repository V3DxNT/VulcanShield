package scenarios

import (
	"fmt"
	"math/rand"

	"github.com/vulcanshield/backend/internal/models"
)

// AccountTakeoverScenario simulates an attacker using a new (high-risk) device
// and high-risk IP to make an unusually large transaction on a legitimate user's
// account — creating clear behavioural deviation from historical patterns.
type AccountTakeoverScenario struct{}

func (a *AccountTakeoverScenario) Type() models.ScenarioType { return models.ScenarioAccountTakeover }

func (a *AccountTakeoverScenario) Next(idx int, rng *rand.Rand, pool *models.EntityPool, targetUserIndex int) models.Transaction {
	userIdx := idx % len(pool.Users)
	if targetUserIndex >= 0 && targetUserIndex < len(pool.Users) {
		userIdx = targetUserIndex
	}
	user := pool.Users[userIdx]

	deviceID := userDeviceID(pool, user)
	ipAddr := userIP(pool, user)
	if len(pool.DeviceIDs) > 0 {
		deviceID = pool.DeviceIDs[len(pool.DeviceIDs)-1]
	}
	if len(pool.IPAddresses) > 0 {
		ipAddr = pool.IPAddresses[len(pool.IPAddresses)-1]
	}
	merchantID := pool.MerchantIDs[(idx+2)%len(pool.MerchantIDs)]

	multiplier := 5.0 + rng.Float64()*10.0
	amount := randAmount(rng, user.TypicalMaxAmount*multiplier*0.8, user.TypicalMaxAmount*multiplier)

	now := progressiveTimestamp(idx, rng)

	return models.Transaction{
		TransactionID: fmt.Sprintf("TX-ATO-%d-%05d", rng.Int63n(9000)+1000, idx),
		UserID:        user.UserID,
		DeviceID:      deviceID,
		IPAddress:     ipAddr,
		MerchantID:    merchantID,
		Amount:        amount,
		Currency:      "INR",
		Channel:       "WEB",
		Status:        models.StatusPending,
		Timestamp:     now,
		CreatedAt:     now,
	}
}
