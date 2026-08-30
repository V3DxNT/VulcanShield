package scenarios

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vulcanshield/backend/internal/models"
)

// VelocityAttackScenario generates rapid high-frequency transactions for a
// single target user from the same device and IP, simulating card testing
// or velocity attacks detectable by Redis sliding windows.
type VelocityAttackScenario struct{}

func (v *VelocityAttackScenario) Type() models.ScenarioType { return models.ScenarioVelocityAttack }

func (v *VelocityAttackScenario) Next(idx int, rng *rand.Rand, pool *models.EntityPool, targetUserIndex int) models.Transaction {
	// Pin to the first user (or specified target) — same user every transaction.
	userIdx := 0
	if targetUserIndex >= 0 && targetUserIndex < len(pool.Users) {
		userIdx = targetUserIndex
	}
	user := pool.Users[userIdx]

	// Pin device and IP — same values every transaction to trigger velocity signals.
	deviceID := pool.DeviceIDs[0]
	ipAddr := pool.IPAddresses[0]
	merchantID := pool.MerchantIDs[0]

	// Amount is low-to-medium, realistic for card testing.
	amount := randAmount(rng, user.TypicalMinAmount, user.TypicalMaxAmount*0.5)
	now := time.Now().UTC()

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
