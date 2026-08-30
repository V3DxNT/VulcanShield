package aiclient

import (
	"encoding/json"
	"testing"
)

func TestInvestigationResponseIncludesRetrievalAndRiskTrace(t *testing.T) {
	payload := []byte(`{
		"investigation_id": "INV-TX-1001",
		"transaction_id": "TX-1001",
		"risk_level": "HIGH",
		"summary": "Summary",
		"evidence": [],
		"similar_cases": [],
		"recommended_action": "CHALLENGE",
		"confidence": 0.9,
		"llm_model": "qwen2.5:7b-instruct",
		"initial_risk_score": 62,
		"final_risk_score": 73,
		"retrieval_trace": [{"source": "customer_history", "query": "recent fraud pattern"}],
		"reasoning_trace": ["initial risk = 62", "policy challenge applied", "final risk = 73"]
	}`)

	var out InvestigationResponse
	if err := json.Unmarshal(payload, &out); err != nil {
		t.Fatalf("json decode failed: %v", err)
	}

	if out.InitialRiskScore != 62 {
		t.Fatalf("expected initial_risk_score 62, got %d", out.InitialRiskScore)
	}
	if out.FinalRiskScore != 73 {
		t.Fatalf("expected final_risk_score 73, got %d", out.FinalRiskScore)
	}
	if len(out.RetrievalTrace) != 1 {
		t.Fatalf("expected 1 retrieval trace item, got %d", len(out.RetrievalTrace))
	}
	if len(out.ReasoningTrace) != 3 {
		t.Fatalf("expected 3 reasoning trace items, got %d", len(out.ReasoningTrace))
	}
}
