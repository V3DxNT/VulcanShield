package mlclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)


type PredictRequest struct {
	TransactionID    string  `json:"transaction_id"`
	UserID           string  `json:"user_id"`
	Amount           float64 `json:"amount"`
	TypicalMaxAmount float64 `json:"typical_max_amount"`
	UserTxCount60s   int64   `json:"user_tx_count_60s"`
	IPTxCount60s     int64   `json:"ip_tx_count_60s"`
	DeviceTxCount60s int64   `json:"device_tx_count_60s"`
	TrustScore       int     `json:"trust_score"`
	IsEmulator       bool    `json:"is_emulator"`
	IsVPN            bool    `json:"is_vpn"`
}


type PredictResponse struct {
	TransactionID    string         `json:"transaction_id"`
	FraudProbability float64        `json:"fraud_probability"`
	AnomalyScore     float64        `json:"anomaly_score"`
	ModelVersion     string         `json:"model_version"`
	FeatureSnapshot  map[string]any `json:"feature_snapshot"`
}


type Client struct {
	baseURL    string
	httpClient *http.Client
}


func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}



func (c *Client) Predict(ctx context.Context, req *PredictRequest) (*PredictResponse, error) {
	url := fmt.Sprintf("%s/predict", c.baseURL)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshalling predict request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("creating http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		
		return fallbackPredict(req), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fallbackPredict(req), nil
	}

	var res PredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fallbackPredict(req), nil
	}

	return &res, nil
}


func fallbackPredict(req *PredictRequest) *PredictResponse {
	var fraudProb float64
	var anomalyScore float64

	ratio := req.Amount / (req.TypicalMaxAmount + 1e-5)
	if ratio > 3.0 {
		fraudProb += 0.5
		anomalyScore += 0.6
	}
	if req.UserTxCount60s > 5 {
		fraudProb += 0.4
	}
	if req.IsEmulator {
		fraudProb += 0.3
		anomalyScore += 0.3
	}

	if fraudProb > 1.0 {
		fraudProb = 1.0
	}
	if anomalyScore > 1.0 {
		anomalyScore = 1.0
	}

	return &PredictResponse{
		TransactionID:    req.TransactionID,
		FraudProbability: fraudProb,
		AnomalyScore:     anomalyScore,
		ModelVersion:     "v1.0-fallback",
		FeatureSnapshot: map[string]any{
			"fallback": true,
			"amount":   req.Amount,
		},
	}
}
