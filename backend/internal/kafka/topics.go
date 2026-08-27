package kafka

// Topic names must match the values in .env.example and PROJECT_SPEC.md §7.
const (
	// TopicTransactionCreated is published when a new transaction is generated.
	TopicTransactionCreated = "transaction.created"

	// TopicRiskEvaluated is published after ML risk assessment is complete.
	TopicRiskEvaluated = "risk.evaluated"

	// TopicTransactionDecisioned is published after the policy engine decides.
	TopicTransactionDecisioned = "transaction.decisioned"

	// TopicAIInvestigationRequested is published to trigger an AI investigation.
	TopicAIInvestigationRequested = "ai.investigation.requested"

	// TopicAIInvestigationCompleted is published when an AI investigation finishes.
	TopicAIInvestigationCompleted = "ai.investigation.completed"
)
