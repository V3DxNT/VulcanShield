package risk

import (
	"fmt"
	"math"
	"time"

	"github.com/vulcanshield/backend/internal/mlclient"
	"github.com/vulcanshield/backend/internal/models"
	appredis "github.com/vulcanshield/backend/internal/redis"
)

// Evaluator combines ML predictions and behavioral signals into a 0-100 normalized risk score.
type Evaluator struct{}

// NewEvaluator creates a new Risk Evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// Evaluate computes the standardized 0-100 risk score and constructs the RiskAssessment domain model.
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

	// Calculate velocity contribution factor (0.0 to 1.0)
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

	// Risk Score Contract (0 - 100) per PROJECT_SPEC.md §18.
	// A fully saturated 60-second burst is itself decisive behavioral evidence,
	// so it receives enough weight to trigger policy step-up verification even
	// when a low-value card-testing attempt has a modest standalone ML score.
	// ML still provides the majority of non-velocity risk.
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
