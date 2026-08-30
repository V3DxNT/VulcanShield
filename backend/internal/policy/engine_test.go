package policy

import (
	"testing"

	"github.com/vulcanshield/backend/internal/models"
)

func TestEvaluate_AllowChallengeBlock(t *testing.T) {
	e := NewEngine()
	tx := &models.Transaction{TransactionID: "TX-1"}

	pd, status := e.Evaluate(tx, &models.RiskAssessment{RiskScore: 40}, nil, nil)
	if pd.Decision != models.DecisionAllow || status != models.StatusApproved {
		t.Fatalf("score 40: got %s / %s", pd.Decision, status)
	}
	if pd.RiskScore != 40 || pd.ChallengeThreshold != 65 || pd.BlockThreshold != 85 {
		t.Fatalf("defaults not recorded: %+v", pd)
	}

	pd, status = e.Evaluate(tx, &models.RiskAssessment{RiskScore: 70}, nil, nil)
	if pd.Decision != models.DecisionChallenge || status != models.StatusChallenged {
		t.Fatalf("score 70: got %s / %s", pd.Decision, status)
	}

	pd, status = e.Evaluate(tx, &models.RiskAssessment{RiskScore: 90}, nil, nil)
	if pd.Decision != models.DecisionBlock || status != models.StatusBlocked {
		t.Fatalf("score 90: got %s / %s", pd.Decision, status)
	}
}

func TestEvaluate_PerUserThresholds(t *testing.T) {
	e := NewEngine()
	tx := &models.Transaction{TransactionID: "TX-2"}
	user := &models.UserProfile{UserID: "C1003", ChallengeThreshold: 75, BlockThreshold: 90}

	pd, status := e.Evaluate(tx, &models.RiskAssessment{RiskScore: 70}, user, nil)
	if pd.Decision != models.DecisionAllow || status != models.StatusApproved {
		t.Fatalf("high-tolerance user at 70 should ALLOW, got %s / %s", pd.Decision, status)
	}
	if pd.ChallengeThreshold != 75 || pd.BlockThreshold != 90 {
		t.Fatalf("user thresholds not recorded: %+v", pd)
	}

	strict := &models.User{ChallengeThreshold: 55, BlockThreshold: 80}
	pd, status = e.Evaluate(tx, &models.RiskAssessment{RiskScore: 70}, nil, strict)
	if pd.Decision != models.DecisionChallenge || status != models.StatusChallenged {
		t.Fatalf("low-tolerance user at 70 should CHALLENGE, got %s / %s", pd.Decision, status)
	}
}
