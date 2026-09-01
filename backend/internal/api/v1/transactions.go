package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/vulcanshield/backend/internal/db/repository"
	"github.com/vulcanshield/backend/internal/models"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)


type TransactionHandlers struct {
	TxRepo        repository.TransactionRepository
	RiskRepo      repository.RiskRepository
	PolicyRepo    repository.PolicyRepository
	ChallengeRepo repository.ChallengeRepository
}



func (h *TransactionHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "transaction id is required")
		return
	}

	tx, err := h.TxRepo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "transaction not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "QUERY_FAILED", err.Error())
		return
	}
	view := transactionView{Transaction: *tx}
	if h.RiskRepo != nil {
		if ra, err := h.RiskRepo.GetByTransactionID(r.Context(), id); err == nil && ra != nil {
			view.RiskScore = ra.RiskScore
			view.FraudProbability = ra.FraudProbability
			view.AnomalyScore = ra.AnomalyScore
		}
	}
	if h.PolicyRepo != nil {
		if pd, err := h.PolicyRepo.GetByTransactionID(r.Context(), id); err == nil && pd != nil {
			view.Decision = string(pd.Decision)
			view.DecisionReason = pd.Reason
		}
	}
	if h.ChallengeRepo != nil {
		if challenge, err := h.ChallengeRepo.GetByTransactionID(r.Context(), id); err == nil && challenge != nil {
			view.ChallengeStatus = string(challenge.Status)
		}
	}
	writeJSON(w, http.StatusOK, view)
}

type transactionView struct {
	models.Transaction
	RiskScore        int     `json:"risk_score"`
	FraudProbability float64 `json:"fraud_probability"`
	AnomalyScore     float64 `json:"anomaly_score"`
	Decision         string  `json:"decision,omitempty"`
	DecisionReason   string  `json:"decision_reason,omitempty"`
	ChallengeStatus  string  `json:"challenge_status,omitempty"`
}




func (h *TransactionHandlers) List(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	limit := defaultLimit
	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	var txs []*models.Transaction
	var err error

	if userID != "" {
		txs, err = h.TxRepo.ListByUser(r.Context(), userID, limit)
	} else {
		txs, err = h.TxRepo.ListRecent(r.Context(), limit)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "QUERY_FAILED", err.Error())
		return
	}

	data := make([]transactionView, 0, len(txs))
	for _, tx := range txs {
		view := transactionView{Transaction: *tx}
		if h.RiskRepo != nil {
			if ra, err := h.RiskRepo.GetByTransactionID(r.Context(), tx.TransactionID); err == nil && ra != nil {
				view.RiskScore = ra.RiskScore
				view.FraudProbability = ra.FraudProbability
				view.AnomalyScore = ra.AnomalyScore
			}
		}
		if h.PolicyRepo != nil {
			if pd, err := h.PolicyRepo.GetByTransactionID(r.Context(), tx.TransactionID); err == nil && pd != nil {
				view.Decision = string(pd.Decision)
				view.DecisionReason = pd.Reason
			}
		}
		if h.ChallengeRepo != nil {
			if challenge, err := h.ChallengeRepo.GetByTransactionID(r.Context(), tx.TransactionID); err == nil && challenge != nil {
				view.ChallengeStatus = string(challenge.Status)
			}
		}
		data = append(data, view)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data":  data,
		"total": len(data),
	})
}
