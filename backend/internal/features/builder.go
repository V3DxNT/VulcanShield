package features

import (
	"github.com/vulcanshield/backend/internal/mlclient"
	"github.com/vulcanshield/backend/internal/models"
	appredis "github.com/vulcanshield/backend/internal/redis"
)


type FeatureBuilder struct{}


func NewFeatureBuilder() *FeatureBuilder {
	return &FeatureBuilder{}
}


func (b *FeatureBuilder) BuildVector(
	tx *models.Transaction,
	user *models.UserProfile,
	velocity *appredis.VelocitySignals,
	isEmulator bool,
	isVPN bool,
) *mlclient.PredictRequest {
	typicalMax := 250.0
	trustScore := 85

	if user != nil {
		if user.TypicalMaxAmount > 0 {
			typicalMax = user.TypicalMaxAmount
		}
		trustScore = user.TrustScore
	}

	var userCount, ipCount, devCount int64
	if velocity != nil {
		userCount = velocity.UserTxCount60s
		ipCount = velocity.IPTxCount60s
		devCount = velocity.DeviceTxCount60s
	}

	return &mlclient.PredictRequest{
		TransactionID:    tx.TransactionID,
		UserID:           tx.UserID,
		Amount:           tx.Amount,
		TypicalMaxAmount: typicalMax,
		UserTxCount60s:   userCount,
		IPTxCount60s:     ipCount,
		DeviceTxCount60s: devCount,
		TrustScore:       trustScore,
		IsEmulator:       isEmulator,
		IsVPN:            isVPN,
	}
}
