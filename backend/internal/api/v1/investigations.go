package v1

import (
	"errors"
	"net/http"

	"github.com/vulcanshield/backend/internal/aiclient"
	"github.com/vulcanshield/backend/internal/db/repository"
)

type InvestigationHandlers struct {
	TxRepo   repository.TransactionRepository
	RiskRepo repository.RiskRepository
	AIClient *aiclient.Client
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

	riskScore := 50
	fraudProb := 0.5
	anomalyScore := 0.5

	if h.RiskRepo != nil {
		ra, err := h.RiskRepo.GetByTransactionID(r.Context(), txID)
		if err == nil && ra != nil {
			riskScore = ra.RiskScore
			fraudProb = ra.FraudProbability
			anomalyScore = ra.AnomalyScore
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
	}

	inv, err := h.AIClient.Investigate(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AI_INVESTIGATION_FAILED", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, inv)
}
