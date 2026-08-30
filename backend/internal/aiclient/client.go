package aiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type InvestigationRequest struct {
	TransactionID    string  `json:"transaction_id"`
	UserID           string  `json:"user_id"`
	DeviceID         string  `json:"device_id"`
	IPAddress        string  `json:"ip_address"`
	Amount           float64 `json:"amount"`
	RiskScore        int     `json:"risk_score"`
	FraudProbability float64 `json:"fraud_probability"`
	AnomalyScore     float64 `json:"anomaly_score"`
	Status           string  `json:"status,omitempty"`
	Decision         string  `json:"decision,omitempty"`
}

type EvidenceItem struct {
	Category string `json:"category"`
	Fact     string `json:"fact"`
	Severity string `json:"severity"`
}

type SimilarCase struct {
	CaseID         string  `json:"case_id"`
	Title          string  `json:"title"`
	RelevanceScore float64 `json:"relevance_score"`
}

type InvestigationResponse struct {
	InvestigationID   string         `json:"investigation_id"`
	TransactionID     string         `json:"transaction_id"`
	RiskLevel         string         `json:"risk_level"`
	Summary           string         `json:"summary"`
	Evidence          []EvidenceItem `json:"evidence"`
	SimilarCases      []SimilarCase  `json:"similar_cases"`
	RecommendedAction string         `json:"recommended_action"`
	Confidence        float64        `json:"confidence"`
	LLMModel          string         `json:"llm_model"`
	RiskScore         int            `json:"risk_score"`
	FraudProbability  float64        `json:"fraud_probability"`
	AnomalyScore      float64        `json:"anomaly_score"`
	PolicyDecision    string         `json:"policy_decision"`
	TransactionStatus string         `json:"transaction_status"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 25 * time.Second,
		},
	}
}

func (c *Client) Investigate(ctx context.Context, req *InvestigationRequest) (*InvestigationResponse, error) {
	url := fmt.Sprintf("%s/ai/investigate", c.baseURL)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fallbackInvestigation(req), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fallbackInvestigation(req), nil
	}

	var res InvestigationResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fallbackInvestigation(req), nil
	}
	return &res, nil
}

func fallbackInvestigation(req *InvestigationRequest) *InvestigationResponse {
	action := req.Decision
	if action == "" {
		action = "ALLOW"
		if req.RiskScore >= 80 {
			action = "BLOCK"
		} else if req.RiskScore >= 60 {
			action = "CHALLENGE"
		}
	}
	riskLevel := "LOW"
	switch action {
	case "BLOCK":
		riskLevel = "CRITICAL"
	case "CHALLENGE":
		riskLevel = "HIGH"
	}

	return &InvestigationResponse{
		InvestigationID:   fmt.Sprintf("INV-%s", req.TransactionID),
		TransactionID:     req.TransactionID,
		RiskLevel:         riskLevel,
		Summary:           fmt.Sprintf("Investigation for transaction %s of $%.2f (Risk Score: %d/100).", req.TransactionID, req.Amount, req.RiskScore),
		Evidence:          []EvidenceItem{{Category: "SYSTEM_SIGNAL", Fact: fmt.Sprintf("Risk Score evaluated at %d/100", req.RiskScore), Severity: riskLevel}},
		SimilarCases:      []SimilarCase{},
		RecommendedAction: action,
		Confidence:        0.85,
		LLMModel:          "rule-based-fallback",
		RiskScore:         req.RiskScore,
		FraudProbability:  req.FraudProbability,
		AnomalyScore:      req.AnomalyScore,
		PolicyDecision:    req.Decision,
		TransactionStatus: req.Status,
	}
}
