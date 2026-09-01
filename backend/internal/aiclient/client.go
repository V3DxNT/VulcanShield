package aiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
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

type RetrievalTraceItem struct {
	Source         string   `json:"source"`
	Query          string   `json:"query"`
	MatchedDocs    []string `json:"matched_documents,omitempty"`
	RelevanceScore float64  `json:"relevance_score"`
}

type InvestigationResponse struct {
	InvestigationID   string               `json:"investigation_id"`
	TransactionID     string               `json:"transaction_id"`
	RiskLevel         string               `json:"risk_level"`
	Summary           string               `json:"summary"`
	Evidence          []EvidenceItem       `json:"evidence"`
	SimilarCases      []SimilarCase        `json:"similar_cases"`
	RecommendedAction string               `json:"recommended_action"`
	Confidence        float64              `json:"confidence"`
	LLMModel          string               `json:"llm_model"`
	LLMPrompt         string               `json:"llm_prompt,omitempty"`
	RiskScore         int                  `json:"risk_score"`
	FraudProbability  float64              `json:"fraud_probability"`
	AnomalyScore      float64              `json:"anomaly_score"`
	PolicyDecision    string               `json:"policy_decision"`
	TransactionStatus string               `json:"transaction_status"`
	InitialRiskScore  int                  `json:"initial_risk_score,omitempty"`
	FinalRiskScore    int                  `json:"final_risk_score,omitempty"`
	RetrievalTrace    []RetrievalTraceItem `json:"retrieval_trace,omitempty"`
	ReasoningTrace    []string             `json:"reasoning_trace,omitempty"`
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
		log.Printf("AI service unreachable at %s; using fallback: %v", url, err)
		return fallbackInvestigation(req), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var body map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&body); err == nil {
			log.Printf("AI service returned %d for %s; using fallback. body=%v", resp.StatusCode, url, body)
		} else {
			log.Printf("AI service returned %d for %s; using fallback.", resp.StatusCode, url)
		}
		return fallbackInvestigation(req), nil
	}

	var res InvestigationResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Printf("AI service returned invalid JSON for %s; using fallback: %v", url, err)
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

	initialRisk := req.RiskScore
	finalRisk := req.RiskScore
	if action == "BLOCK" {
		finalRisk = min(initialRisk+10, 100)
	} else if action == "ALLOW" {
		finalRisk = max(initialRisk-8, 0)
	} else if action == "CHALLENGE" {
		finalRisk = max(initialRisk, 60)
	}

	return &InvestigationResponse{
		InvestigationID:   fmt.Sprintf("INV-%s", req.TransactionID),
		TransactionID:     req.TransactionID,
		RiskLevel:         riskLevel,
		Summary:           fmt.Sprintf("Investigation for transaction %s of $%.2f (Initial risk: %d, Final risk: %d).", req.TransactionID, req.Amount, initialRisk, finalRisk),
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
		InitialRiskScore:  initialRisk,
		FinalRiskScore:    finalRisk,
		RetrievalTrace: []RetrievalTraceItem{{
			Source:         "system",
			Query:          "risk evaluation fallback",
			MatchedDocs:    []string{"local policy engine", "deterministic decision fallback"},
			RelevanceScore: 0.86,
		}},
		ReasoningTrace: []string{
			fmt.Sprintf("initial risk = %d/100", initialRisk),
			fmt.Sprintf("policy action = %s", action),
			fmt.Sprintf("final risk = %d/100", finalRisk),
		},
	}
}
