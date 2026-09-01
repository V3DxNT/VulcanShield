package v1

import (
	"net/http"
	"strconv"

	"github.com/vulcanshield/backend/internal/db/repository"
	"github.com/vulcanshield/backend/internal/graph"
	"github.com/vulcanshield/backend/internal/models"
)

func setNoStoreHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

type GraphHandlers struct {
	GraphRepo repository.GraphRepository
	Engine    *graph.Engine
}


func (h *GraphHandlers) ListRelationships(w http.ResponseWriter, r *http.Request) {
	setNoStoreHeaders(w)

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

	if userID := r.URL.Query().Get("user_id"); userID != "" {
		filtered := make([]models.FraudRelationship, 0, len(rels))
		for _, rel := range rels {
			if rel.SourceID == userID || rel.TargetID == userID {
				filtered = append(filtered, rel)
			}
		}
		rels = filtered
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data":  rels,
		"total": len(rels),
	})
}


func (h *GraphHandlers) GetNeighbors(w http.ResponseWriter, r *http.Request) {
	setNoStoreHeaders(w)

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
