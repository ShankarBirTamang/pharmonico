// Package services provides business logic services
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/phil-my-meds/backend-gogit/internal/database"
)

// PharmacyCapacityService handles pharmacy capacity tracking in Redis
type PharmacyCapacityService struct {
	redis *database.RedisClient
}

// CapacityData represents the capacity information for a pharmacy
type CapacityData struct {
	PharmacyID     string    `json:"pharmacy_id"`
	CurrentDailyRx int       `json:"current_daily_rx"`
	MaxDailyRx     int       `json:"max_daily_rx"`
	Utilization    float64   `json:"utilization"` // current_daily_rx / max_daily_rx
	LastUpdated    time.Time `json:"last_updated"`
}

// NewPharmacyCapacityService creates a new pharmacy capacity service
func NewPharmacyCapacityService(redis *database.RedisClient) *PharmacyCapacityService {
	return &PharmacyCapacityService{
		redis: redis,
	}
}

// GetCapacity retrieves the current capacity for a pharmacy
func (s *PharmacyCapacityService) GetCapacity(ctx context.Context, pharmacyID string) (*CapacityData, error) {
	key := fmt.Sprintf("pharmacy_capacity:%s", pharmacyID)

	val, err := s.redis.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("capacity not found for pharmacy %s: %w", pharmacyID, err)
	}

	var data CapacityData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal capacity data: %w", err)
	}

	return &data, nil
}

// SetCapacity sets the capacity for a pharmacy
// ttl: time to live (default 5 minutes for cache)
func (s *PharmacyCapacityService) SetCapacity(ctx context.Context, pharmacyID string, currentDailyRx, maxDailyRx int, ttl time.Duration) error {
	if ttl == 0 {
		ttl = 5 * time.Minute // Default 5 minutes
	}

	utilization := float64(currentDailyRx) / float64(maxDailyRx)
	if maxDailyRx == 0 {
		utilization = 0
	}

	data := CapacityData{
		PharmacyID:     pharmacyID,
		CurrentDailyRx: currentDailyRx,
		MaxDailyRx:     maxDailyRx,
		Utilization:    utilization,
		LastUpdated:    time.Now(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal capacity data: %w", err)
	}

	key := fmt.Sprintf("pharmacy_capacity:%s", pharmacyID)
	if err := s.redis.Set(ctx, key, string(jsonData), ttl); err != nil {
		return fmt.Errorf("failed to store capacity: %w", err)
	}

	return nil
}

// IncrementCapacity increments the current daily RX count for a pharmacy
// Returns the new current count and utilization
func (s *PharmacyCapacityService) IncrementCapacity(ctx context.Context, pharmacyID string) (int, float64, error) {
	key := fmt.Sprintf("pharmacy_capacity:%s", pharmacyID)

	// Try to get existing capacity
	val, err := s.redis.Get(ctx, key)
	var data CapacityData

	if err != nil {
		// If not found, initialize with default values
		data = CapacityData{
			PharmacyID:     pharmacyID,
			CurrentDailyRx: 0,
			MaxDailyRx:     100, // Default max
			Utilization:    0,
			LastUpdated:    time.Now(),
		}
	} else {
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			return 0, 0, fmt.Errorf("failed to unmarshal capacity data: %w", err)
		}
	}

	// Increment
	data.CurrentDailyRx++
	data.Utilization = float64(data.CurrentDailyRx) / float64(data.MaxDailyRx)
	if data.MaxDailyRx == 0 {
		data.Utilization = 0
	}
	data.LastUpdated = time.Now()

	// Save back to Redis
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to marshal updated capacity data: %w", err)
	}

	ttl := 5 * time.Minute
	if err := s.redis.Set(ctx, key, string(jsonData), ttl); err != nil {
		return 0, 0, fmt.Errorf("failed to update capacity: %w", err)
	}

	return data.CurrentDailyRx, data.Utilization, nil
}

// DecrementCapacity decrements the current daily RX count for a pharmacy
func (s *PharmacyCapacityService) DecrementCapacity(ctx context.Context, pharmacyID string) (int, float64, error) {
	key := fmt.Sprintf("pharmacy_capacity:%s", pharmacyID)

	val, err := s.redis.Get(ctx, key)
	if err != nil {
		return 0, 0, fmt.Errorf("capacity not found for pharmacy %s: %w", pharmacyID, err)
	}

	var data CapacityData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal capacity data: %w", err)
	}

	// Decrement (but don't go below 0)
	if data.CurrentDailyRx > 0 {
		data.CurrentDailyRx--
	}
	data.Utilization = float64(data.CurrentDailyRx) / float64(data.MaxDailyRx)
	if data.MaxDailyRx == 0 {
		data.Utilization = 0
	}
	data.LastUpdated = time.Now()

	// Save back to Redis
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to marshal updated capacity data: %w", err)
	}

	ttl := 5 * time.Minute
	if err := s.redis.Set(ctx, key, string(jsonData), ttl); err != nil {
		return 0, 0, fmt.Errorf("failed to update capacity: %w", err)
	}

	return data.CurrentDailyRx, data.Utilization, nil
}

// HasCapacity checks if a pharmacy has available capacity
// Returns true if utilization is below the threshold (default 0.95)
func (s *PharmacyCapacityService) HasCapacity(ctx context.Context, pharmacyID string, threshold float64) (bool, error) {
	if threshold == 0 {
		threshold = 0.95 // Default 95% threshold
	}

	capacity, err := s.GetCapacity(ctx, pharmacyID)
	if err != nil {
		return false, err
	}

	return capacity.Utilization < threshold, nil
}

// DeleteCapacity removes capacity data for a pharmacy
func (s *PharmacyCapacityService) DeleteCapacity(ctx context.Context, pharmacyID string) error {
	key := fmt.Sprintf("pharmacy_capacity:%s", pharmacyID)
	return s.redis.Delete(ctx, key)
}
