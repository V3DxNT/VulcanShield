package models

import "time"



type TransactionStatus string

const (
	StatusPending    TransactionStatus = "PENDING"
	StatusApproved   TransactionStatus = "APPROVED"
	StatusChallenged TransactionStatus = "CHALLENGED"
	StatusBlocked    TransactionStatus = "BLOCKED"
	StatusCancelled  TransactionStatus = "CANCELLED"
)



type Transaction struct {
	TransactionID string            `json:"transaction_id"`
	UserID        string            `json:"user_id"`
	DeviceID      string            `json:"device_id"`
	IPAddress     string            `json:"ip_address"`
	MerchantID    string            `json:"merchant_id"`
	Amount        float64           `json:"amount"`
	Currency      string            `json:"currency"`
	Channel       string            `json:"channel"`
	Status        TransactionStatus `json:"status"`
	Timestamp     time.Time         `json:"timestamp"`
	CreatedAt     time.Time         `json:"created_at"`
}




type User struct {
	UserID             string    `json:"user_id"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	RiskTolerance      string    `json:"risk_tolerance"`
	ChallengeThreshold int       `json:"challenge_threshold"`
	BlockThreshold     int       `json:"block_threshold"`
	AccountAgeDays     int       `json:"account_age_days"`
	TrustScore         int       `json:"trust_score"`
	TypicalMinAmount   float64   `json:"typical_min_amount"`
	TypicalMaxAmount   float64   `json:"typical_max_amount"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
