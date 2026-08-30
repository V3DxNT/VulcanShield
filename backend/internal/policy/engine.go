package policy

import (
	"fmt"
	"time"

	"github.com/vulcanshield/backend/internal/models"
)

// Engine evaluates standardized risk scores against per-user authoritative policy thresholds.
type Engine struct{}

// NewEngine creates a new Policy Engine.
func NewEngine() *Engine {
	return &Engine{}
}

// Evaluate determines the authoritative payment authorization decision.
// Default thresholds: challenge_threshold=65, block_threshold=85.
func (e *Engine) Evaluate(
	tx *models.Transaction,
	ra *models.RiskAssessment,
	user *models.UserProfile,
	userThresholds *models.User,
) (*models.PolicyDecision, models.TransactionStatus) {
	challengeThreshold := 65
	blockThreshold := 85

	if userThresholds != nil {
		if userThresholds.ChallengeThreshold > 0 {
			challengeThreshold = userThresholds.ChallengeThreshold
		}
		if userThresholds.BlockThreshold > 0 {
			blockThreshold = userThresholds.BlockThreshold
		}
	}

	riskScore := 0
	if ra != nil {
		riskScore = ra.RiskScore
	}

	var decision models.PolicyDecisionType
	var status models.TransactionStatus
	var reason string
	rules := []string{}

	if riskScore >= blockThreshold {
		decision = models.DecisionBlock
		status = models.StatusBlocked
		reason = fmt.Sprintf("Risk score %d exceeded user block threshold %d", riskScore, blockThreshold)
		rules = append(rules, "RULE_BLOCK_THRESHOLD_EXCEEDED")
	} else if riskScore >= challengeThreshold {
		decision = models.DecisionChallenge
		status = models.StatusChallenged
		reason = fmt.Sprintf("Risk score %d exceeded user challenge threshold %d", riskScore, challengeThreshold)
		rules = append(rules, "RULE_CHALLENGE_THRESHOLD_EXCEEDED")
	} else {
		decision = models.DecisionAllow
		status = models.StatusApproved
		reason = fmt.Sprintf("Risk score %d within acceptable threshold (%d)", riskScore, challengeThreshold)
		rules = append(rules, "RULE_ALLOW_PASSTHROUGH")
	}

	now := time.Now().UTC()
	return &models.PolicyDecision{
		DecisionID:     fmt.Sprintf("PD-%s", tx.TransactionID),
		TransactionID:  tx.TransactionID,
		Decision:       decision,
		Reason:         reason,
		RulesTriggered: rules,
		CreatedAt:      now,
	}, status
}
