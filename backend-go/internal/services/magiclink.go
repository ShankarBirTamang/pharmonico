// Package services provides business logic services
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pharmonico/backend-gogit/internal/database"
)

// MagicLinkService handles magic link token operations
type MagicLinkService struct {
	redis *database.RedisClient
}

// MagicLinkData represents the data stored for a magic link token
type MagicLinkData struct {
	PrescriptionID string    `json:"prescription_id"`
	PatientID      string    `json:"patient_id"`
	ExpiresAt      time.Time `json:"expires_at"`
	Used           bool      `json:"used"`
	CreatedAt      time.Time `json:"created_at"`
}

// NewMagicLinkService creates a new magic link service
func NewMagicLinkService(redis *database.RedisClient) *MagicLinkService {
	return &MagicLinkService{
		redis: redis,
	}
}

// GenerateToken creates a new magic link token and stores it in Redis
// token: the unique token (typically UUID)
// data: the enrollment data to store
// ttl: time to live for the token (default 48 hours)
func (s *MagicLinkService) GenerateToken(ctx context.Context, token string, prescriptionID, patientID string, ttl time.Duration) error {
	if ttl == 0 {
		ttl = 48 * time.Hour // Default 48 hours
	}

	now := time.Now()
	data := MagicLinkData{
		PrescriptionID: prescriptionID,
		PatientID:      patientID,
		ExpiresAt:      now.Add(ttl),
		Used:           false,
		CreatedAt:      now,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal magic link data: %w", err)
	}

	key := fmt.Sprintf("magic_link:%s", token)
	if err := s.redis.Set(ctx, key, string(jsonData), ttl); err != nil {
		return fmt.Errorf("failed to store magic link token: %w", err)
	}

	return nil
}

// ValidateToken retrieves and validates a magic link token
// Returns the enrollment data if valid, or an error if invalid/expired/used
func (s *MagicLinkService) ValidateToken(ctx context.Context, token string) (*MagicLinkData, error) {
	key := fmt.Sprintf("magic_link:%s", token)

	val, err := s.redis.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("token not found or expired: %w", err)
	}

	var data MagicLinkData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal magic link data: %w", err)
	}

	// Check if already used
	if data.Used {
		return nil, fmt.Errorf("token already used")
	}

	// Check expiration
	if time.Now().After(data.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	return &data, nil
}

// MarkAsUsed marks a magic link token as used
func (s *MagicLinkService) MarkAsUsed(ctx context.Context, token string) error {
	key := fmt.Sprintf("magic_link:%s", token)

	// Get existing data
	val, err := s.redis.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("token not found: %w", err)
	}

	var data MagicLinkData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return fmt.Errorf("failed to unmarshal magic link data: %w", err)
	}

	// Mark as used
	data.Used = true

	// Calculate remaining TTL
	remainingTTL := time.Until(data.ExpiresAt)
	if remainingTTL < 0 {
		remainingTTL = 1 * time.Minute // Keep for at least 1 minute for audit purposes
	}

	// Update in Redis
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal updated data: %w", err)
	}

	if err := s.redis.Set(ctx, key, string(jsonData), remainingTTL); err != nil {
		return fmt.Errorf("failed to update token: %w", err)
	}

	return nil
}

// DeleteToken removes a magic link token from Redis
func (s *MagicLinkService) DeleteToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("magic_link:%s", token)
	return s.redis.Delete(ctx, key)
}
