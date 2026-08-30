package challenge

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vulcanshield/backend/internal/models"
)

// Service manages OTP challenge generation, Redis state, and verification.
type Service struct {
	redisClient *redis.Client
}

// NewService returns a new OTP Challenge Service.
func NewService(redisClient *redis.Client) *Service {
	return &Service{redisClient: redisClient}
}

// GenerateChallenge creates a new 6-digit OTP challenge expiring in 60s.
// Plaintext OTP is stored ONLY in short-lived Redis key for demo UI display;
// ONLY the SHA-256 hash is persisted in PostgreSQL.
func (s *Service) GenerateChallenge(ctx context.Context, txID string) (*models.OTPChallenge, string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return nil, "", fmt.Errorf("generating otp random number: %w", err)
	}
	otpCode := fmt.Sprintf("%06d", n.Int64()+100000)
	otpHash := HashOTP(otpCode)

	now := time.Now().UTC()
	expiresAt := now.Add(60 * time.Second)
	challengeID := fmt.Sprintf("CH-%s", txID)

	challenge := &models.OTPChallenge{
		ChallengeID:   challengeID,
		TransactionID: txID,
		OTPCodeHash:   otpHash,
		DemoOTP:       otpCode,
		Status:        models.ChallengePending,
		Attempts:      0,
		MaxAttempts:   3,
		CreatedAt:     now,
		ExpiresAt:     expiresAt,
	}

	// Store demo OTP in Redis with 60s TTL
	if s.redisClient != nil {
		redisKey := fmt.Sprintf("otp:%s", challengeID)
		_ = s.redisClient.Set(ctx, redisKey, otpCode, 60*time.Second).Err()
	}

	return challenge, otpCode, nil
}

// VerifyChecks whether the submitted code matches the SHA-256 hash within the 60s window.
func (s *Service) Verify(challenge *models.OTPChallenge, inputCode string) bool {
	if time.Now().UTC().After(challenge.ExpiresAt) {
		return false
	}
	return HashOTP(inputCode) == challenge.OTPCodeHash
}

// HashOTP returns the hex SHA-256 hash of plaintext OTP.
func HashOTP(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}
