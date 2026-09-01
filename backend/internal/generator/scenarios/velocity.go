package scenarios

import (
	"fmt"
	"math/rand"

	"github.com/vulcanshield/backend/internal/models"
)




type VelocityAttackScenario struct{}

func (v *VelocityAttackScenario) Type() models.ScenarioType { return models.ScenarioVelocityAttack }

func (v *VelocityAttackScenario) Next(idx int, rng *rand.Rand, pool *models.EntityPool, targetUserIndex int) models.Transaction {
	userIdx := 0
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
