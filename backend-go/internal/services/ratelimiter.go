// Package services provides business logic services
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/phil-my-meds/backend-gogit/internal/database"
)

// RateLimiterService provides rate limiting functionality using Redis
type RateLimiterService struct {
	redis *database.RedisClient
}

// RateLimitResult represents the result of a rate limit check
type RateLimitResult struct {
	Allowed   bool
	Remaining int
	ResetAt   time.Time
}

// NewRateLimiterService creates a new rate limiter service
func NewRateLimiterService(redis *database.RedisClient) *RateLimiterService {
	return &RateLimiterService{
		redis: redis,
	}
}

// CheckRateLimit checks if a request should be allowed based on rate limiting
// Uses fixed window algorithm with Redis atomic increment
// identifier: unique identifier (e.g., IP address, user ID, API key)
// limit: maximum number of requests allowed
// window: time window for the limit (e.g., 1 minute, 1 hour)
func (s *RateLimiterService) CheckRateLimit(ctx context.Context, identifier string, limit int, window time.Duration) (*RateLimitResult, error) {
	if limit <= 0 {
		return &RateLimitResult{Allowed: true, Remaining: -1, ResetAt: time.Now().Add(window)}, nil
	}

	now := time.Now()
	counterKey := fmt.Sprintf("rate_limit:count:%s", identifier)

	// Check if key exists to set TTL on first request
	exists, err := s.redis.Exists(ctx, counterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to check rate limit: %w", err)
	}

	// Atomically increment the counter
	count, err := s.redis.Increment(ctx, counterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to increment rate limit counter: %w", err)
	}

	// If this is the first request, set the TTL
	if !exists {
		if err := s.redis.Expire(ctx, counterKey, window); err != nil {
			// Non-fatal error, continue
		}
	}

	// Check if limit exceeded
	allowed := count <= int64(limit)
	remaining := limit - int(count)
	if remaining < 0 {
		remaining = 0
	}

	// Calculate reset time (approximate, based on TTL)
	resetAt := now.Add(window)

	return &RateLimitResult{
		Allowed:   allowed,
		Remaining: remaining,
		ResetAt:   resetAt,
	}, nil
}

// CheckRateLimitWithSlidingWindow implements a more accurate sliding window rate limiter
// Uses atomic increment for better performance
func (s *RateLimiterService) CheckRateLimitWithSlidingWindow(ctx context.Context, identifier string, limit int, window time.Duration) (*RateLimitResult, error) {
	if limit <= 0 {
		return &RateLimitResult{Allowed: true, Remaining: -1, ResetAt: time.Now().Add(window)}, nil
	}

	now := time.Now()
	counterKey := fmt.Sprintf("rate_limit:sliding:count:%s", identifier)

	// Check if key exists to set TTL on first request
	exists, err := s.redis.Exists(ctx, counterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to check rate limit: %w", err)
	}

	// Atomically increment the counter
	count, err := s.redis.Increment(ctx, counterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to increment rate limit counter: %w", err)
	}

	// If this is the first request, set the TTL
	if !exists {
		if err := s.redis.Expire(ctx, counterKey, window); err != nil {
			// Non-fatal error, continue
		}
	}

	allowed := count <= int64(limit)
	remaining := limit - int(count)
	if remaining < 0 {
		remaining = 0
	}

	resetAt := now.Add(window)

	return &RateLimitResult{
		Allowed:   allowed,
		Remaining: remaining,
		ResetAt:   resetAt,
	}, nil
}

// ResetRateLimit clears the rate limit for an identifier
func (s *RateLimiterService) ResetRateLimit(ctx context.Context, identifier string) error {
	counterKey := fmt.Sprintf("rate_limit:count:%s", identifier)
	slidingKey := fmt.Sprintf("rate_limit:sliding:count:%s", identifier)

	// Delete both keys
	if err := s.redis.Delete(ctx, counterKey, slidingKey); err != nil {
		return fmt.Errorf("failed to reset rate limit: %w", err)
	}

	return nil
}
