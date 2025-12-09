// Package database provides Redis connection tests
package database

import (
	"context"
	"os"
	"testing"
	"time"
)

// getTestRedisURL returns Redis URL for testing
// Uses REDIS_URL env var or defaults to localhost
func getTestRedisURL() string {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		url = "redis://localhost:6379"
	}
	return url
}

// TestConnectRedis tests Redis connection
func TestConnectRedis(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	url := getTestRedisURL()
	client, err := ConnectRedis(url)
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	
	// Test ping by trying to set a key
	testKey := "test:connection"
	if err := client.Set(ctx, testKey, "test", 1*time.Second); err != nil {
		t.Fatalf("Failed to set test key: %v", err)
	}
	
	// Cleanup
	_ = client.Delete(ctx, testKey)
}

// TestRedisSetGet tests basic set and get operations
func TestRedisSetGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	url := getTestRedisURL()
	client, err := ConnectRedis(url)
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	testKey := "test:setget"
	testValue := "test-value-123"

	// Test Set
	if err := client.Set(ctx, testKey, testValue, 10*time.Second); err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// Test Get
	val, err := client.Get(ctx, testKey)
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	if val != testValue {
		t.Errorf("Expected value %s, got %s", testValue, val)
	}

	// Cleanup
	if err := client.Delete(ctx, testKey); err != nil {
		t.Errorf("Failed to delete test key: %v", err)
	}
}

// TestRedisDelete tests delete operation
func TestRedisDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	url := getTestRedisURL()
	client, err := ConnectRedis(url)
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	testKey := "test:delete"
	testValue := "test-value"

	// Set a key
	if err := client.Set(ctx, testKey, testValue, 10*time.Second); err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// Verify it exists
	exists, err := client.Exists(ctx, testKey)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if !exists {
		t.Fatal("Key should exist before deletion")
	}

	// Delete it
	if err := client.Delete(ctx, testKey); err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	// Verify it doesn't exist
	exists, err = client.Exists(ctx, testKey)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if exists {
		t.Fatal("Key should not exist after deletion")
	}
}

// TestRedisTTL tests TTL (time to live) behavior
func TestRedisTTL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	url := getTestRedisURL()
	client, err := ConnectRedis(url)
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	testKey := "test:ttl"
	testValue := "test-value"
	ttl := 2 * time.Second

	// Set a key with TTL
	if err := client.Set(ctx, testKey, testValue, ttl); err != nil {
		t.Fatalf("Failed to set key with TTL: %v", err)
	}

	// Verify it exists immediately
	exists, err := client.Exists(ctx, testKey)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if !exists {
		t.Fatal("Key should exist immediately after setting")
	}

	// Wait for TTL to expire
	time.Sleep(ttl + 500*time.Millisecond)

	// Verify it no longer exists
	exists, err = client.Exists(ctx, testKey)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if exists {
		t.Fatal("Key should not exist after TTL expiration")
	}

	// Try to get it (should fail)
	_, err = client.Get(ctx, testKey)
	if err == nil {
		t.Fatal("Get should fail for expired key")
	}
}

// TestRedisSetNX tests SetNX (set if not exists) operation
func TestRedisSetNX(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	url := getTestRedisURL()
	client, err := ConnectRedis(url)
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	testKey := "test:setnx"
	testValue := "test-value"

	// First SetNX should succeed
	created, err := client.SetNX(ctx, testKey, testValue, 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to SetNX: %v", err)
	}
	if !created {
		t.Fatal("First SetNX should succeed")
	}

	// Second SetNX should fail (key already exists)
	created, err = client.SetNX(ctx, testKey, "different-value", 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to SetNX: %v", err)
	}
	if created {
		t.Fatal("Second SetNX should fail (key already exists)")
	}

	// Verify original value is still there
	val, err := client.Get(ctx, testKey)
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}
	if val != testValue {
		t.Errorf("Expected value %s, got %s", testValue, val)
	}

	// Cleanup
	_ = client.Delete(ctx, testKey)
}

// TestRedisIncrement tests atomic increment operations
func TestRedisIncrement(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	url := getTestRedisURL()
	client, err := ConnectRedis(url)
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	testKey := "test:increment"

	// Cleanup before test
	_ = client.Delete(ctx, testKey)

	// First increment should return 1
	val, err := client.Increment(ctx, testKey)
	if err != nil {
		t.Fatalf("Failed to increment: %v", err)
	}
	if val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}

	// Second increment should return 2
	val, err = client.Increment(ctx, testKey)
	if err != nil {
		t.Fatalf("Failed to increment: %v", err)
	}
	if val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}

	// Increment by 5 should return 7
	val, err = client.IncrementBy(ctx, testKey, 5)
	if err != nil {
		t.Fatalf("Failed to increment by 5: %v", err)
	}
	if val != 7 {
		t.Errorf("Expected 7, got %d", val)
	}

	// Cleanup
	_ = client.Delete(ctx, testKey)
}

// TestRedisExpire tests setting expiration on existing keys
func TestRedisExpire(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	url := getTestRedisURL()
	client, err := ConnectRedis(url)
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	testKey := "test:expire"
	testValue := "test-value"

	// Set a key without TTL
	if err := client.Set(ctx, testKey, testValue, 0); err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// Set expiration
	ttl := 2 * time.Second
	if err := client.Expire(ctx, testKey, ttl); err != nil {
		t.Fatalf("Failed to set expiration: %v", err)
	}

	// Verify it exists
	exists, err := client.Exists(ctx, testKey)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if !exists {
		t.Fatal("Key should exist after setting expiration")
	}

	// Wait for expiration
	time.Sleep(ttl + 500*time.Millisecond)

	// Verify it no longer exists
	exists, err = client.Exists(ctx, testKey)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if exists {
		t.Fatal("Key should not exist after expiration")
	}
}

