package scenarios

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vulcanshield/backend/internal/models"
)

// DeviceFarmScenario simulates multiple distinct users sharing the same device,
// which creates suspicious graph relationships detectable in Phase 11.
type DeviceFarmScenario struct{}

func (d *DeviceFarmScenario) Type() models.ScenarioType { return models.ScenarioDeviceFarm }

func (d *DeviceFarmScenario) Next(idx int, rng *rand.Rand, pool *models.EntityPool, targetUserIndex int) models.Transaction {
	// Cycle through all users — each user uses the same shared device.
	userIdx := idx % len(pool.Users)
	user := pool.Users[userIdx]

	// All transactions share the last device (the high-risk emulator).
	// This creates a SHARED_DEVICE graph relationship across multiple users.
	sharedDeviceID := pool.DeviceIDs[len(pool.DeviceIDs)-1]

	ipAddr := pool.IPAddresses[rng.Intn(len(pool.IPAddresses))]
	merchantID := pool.MerchantIDs[rng.Intn(len(pool.MerchantIDs))]
	amount := randAmount(rng, user.TypicalMinAmount, user.TypicalMaxAmount)
	now := time.Now().UTC()

	return models.Transaction{
		TransactionID: fmt.Sprintf("TX-DEV-%d-%05d", rng.Int63n(9000)+1000, idx),
		UserID:        user.UserID,
		DeviceID:      sharedDeviceID,
		IPAddress:     ipAddr,
		MerchantID:    merchantID,
		Amount:        amount,
		Currency:      "INR",
		Channel:       "MOBILE",
		Status:        models.StatusPending,
		Timestamp:     now,
		CreatedAt:     now,
	}
}
