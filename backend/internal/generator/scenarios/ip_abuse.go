package scenarios

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vulcanshield/backend/internal/models"
)

// IPAbuseScenario simulates multiple distinct users making transactions from
// the same high-risk IP (seeded as VPN/proxy), creating shared-IP fraud graph
// relationships exploited in Phase 11.
type IPAbuseScenario struct{}

func (ip *IPAbuseScenario) Type() models.ScenarioType { return models.ScenarioIPAbuse }

func (ip *IPAbuseScenario) Next(idx int, rng *rand.Rand, pool *models.EntityPool, targetUserIndex int) models.Transaction {
	// Cycle through all users — each uses the same shared high-risk IP.
	userIdx := idx % len(pool.Users)
	user := pool.Users[userIdx]

	// All transactions share the last IP (seeded as VPN/proxy/high-risk).
	sharedIP := pool.IPAddresses[len(pool.IPAddresses)-1]

	deviceID := pool.DeviceIDs[rng.Intn(len(pool.DeviceIDs))]
	merchantID := pool.MerchantIDs[rng.Intn(len(pool.MerchantIDs))]
	amount := randAmount(rng, user.TypicalMinAmount, user.TypicalMaxAmount)
	now := time.Now().UTC()

	return models.Transaction{
		TransactionID: fmt.Sprintf("TX-IP-%d-%05d", rng.Int63n(9000)+1000, idx),
		UserID:        user.UserID,
		DeviceID:      deviceID,
		IPAddress:     sharedIP,
		MerchantID:    merchantID,
		Amount:        amount,
		Currency:      "USD",
		Channel:       "WEB",
		Status:        models.StatusPending,
		Timestamp:     now,
		CreatedAt:     now,
	}
}
