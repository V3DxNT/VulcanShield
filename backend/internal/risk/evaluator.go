package risk

import (
	"fmt"
	"math"
	"time"

	"github.com/vulcanshield/backend/internal/mlclient"
	"github.com/vulcanshield/backend/internal/models"
	appredis "github.com/vulcanshield/backend/internal/redis"
)


type Evaluator struct{}


func NewEvaluator() *Evaluator {
	return &Evaluator{}
}


func (e *Evaluator) Evaluate(
	tx *models.Transaction,
	mlResp *mlclient.PredictResponse,
	velocity *appredis.VelocitySignals,
) *models.RiskAssessment {
	fraudProb := 0.05
	anomalyScore := 0.05
	fraudModelVer := "xgboost-v1"
	anomalyModelVer := "isoforest-v1"
	snapshot := make(map[string]any)

	if mlResp != nil {
		fraudProb = mlResp.FraudProbability
		anomalyScore = mlResp.AnomalyScore
		if mlResp.ModelVersion != "" {
			fraudModelVer = mlResp.ModelVersion
			anomalyModelVer = mlResp.ModelVersion
		}
		if mlResp.FeatureSnapshot != nil {
			snapshot = mlResp.FeatureSnapshot
		}
	}

	
	var velocityScore float64
	if velocity != nil {
		snapshot["user_tx_count_60s"] = velocity.UserTxCount60s
		snapshot["ip_tx_count_60s"] = velocity.IPTxCount60s
		snapshot["device_tx_count_60s"] = velocity.DeviceTxCount60s

		if velocity.UserTxCount60s > 3 {
			velocityScore += 0.4
		}
		if velocity.IPTxCount60s > 3 {
			velocityScore += 0.3
		}
		if velocity.DeviceTxCount60s > 2 {
			velocityScore += 0.3
		}
		if velocityScore > 1.0 {
			velocityScore = 1.0
		}
	}

	
	
	
	
	
	rawScore := (fraudProb * 40.0) + (anomalyScore * 20.0) + (velocityScore * 40.0)
	normalizedRiskScore := int(math.Round(math.Max(0.0, math.Min(100.0, rawScore))))

	now := time.Now().UTC()
	return &models.RiskAssessment{
		AssessmentID:        fmt.Sprintf("RA-%s", tx.TransactionID),
		TransactionID:       tx.TransactionID,
		FraudProbability:    fraudProb,
		AnomalyScore:        anomalyScore,
		FraudModelVersion:   fraudModelVer,
		AnomalyModelVersion: anomalyModelVer,
		RiskScore:           normalizedRiskScore,
		FeatureSnapshot:     snapshot,
		CreatedAt:           now,
	}
}
