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

// TransactionHandlers handles transaction query endpoints.
type TransactionHandlers struct {
	TxRepo repository.TransactionRepository
}

// GetByID handles GET /api/v1/transactions/{id}.
// Returns 404 if the transaction does not exist.
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
	writeJSON(w, http.StatusOK, tx)
}

// List handles GET /api/v1/transactions.
// Supports query parameters: user_id (optional), limit (default 20, max 100).
// Enforces sane pagination per implementation constraints.
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

	writeJSON(w, http.StatusOK, map[string]any{
		"data":  txs,
		"total": len(txs),
	})
}
