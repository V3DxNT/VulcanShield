package models

import "time"

// FraudRelationship represents an edge in the relational fraud graph.
// Matches the Phase 2 schema `fraud_relationships` table.
type FraudRelationship struct {
	RelationshipID   string    `json:"relationship_id"`
	SourceType       string    `json:"source_type"`
	SourceID         string    `json:"source_id"`
	TargetType       string    `json:"target_type"`
	TargetID         string    `json:"target_id"`
	RelationshipType string    `json:"relationship_type"`
	Weight           float64   `json:"weight"`
	FraudLinked      bool      `json:"fraud_linked"`
	CreatedAt        time.Time `json:"created_at"`
}

// GraphFeatures contains graph-derived features for a user or transaction.
type GraphFeatures struct {
	SharedDeviceAccounts int  `json:"shared_device_accounts"`
	SharedIPAccounts     int  `json:"shared_ip_accounts"`
	FraudNeighborCount   int  `json:"fraud_neighbor_count"`
	HasFraudLink         bool `json:"has_fraud_link"`
}

type GraphNeighborsResponse struct {
	UserID        string              `json:"user_id"`
	Relationships []FraudRelationship `json:"relationships"`
	Features      GraphFeatures       `json:"features"`
}
