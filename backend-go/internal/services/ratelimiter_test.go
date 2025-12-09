// Package services provides service layer tests
package services

import (
	"context"
	"testing"
	"time"
)

// TestRateLimiterService_CheckRateLimit tests basic rate limiting
func TestRateLimiterService_CheckRateLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	service := NewRateLimiterService(redisClient)
	ctx := context.Background()

	identifier := "test-ip-123"
	limit := 5
	window := 10 * time.Second

	// Cleanup before test
	_ = service.ResetRateLimit(ctx, identifier)

	// Make requests up to the limit
	for i := 0; i < limit; i++ {
		result, err := service.CheckRateLimit(ctx, identifier, limit, window)
		if err != nil {
			t.Fatalf("Failed to check rate limit: %v", err)
		}

		if !result.Allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}

		expectedRemaining := limit - (i + 1)
		if result.Remaining != expectedRemaining {
			t.Errorf("Expected remaining %d, got %d", expectedRemaining, result.Remaining)
		}
	}

	// Next request should be rate limited
	result, err := service.CheckRateLimit(ctx, identifier, limit, window)
	if err != nil {
		t.Fatalf("Failed to check rate limit: %v", err)
	}

	if result.Allowed {
		t.Error("Request should be rate limited")
	}

	if result.Remaining != 0 {
		t.Errorf("Expected remaining 0, got %d", result.Remaining)
	}

	// Cleanup
	_ = service.ResetRateLimit(ctx, identifier)
}

// TestRateLimiterService_CheckRateLimit_WindowExpiry tests rate limit window expiry
func TestRateLimiterService_CheckRateLimit_WindowExpiry(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	service := NewRateLimiterService(redisClient)
	ctx := context.Background()

	identifier := "test-ip-window"
	limit := 3
	window := 2 * time.Second

	// Cleanup before test
	_ = service.ResetRateLimit(ctx, identifier)

	// Exhaust the limit
	for i := 0; i < limit; i++ {
		_, err := service.CheckRateLimit(ctx, identifier, limit, window)
		if err != nil {
			t.Fatalf("Failed to check rate limit: %v", err)
		}
	}

	// Verify we're rate limited
	result, err := service.CheckRateLimit(ctx, identifier, limit, window)
	if err != nil {
		t.Fatalf("Failed to check rate limit: %v", err)
	}
	if result.Allowed {
		t.Error("Should be rate limited")
	}

	// Wait for window to expire
	time.Sleep(window + 500*time.Millisecond)

	// Should be able to make requests again
	result, err = service.CheckRateLimit(ctx, identifier, limit, window)
	if err != nil {
		t.Fatalf("Failed to check rate limit: %v", err)
	}
	if !result.Allowed {
		t.Error("Should be able to make requests after window expiry")
	}

	// Cleanup
	_ = service.ResetRateLimit(ctx, identifier)
}

// TestRateLimiterService_CheckRateLimitWithSlidingWindow tests sliding window rate limiting
func TestRateLimiterService_CheckRateLimitWithSlidingWindow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	service := NewRateLimiterService(redisClient)
	ctx := context.Background()

	identifier := "test-ip-sliding"
	limit := 5
	window := 10 * time.Second

	// Cleanup before test
	_ = service.ResetRateLimit(ctx, identifier)

	// Make requests up to the limit
	for i := 0; i < limit; i++ {
		result, err := service.CheckRateLimitWithSlidingWindow(ctx, identifier, limit, window)
		if err != nil {
			t.Fatalf("Failed to check rate limit: %v", err)
		}

		if !result.Allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Next request should be rate limited
	result, err := service.CheckRateLimitWithSlidingWindow(ctx, identifier, limit, window)
	if err != nil {
		t.Fatalf("Failed to check rate limit: %v", err)
	}

	if result.Allowed {
		t.Error("Request should be rate limited")
	}

	// Cleanup
	_ = service.ResetRateLimit(ctx, identifier)
}

// TestRateLimiterService_ResetRateLimit tests resetting rate limits
func TestRateLimiterService_ResetRateLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	service := NewRateLimiterService(redisClient)
	ctx := context.Background()

	identifier := "test-ip-reset"
	limit := 3
	window := 10 * time.Second

	// Exhaust the limit
	for i := 0; i < limit; i++ {
		_, err := service.CheckRateLimit(ctx, identifier, limit, window)
		if err != nil {
			t.Fatalf("Failed to check rate limit: %v", err)
		}
	}

	// Verify we're rate limited
	result, err := service.CheckRateLimit(ctx, identifier, limit, window)
	if err != nil {
		t.Fatalf("Failed to check rate limit: %v", err)
	}
	if result.Allowed {
		t.Error("Should be rate limited")
	}

	// Reset rate limit
	err = service.ResetRateLimit(ctx, identifier)
	if err != nil {
		t.Fatalf("Failed to reset rate limit: %v", err)
	}

	// Should be able to make requests again
	result, err = service.CheckRateLimit(ctx, identifier, limit, window)
	if err != nil {
		t.Fatalf("Failed to check rate limit: %v", err)
	}
	if !result.Allowed {
		t.Error("Should be able to make requests after reset")
	}
}
