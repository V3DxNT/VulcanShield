package scenarios

import (
	"fmt"
	"math/rand"

	"github.com/vulcanshield/backend/internal/models"
)

// VelocityAttackScenario generates rapid high-frequency transactions for a
// single target user from the same device and IP, simulating card testing
// or velocity attacks detectable by Redis sliding windows.
type VelocityAttackScenario struct{}

func (v *VelocityAttackScenario) Type() models.ScenarioType { return models.ScenarioVelocityAttack }

func (v *VelocityAttackScenario) Next(idx int, rng *rand.Rand, pool *models.EntityPool, targetUserIndex int) models.Transaction {
	userIdx := 0
	if targetUserIndex >= 0 && targetUserIndex < len(pool.Users) {
		userIdx = targetUserIndex
	}
	user := pool.Users[userIdx]

	deviceID := pool.DeviceIDs[userIdx%len(pool.DeviceIDs)]
	ipAddr := pool.IPAddresses[userIdx%len(pool.IPAddresses)]
	merchantID := pool.MerchantIDs[(idx+1)%len(pool.MerchantIDs)]

	amount := randAmount(rng, user.TypicalMinAmount, user.TypicalMaxAmount*0.5)
	now := progressiveTimestamp(idx, rng)

	return models.Transaction{
		TransactionID: fmt.Sprintf("TX-VEL-%d-%05d", rng.Int63n(9000)+1000, idx),
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
