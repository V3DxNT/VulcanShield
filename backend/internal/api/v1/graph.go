package v1

import (
	"net/http"
	"strconv"

	"github.com/vulcanshield/backend/internal/db/repository"
	"github.com/vulcanshield/backend/internal/graph"
	"github.com/vulcanshield/backend/internal/models"
)

type GraphHandlers struct {
	GraphRepo repository.GraphRepository
	Engine    *graph.Engine
}

// ListRelationships handles GET /api/v1/graph/relationships.
func (h *GraphHandlers) ListRelationships(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	rels, err := h.GraphRepo.ListRelationships(r.Context(), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "QUERY_FAILED", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data":  rels,
		"total": len(rels),
	})
}

// GetNeighbors handles GET /api/v1/graph/neighbors/{user_id}.
func (h *GraphHandlers) GetNeighbors(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "MISSING_USER_ID", "user_id is required")
		return
	}

	rels, err := h.GraphRepo.GetNeighbors(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "QUERY_FAILED", err.Error())
		return
	}

	features, _ := h.Engine.ExtractFeatures(r.Context(), userID, "", "")

	writeJSON(w, http.StatusOK, models.GraphNeighborsResponse{
		UserID:        userID,
		Relationships: rels,
		Features:      *features,
	})
}
