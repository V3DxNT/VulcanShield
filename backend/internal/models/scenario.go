package models

import "time"


type ScenarioType string

const (
	ScenarioNormal          ScenarioType = "normal"
	ScenarioVelocityAttack  ScenarioType = "velocity_attack"
	ScenarioAccountTakeover ScenarioType = "account_takeover"
	ScenarioDeviceFarm      ScenarioType = "device_farm"
	ScenarioIPAbuse         ScenarioType = "ip_abuse"
	ScenarioAmountAnomaly   ScenarioType = "amount_anomaly"
)


type ScenarioStatus string

const (
	ScenarioRunning   ScenarioStatus = "RUNNING"
	ScenarioCompleted ScenarioStatus = "COMPLETED"
	ScenarioStopped   ScenarioStatus = "STOPPED"
)


type ScenarioStartRequest struct {
	Scenario     ScenarioType `json:"scenario"`
	Transactions int          `json:"transactions"` 
	IntervalMS   int          `json:"interval_ms"`  
	Seed         int64        `json:"seed"`         
	CustomerID   string       `json:"customer_id"`  
}



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



type EntityPool struct {
	Users           []UserProfile
	UserDeviceMap   map[string]string
	UserIPMap       map[string]string
	DeviceIDs       []string
	IPAddresses     []string
	MerchantIDs     []string
}


type UserProfile struct {
	UserID             string
	TypicalMinAmount   float64
	TypicalMaxAmount   float64
	TrustScore         int
	ChallengeThreshold int
	BlockThreshold     int
}
