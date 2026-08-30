package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/vulcanshield/backend/internal/challenge"
	"github.com/vulcanshield/backend/internal/db/repository"
	"github.com/vulcanshield/backend/internal/models"
)

type ChallengeHandlers struct {
	ChallengeRepo repository.ChallengeRepository
	TxRepo        repository.TransactionRepository
	PolicyRepo    repository.PolicyRepository
	OTPService    *challenge.Service
}

// GetByID handles GET /api/v1/challenges/{id}.
func (h *ChallengeHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "challenge id is required")
		return
	}

	c, err := h.ChallengeRepo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "challenge not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "QUERY_FAILED", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, c)
}

// Verify handles POST /api/v1/challenges/{id}/verify.
// Validates 6-digit OTP against 60s TTL / SHA-256 hash.
func (h *ChallengeHandlers) Verify(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "MISSING_ID", "challenge id is required")
		return
	}

	var req models.OTPVerifyRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	c, err := h.ChallengeRepo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "challenge not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "QUERY_FAILED", err.Error())
		return
	}

	if c.Status != models.ChallengePending {
		writeError(w, http.StatusBadRequest, "INVALID_STATE", "challenge is no longer pending")
		return
	}

	now := time.Now().UTC()
	if now.After(c.ExpiresAt) {
		c.Status = models.ChallengeExpired
		_ = h.ChallengeRepo.Update(r.Context(), c)
		_ = h.TxRepo.UpdateStatus(r.Context(), c.TransactionID, models.StatusBlocked)

		writeJSON(w, http.StatusOK, models.OTPVerifyResponse{
			ChallengeID:   c.ChallengeID,
			TransactionID: c.TransactionID,
			Status:        models.ChallengeExpired,
			Message:       "OTP challenge expired (>60s)",
			FinalStatus:   models.StatusBlocked,
		})
		return
	}

	c.Attempts++
	success := h.OTPService.Verify(c, req.OTPCode)

	if success {
		c.Status = models.ChallengeVerified
		c.VerifiedAt = &now
		_ = h.ChallengeRepo.Update(r.Context(), c)
		_ = h.TxRepo.UpdateStatus(r.Context(), c.TransactionID, models.StatusApproved)

		// Insert policy decision update ALLOW
		if h.PolicyRepo != nil {
			_ = h.PolicyRepo.Create(r.Context(), &models.PolicyDecision{
				DecisionID:     "PD-VERIFIED-" + c.TransactionID,
				TransactionID:  c.TransactionID,
				Decision:       models.DecisionAllow,
				Reason:         "Step-up OTP successfully verified",
				RulesTriggered: []string{"RULE_OTP_VERIFIED"},
				CreatedAt:      now,
			})
		}

		writeJSON(w, http.StatusOK, models.OTPVerifyResponse{
			ChallengeID:   c.ChallengeID,
			TransactionID: c.TransactionID,
			Status:        models.ChallengeVerified,
			Message:       "OTP verified successfully — transaction approved",
			FinalStatus:   models.StatusApproved,
		})
		return
	}

	// Invalid OTP
	if c.Attempts >= c.MaxAttempts {
		c.Status = models.ChallengeFailed
		_ = h.ChallengeRepo.Update(r.Context(), c)
		_ = h.TxRepo.UpdateStatus(r.Context(), c.TransactionID, models.StatusBlocked)

		writeJSON(w, http.StatusOK, models.OTPVerifyResponse{
			ChallengeID:   c.ChallengeID,
			TransactionID: c.TransactionID,
			Status:        models.ChallengeFailed,
			Message:       "Max attempts exceeded — transaction blocked",
			FinalStatus:   models.StatusBlocked,
		})
		return
	}

	_ = h.ChallengeRepo.Update(r.Context(), c)
	writeJSON(w, http.StatusBadRequest, models.OTPVerifyResponse{
		ChallengeID:   c.ChallengeID,
		TransactionID: c.TransactionID,
		Status:        models.ChallengePending,
		Message:       "Incorrect OTP code",
		FinalStatus:   models.StatusChallenged,
	})
}
