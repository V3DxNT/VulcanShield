package v1

import (
	"errors"
	"net/http"

	"github.com/vulcanshield/backend/internal/aiclient"
	"github.com/vulcanshield/backend/internal/db/repository"
	"github.com/vulcanshield/backend/internal/models"
)

type InvestigationHandlers struct {
	TxRepo     repository.TransactionRepository
	RiskRepo   repository.RiskRepository
	PolicyRepo repository.PolicyRepository
	AIClient   *aiclient.Client
}

// Investigate handles GET /api/v1/investigations/{transaction_id}.
func (h *InvestigationHandlers) Investigate(w http.ResponseWriter, r *http.Request) {
	txID := r.PathValue("transaction_id")
	if txID == "" {
		writeError(w, http.StatusBadRequest, "MISSING_TX_ID", "transaction_id is required")
		return
	}

	tx, err := h.TxRepo.GetByID(r.Context(), txID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "transaction not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "QUERY_FAILED", err.Error())
		return
	}

	riskScore := 0
	fraudProb := 0.0
	anomalyScore := 0.0

	if h.RiskRepo != nil {
		ra, err := h.RiskRepo.GetByTransactionID(r.Context(), txID)
		if err == nil && ra != nil {
			riskScore = ra.RiskScore
			fraudProb = ra.FraudProbability
			anomalyScore = ra.AnomalyScore
		}
	}

	decision := ""
	switch tx.Status {
	case models.StatusApproved:
		decision = "ALLOW"
	case models.StatusChallenged:
		decision = "CHALLENGE"
	case models.StatusBlocked:
		decision = "BLOCK"
	}
	if h.PolicyRepo != nil {
		if pd, err := h.PolicyRepo.GetByTransactionID(r.Context(), txID); err == nil && pd != nil {
			decision = string(pd.Decision)
		}
	}

	req := &aiclient.InvestigationRequest{
		TransactionID:    tx.TransactionID,
		UserID:           tx.UserID,
		DeviceID:         tx.DeviceID,
		IPAddress:        tx.IPAddress,
		Amount:           tx.Amount,
		RiskScore:        riskScore,
		FraudProbability: fraudProb,
		AnomalyScore:     anomalyScore,
		Status:           string(tx.Status),
		Decision:         decision,
	}

	inv, err := h.AIClient.Investigate(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AI_INVESTIGATION_FAILED", err.Error())
		return
	}
	inv.RiskScore = riskScore
	inv.FraudProbability = fraudProb
	inv.AnomalyScore = anomalyScore
	inv.PolicyDecision = decision
	inv.TransactionStatus = string(tx.Status)
	inv.RecommendedAction = decision

	writeJSON(w, http.StatusOK, inv)
}
