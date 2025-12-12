// Package main provides a test script to validate Redis functionality and TTL behavior
package main

import (
	"context"
	"log"
	"time"

	"github.com/phil-my-meds/backend-gogit/internal/config"
	"github.com/phil-my-meds/backend-gogit/internal/database"
	"github.com/phil-my-meds/backend-gogit/internal/services"
)

func main() {
	log.Println("🧪 Testing Redis Functionality and TTL Behavior...")
	log.Println("=" + string(make([]byte, 60)))

	cfg := config.Load()

	// Connect to Redis
	log.Println("\n📡 Connecting to Redis...")
	redisClient, err := database.ConnectRedis(cfg.RedisURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	log.Println("✅ Redis connected successfully")

	ctx := context.Background()

	// Test 1: Basic Set/Get
	log.Println("\n📝 Test 1: Basic Set/Get Operations")
	testKey := "test:basic:setget"
	testValue := "test-value-123"
	if err := redisClient.Set(ctx, testKey, testValue, 10*time.Second); err != nil {
		log.Fatalf("❌ Failed to set key: %v", err)
	}
	log.Printf("   ✓ Set key: %s = %s", testKey, testValue)

	val, err := redisClient.Get(ctx, testKey)
	if err != nil {
		log.Fatalf("❌ Failed to get key: %v", err)
	}
	log.Printf("   ✓ Got value: %s", val)

	if val != testValue {
		log.Fatalf("❌ Value mismatch: expected %s, got %s", testValue, val)
	}
	log.Println("   ✅ Set/Get test passed")

	// Test 2: TTL Behavior
	log.Println("\n⏱️  Test 2: TTL (Time To Live) Behavior")
	ttlKey := "test:ttl:expiry"
	ttlValue := "ttl-test-value"
	ttl := 3 * time.Second

	if err := redisClient.Set(ctx, ttlKey, ttlValue, ttl); err != nil {
		log.Fatalf("❌ Failed to set key with TTL: %v", err)
	}
	log.Printf("   ✓ Set key with TTL: %s (expires in %v)", ttlKey, ttl)

	// Check immediately
	exists, err := redisClient.Exists(ctx, ttlKey)
	if err != nil {
		log.Fatalf("❌ Failed to check existence: %v", err)
	}
	if !exists {
		log.Fatalf("❌ Key should exist immediately after setting")
	}
	log.Println("   ✓ Key exists immediately")

	// Wait for TTL to expire
	log.Printf("   ⏳ Waiting %v for TTL to expire...", ttl+500*time.Millisecond)
	time.Sleep(ttl + 500*time.Millisecond)

	// Check after expiration
	exists, err = redisClient.Exists(ctx, ttlKey)
	if err != nil {
		log.Fatalf("❌ Failed to check existence: %v", err)
	}
	if exists {
		log.Fatalf("❌ Key should not exist after TTL expiration")
	}
	log.Println("   ✓ Key expired after TTL")
	log.Println("   ✅ TTL test passed")

	// Test 3: Magic Link Service
	log.Println("\n🔗 Test 3: Magic Link Token Store")
	magicLinkService := services.NewMagicLinkService(redisClient)
	token := "test-magic-link-token"
	prescriptionID := "rx_test_123"
	patientID := "pat_test_456"

	if err := magicLinkService.GenerateToken(ctx, token, prescriptionID, patientID, 1*time.Hour); err != nil {
		log.Fatalf("❌ Failed to generate magic link token: %v", err)
	}
	log.Printf("   ✓ Generated magic link token: %s", token)

	data, err := magicLinkService.ValidateToken(ctx, token)
	if err != nil {
		log.Fatalf("❌ Failed to validate token: %v", err)
	}
	log.Printf("   ✓ Validated token - Prescription: %s, Patient: %s", data.PrescriptionID, data.PatientID)

	if err := magicLinkService.MarkAsUsed(ctx, token); err != nil {
		log.Fatalf("❌ Failed to mark token as used: %v", err)
	}
	log.Println("   ✓ Marked token as used")

	_, err = magicLinkService.ValidateToken(ctx, token)
	if err == nil {
		log.Fatalf("❌ Should fail to validate used token")
	}
	log.Println("   ✓ Used token correctly rejected")
	log.Println("   ✅ Magic Link service test passed")

	// Test 4: Pharmacy Capacity Service
	log.Println("\n🏥 Test 4: Pharmacy Capacity Store")
	capacityService := services.NewPharmacyCapacityService(redisClient)
	pharmacyID := "pharm_test_123"
	currentDailyRx := 45
	maxDailyRx := 100

	if err := capacityService.SetCapacity(ctx, pharmacyID, currentDailyRx, maxDailyRx, 5*time.Minute); err != nil {
		log.Fatalf("❌ Failed to set capacity: %v", err)
	}
	log.Printf("   ✓ Set capacity: %d/%d", currentDailyRx, maxDailyRx)

	capacity, err := capacityService.GetCapacity(ctx, pharmacyID)
	if err != nil {
		log.Fatalf("❌ Failed to get capacity: %v", err)
	}
	log.Printf("   ✓ Got capacity: %d/%d (%.2f%% utilization)", capacity.CurrentDailyRx, capacity.MaxDailyRx, capacity.Utilization*100)

	newCount, utilization, err := capacityService.IncrementCapacity(ctx, pharmacyID)
	if err != nil {
		log.Fatalf("❌ Failed to increment capacity: %v", err)
	}
	log.Printf("   ✓ Incremented capacity: %d/%d (%.2f%% utilization)", newCount, maxDailyRx, utilization*100)

	hasCapacity, err := capacityService.HasCapacity(ctx, pharmacyID, 0.95)
	if err != nil {
		log.Fatalf("❌ Failed to check capacity: %v", err)
	}
	log.Printf("   ✓ Has capacity: %v", hasCapacity)
	log.Println("   ✅ Pharmacy Capacity service test passed")

	// Test 5: Rate Limiter Service
	log.Println("\n🚦 Test 5: Rate Limiter")
	rateLimiterService := services.NewRateLimiterService(redisClient)
	identifier := "test-ip-123"
	limit := 5
	window := 10 * time.Second

	// Cleanup
	_ = rateLimiterService.ResetRateLimit(ctx, identifier)

	// Make requests up to limit
	for i := 0; i < limit; i++ {
		result, err := rateLimiterService.CheckRateLimit(ctx, identifier, limit, window)
		if err != nil {
			log.Fatalf("❌ Failed to check rate limit: %v", err)
		}
		if !result.Allowed {
			log.Fatalf("❌ Request %d should be allowed", i+1)
		}
		log.Printf("   ✓ Request %d/%d allowed (remaining: %d)", i+1, limit, result.Remaining)
	}

	// Next request should be rate limited
	result, err := rateLimiterService.CheckRateLimit(ctx, identifier, limit, window)
	if err != nil {
		log.Fatalf("❌ Failed to check rate limit: %v", err)
	}
	if result.Allowed {
		log.Fatalf("❌ Request should be rate limited")
	}
	log.Printf("   ✓ Request rate limited (remaining: %d)", result.Remaining)
	log.Println("   ✅ Rate Limiter service test passed")

	// Cleanup
	log.Println("\n🧹 Cleaning up test keys...")
	_ = redisClient.Delete(ctx, testKey)
	_ = magicLinkService.DeleteToken(ctx, token)
	_ = capacityService.DeleteCapacity(ctx, pharmacyID)
	_ = rateLimiterService.ResetRateLimit(ctx, identifier)
	log.Println("   ✓ Cleanup complete")

	log.Println("\n" + "=" + string(make([]byte, 60)))
	log.Println("✅ All Redis tests passed successfully!")
	log.Println("\n📊 Summary:")
	log.Println("   ✓ Basic Set/Get operations")
	log.Println("   ✓ TTL expiration behavior")
	log.Println("   ✓ Magic Link token store")
	log.Println("   ✓ Pharmacy capacity store")
	log.Println("   ✓ Rate limiter functionality")
	log.Println("\n💡 To view Redis logs, check the Redis container:")
	log.Println("   docker logs phil-my-meds-redis")
}

