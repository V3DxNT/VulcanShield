package models

import "time"

type ChallengeStatus string

const (
	ChallengePending  ChallengeStatus = "PENDING"
	ChallengeVerified ChallengeStatus = "VERIFIED"
	ChallengeExpired  ChallengeStatus = "EXPIRED"
	ChallengeFailed   ChallengeStatus = "FAILED"
)



type OTPChallenge struct {
	ChallengeID   string          `json:"challenge_id"`
	TransactionID string          `json:"transaction_id"`
	OTPCodeHash   string          `json:"-"`                  
	DemoOTP       string          `json:"demo_otp,omitempty"` 
	Status        ChallengeStatus `json:"status"`
	Attempts      int             `json:"attempts"`
	MaxAttempts   int             `json:"max_attempts"`
	CreatedAt     time.Time       `json:"created_at"`
	ExpiresAt     time.Time       `json:"expires_at"`
	VerifiedAt    *time.Time      `json:"verified_at,omitempty"`
}

type OTPVerifyRequest struct {
	OTPCode string `json:"otp_code"`
}

type OTPVerifyResponse struct {
	ChallengeID   string            `json:"challenge_id"`
	TransactionID string            `json:"transaction_id"`
	Status        ChallengeStatus   `json:"status"`
	Message       string            `json:"message"`
	FinalStatus   TransactionStatus `json:"final_status"`
}
