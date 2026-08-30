package scenarios

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vulcanshield/backend/internal/models"
)

// AmountAnomalyScenario generates transactions with amounts far exceeding
// a user's typical_max_amount — a strong ML anomaly signal for Isolation
// Forest and XGBoost in Phase 6.
type AmountAnomalyScenario struct{}

func (a *AmountAnomalyScenario) Type() models.ScenarioType { return models.ScenarioAmountAnomaly }

func (a *AmountAnomalyScenario) Next(idx int, rng *rand.Rand, pool *models.EntityPool, targetUserIndex int) models.Transaction {
	// Alternate between users to vary patterns.
	userIdx := idx % len(pool.Users)
	if targetUserIndex >= 0 && targetUserIndex < len(pool.Users) {
		userIdx = targetUserIndex
	}
	user := pool.Users[userIdx]

	// Amount is 10–20x the user's typical maximum — extreme anomaly.
	multiplier := 10.0 + rng.Float64()*10.0
	baseAmount := user.TypicalMaxAmount * multiplier
	amount := randAmount(rng, baseAmount*0.9, baseAmount)

	deviceID := pool.DeviceIDs[rng.Intn(len(pool.DeviceIDs))]
	ipAddr := pool.IPAddresses[rng.Intn(len(pool.IPAddresses))]
	// High-value merchants preferred for realism.
	merchantID := pool.MerchantIDs[len(pool.MerchantIDs)-1]
	now := time.Now().UTC()

	return models.Transaction{
		TransactionID: fmt.Sprintf("TX-AMT-%d-%05d", rng.Int63n(9000)+1000, idx),
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
