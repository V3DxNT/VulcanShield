package scenarios

import (
	"fmt"
	"math/rand"

	"github.com/vulcanshield/backend/internal/models"
)




type IPAbuseScenario struct{}

func (ip *IPAbuseScenario) Type() models.ScenarioType { return models.ScenarioIPAbuse }

func (ip *IPAbuseScenario) Next(idx int, rng *rand.Rand, pool *models.EntityPool, targetUserIndex int) models.Transaction {
	userIdx := idx % len(pool.Users)
	user := pool.Users[userIdx]

	sharedIP := pool.IPAddresses[len(pool.IPAddresses)-1]
	deviceID := pool.DeviceIDs[(idx+1)%len(pool.DeviceIDs)]
	merchantID := pool.MerchantIDs[(idx+2)%len(pool.MerchantIDs)]
	amount := randAmount(rng, user.TypicalMinAmount, user.TypicalMaxAmount)
	now := progressiveTimestamp(idx, rng)

	return models.Transaction{
		TransactionID: fmt.Sprintf("TX-IP-%d-%05d", rng.Int63n(9000)+1000, idx),
		UserID:        user.UserID,
		DeviceID:      deviceID,
		IPAddress:     sharedIP,
		MerchantID:    merchantID,
		Amount:        amount,
		Currency:      "INR",
		Channel:       "WEB",
		Status:        models.StatusPending,
		Timestamp:     now,
		CreatedAt:     now,
	}
}
