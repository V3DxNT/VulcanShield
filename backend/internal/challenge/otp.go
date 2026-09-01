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


type Service struct {
	redisClient *redis.Client
}


func NewService(redisClient *redis.Client) *Service {
	return &Service{redisClient: redisClient}
}




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

	
	if s.redisClient != nil {
		redisKey := fmt.Sprintf("otp:%s", challengeID)
		_ = s.redisClient.Set(ctx, redisKey, otpCode, 60*time.Second).Err()
	}

	return challenge, otpCode, nil
}


func (s *Service) DemoOTP(ctx context.Context, challengeID string) string {
	if s.redisClient == nil {
		return ""
	}
	val, err := s.redisClient.Get(ctx, fmt.Sprintf("otp:%s", challengeID)).Result()
	if err != nil {
		return ""
	}
	return val
}


func (s *Service) Verify(challenge *models.OTPChallenge, inputCode string) bool {
	if time.Now().UTC().After(challenge.ExpiresAt) {
		return false
	}
	return HashOTP(inputCode) == challenge.OTPCodeHash
}


func HashOTP(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}
