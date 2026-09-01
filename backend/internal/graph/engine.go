package graph

import (
	"context"
	"fmt"
	"time"

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
		return &models.GraphFeatures{}, nil 
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



func (e *Engine) RecordTransactionEdges(ctx context.Context, tx *models.Transaction, isEmulator, isVPN bool) error {
	if e.graphRepo == nil || tx == nil {
		return nil
	}
	now := time.Now().UTC()
	edges := []models.FraudRelationship{
		{
			RelationshipID:   fmt.Sprintf("REL-USED-%s-%s", tx.UserID, tx.DeviceID),
			SourceType:       "USER",
			SourceID:         tx.UserID,
			TargetType:       "DEVICE",
			TargetID:         tx.DeviceID,
			RelationshipType: "USED",
			Weight:           0.95,
			FraudLinked:      isEmulator,
			CreatedAt:        now,
		},
		{
			RelationshipID:   fmt.Sprintf("REL-CONN-%s-%s", tx.UserID, tx.IPAddress),
			SourceType:       "USER",
			SourceID:         tx.UserID,
			TargetType:       "IP",
			TargetID:         tx.IPAddress,
			RelationshipType: "CONNECTED",
			Weight:           0.9,
			FraudLinked:      isVPN,
			CreatedAt:        now,
		},
		{
			RelationshipID:   fmt.Sprintf("REL-TX-%s-%s", tx.UserID, tx.MerchantID),
			SourceType:       "USER",
			SourceID:         tx.UserID,
			TargetType:       "MERCHANT",
			TargetID:         tx.MerchantID,
			RelationshipType: "TRANSACTED_WITH",
			Weight:           1.0,
			FraudLinked:      false,
			CreatedAt:        now,
		},
	}
	for i := range edges {
		if err := e.graphRepo.CreateRelationship(ctx, &edges[i]); err != nil {
			return err
		}
	}
	return nil
}
