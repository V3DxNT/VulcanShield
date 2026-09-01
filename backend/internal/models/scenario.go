package models

import "time"

// ScenarioType identifies which scenario pattern the generator will run.
type ScenarioType string

const (
	ScenarioNormal          ScenarioType = "normal"
	ScenarioVelocityAttack  ScenarioType = "velocity_attack"
	ScenarioAccountTakeover ScenarioType = "account_takeover"
	ScenarioDeviceFarm      ScenarioType = "device_farm"
	ScenarioIPAbuse         ScenarioType = "ip_abuse"
	ScenarioAmountAnomaly   ScenarioType = "amount_anomaly"
)

// ScenarioStatus mirrors the Phase 2 schema constraint on scenarios.status.
type ScenarioStatus string

const (
	ScenarioRunning   ScenarioStatus = "RUNNING"
	ScenarioCompleted ScenarioStatus = "COMPLETED"
	ScenarioStopped   ScenarioStatus = "STOPPED"
)

// ScenarioStartRequest is the JSON body for POST /api/v1/scenarios/start.
type ScenarioStartRequest struct {
	Scenario     ScenarioType `json:"scenario"`
	Transactions int          `json:"transactions"` // Default: 100
	IntervalMS   int          `json:"interval_ms"`  // Default: 1000
	Seed         int64        `json:"seed"`         // Default: 42
	CustomerID   string       `json:"customer_id"`  // Optional: target a specific user
}

// ScenarioRun is the in-memory and HTTP response representation of an active
// or completed scenario execution. Mirrors the Phase 2 scenarios table.
type ScenarioRun struct {
	ScenarioID       string         `json:"scenario_id"`
	ScenarioType     ScenarioType   `json:"scenario_type"`
	TransactionCount int            `json:"transaction_count"`
	GeneratedCount   int            `json:"generated_count"`
	Seed             int64          `json:"seed"`
	Status           ScenarioStatus `json:"status"`
	StartedAt        time.Time      `json:"started_at"`
	EndedAt          *time.Time     `json:"ended_at,omitempty"`
}

// EntityPool holds the seeded synthetic entity IDs available for generation.
// Loaded once at startup from PostgreSQL seed data.
type EntityPool struct {
	Users           []UserProfile
	UserDeviceMap   map[string]string
	UserIPMap       map[string]string
	DeviceIDs       []string
	IPAddresses     []string
	MerchantIDs     []string
}

// UserProfile is a minimal projection used by the generator.
type UserProfile struct {
	UserID             string
	TypicalMinAmount   float64
	TypicalMaxAmount   float64
	TrustScore         int
	ChallengeThreshold int
	BlockThreshold     int
}
