package graph

import (
	"context"

	"github.com/vulcanshield/backend/internal/db/repository"
	"github.com/vulcanshield/backend/internal/models"
)

type Engine struct {
	graphRepo repository.GraphRepository
}

func NewEngine(graphRepo repository.GraphRepository) *Engine {
	return &Engine{graphRepo: graphRepo}
}

func (e *Engine) ExtractFeatures(ctx context.Context, userID, deviceID, ipAddress string) (*models.GraphFeatures, error) {
	if e.graphRepo == nil {
		return &models.GraphFeatures{}, nil
	}

	rels, err := e.graphRepo.GetNeighbors(ctx, userID)
	if err != nil {
		return &models.GraphFeatures{}, nil // Fallback
	}

	gf := &models.GraphFeatures{}
	for _, r := range rels {
		if r.RelationshipType == "SHARED_DEVICE" {
			gf.SharedDeviceAccounts++
		}
		if r.RelationshipType == "SHARED_IP" {
			gf.SharedIPAccounts++
		}
		if r.FraudLinked {
			gf.FraudNeighborCount++
			gf.HasFraudLink = true
		}
	}
	return gf, nil
}
