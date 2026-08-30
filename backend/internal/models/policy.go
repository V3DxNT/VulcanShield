package models

import "time"

type PolicyDecisionType string

const (
	DecisionAllow     PolicyDecisionType = "ALLOW"
	DecisionChallenge PolicyDecisionType = "CHALLENGE"
	DecisionBlock     PolicyDecisionType = "BLOCK"
)

// PolicyDecision represents the authoritative authorization decision made by the Policy Engine.
// Matches the Phase 2 schema `policy_decisions` table.
type PolicyDecision struct {
	DecisionID     string             `json:"decision_id"`
	TransactionID  string             `json:"transaction_id"`
	Decision       PolicyDecisionType `json:"decision"`
	Reason         string             `json:"reason"`
	RulesTriggered []string           `json:"rules_triggered"`
	CreatedAt      time.Time          `json:"created_at"`
}
