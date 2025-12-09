// Package services provides test helpers
package services

import (
	"os"
	"testing"

	"github.com/pharmonico/backend-gogit/internal/database"
)

// getTestRedisClient returns a Redis client for testing
func getTestRedisClient(t *testing.T) *database.RedisClient {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		url = "redis://localhost:6379"
	}

	client, err := database.ConnectRedis(url)
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	return client
}

