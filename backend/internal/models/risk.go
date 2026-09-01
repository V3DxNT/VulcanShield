package models

import "time"



type RiskAssessment struct {
	AssessmentID        string         `json:"assessment_id"`
	TransactionID       string         `json:"transaction_id"`
	FraudProbability    float64        `json:"fraud_probability"`
	AnomalyScore        float64        `json:"anomaly_score"`
	FraudModelVersion   string         `json:"fraud_model_version"`
	AnomalyModelVersion string         `json:"anomaly_model_version"`
	RiskScore           int            `json:"risk_score"`
	FeatureSnapshot     map[string]any `json:"feature_snapshot"`
	CreatedAt           time.Time      `json:"created_at"`
}
