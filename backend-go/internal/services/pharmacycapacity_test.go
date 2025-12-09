// Package services provides service layer tests
package services

import (
	"context"
	"testing"
)

// TestPharmacyCapacityService_SetGetCapacity tests setting and getting capacity
func TestPharmacyCapacityService_SetGetCapacity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	service := NewPharmacyCapacityService(redisClient)
	ctx := context.Background()

	pharmacyID := "pharm_test_123"
	currentDailyRx := 45
	maxDailyRx := 100

	// Set capacity
	err := service.SetCapacity(ctx, pharmacyID, currentDailyRx, maxDailyRx, 0)
	if err != nil {
		t.Fatalf("Failed to set capacity: %v", err)
	}
	defer service.DeleteCapacity(ctx, pharmacyID)

	// Get capacity
	capacity, err := service.GetCapacity(ctx, pharmacyID)
	if err != nil {
		t.Fatalf("Failed to get capacity: %v", err)
	}

	if capacity.PharmacyID != pharmacyID {
		t.Errorf("Expected pharmacy ID %s, got %s", pharmacyID, capacity.PharmacyID)
	}

	if capacity.CurrentDailyRx != currentDailyRx {
		t.Errorf("Expected current daily RX %d, got %d", currentDailyRx, capacity.CurrentDailyRx)
	}

	if capacity.MaxDailyRx != maxDailyRx {
		t.Errorf("Expected max daily RX %d, got %d", maxDailyRx, capacity.MaxDailyRx)
	}

	expectedUtilization := float64(currentDailyRx) / float64(maxDailyRx)
	if capacity.Utilization != expectedUtilization {
		t.Errorf("Expected utilization %.2f, got %.2f", expectedUtilization, capacity.Utilization)
	}
}

// TestPharmacyCapacityService_IncrementCapacity tests incrementing capacity
func TestPharmacyCapacityService_IncrementCapacity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	service := NewPharmacyCapacityService(redisClient)
	ctx := context.Background()

	pharmacyID := "pharm_test_increment"
	currentDailyRx := 10
	maxDailyRx := 100

	// Set initial capacity
	err := service.SetCapacity(ctx, pharmacyID, currentDailyRx, maxDailyRx, 0)
	if err != nil {
		t.Fatalf("Failed to set capacity: %v", err)
	}
	defer service.DeleteCapacity(ctx, pharmacyID)

	// Increment capacity
	newCount, utilization, err := service.IncrementCapacity(ctx, pharmacyID)
	if err != nil {
		t.Fatalf("Failed to increment capacity: %v", err)
	}

	if newCount != currentDailyRx+1 {
		t.Errorf("Expected count %d, got %d", currentDailyRx+1, newCount)
	}

	expectedUtilization := float64(newCount) / float64(maxDailyRx)
	if utilization != expectedUtilization {
		t.Errorf("Expected utilization %.2f, got %.2f", expectedUtilization, utilization)
	}
}

// TestPharmacyCapacityService_DecrementCapacity tests decrementing capacity
func TestPharmacyCapacityService_DecrementCapacity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	service := NewPharmacyCapacityService(redisClient)
	ctx := context.Background()

	pharmacyID := "pharm_test_decrement"
	currentDailyRx := 10
	maxDailyRx := 100

	// Set initial capacity
	err := service.SetCapacity(ctx, pharmacyID, currentDailyRx, maxDailyRx, 0)
	if err != nil {
		t.Fatalf("Failed to set capacity: %v", err)
	}
	defer service.DeleteCapacity(ctx, pharmacyID)

	// Decrement capacity
	newCount, utilization, err := service.DecrementCapacity(ctx, pharmacyID)
	if err != nil {
		t.Fatalf("Failed to decrement capacity: %v", err)
	}

	if newCount != currentDailyRx-1 {
		t.Errorf("Expected count %d, got %d", currentDailyRx-1, newCount)
	}

	expectedUtilization := float64(newCount) / float64(maxDailyRx)
	if utilization != expectedUtilization {
		t.Errorf("Expected utilization %.2f, got %.2f", expectedUtilization, utilization)
	}
}

// TestPharmacyCapacityService_HasCapacity tests capacity check
func TestPharmacyCapacityService_HasCapacity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	service := NewPharmacyCapacityService(redisClient)
	ctx := context.Background()

	pharmacyID := "pharm_test_has_capacity"
	maxDailyRx := 100

	// Test with low utilization (should have capacity)
	err := service.SetCapacity(ctx, pharmacyID, 50, maxDailyRx, 0)
	if err != nil {
		t.Fatalf("Failed to set capacity: %v", err)
	}
	defer service.DeleteCapacity(ctx, pharmacyID)

	hasCapacity, err := service.HasCapacity(ctx, pharmacyID, 0.95)
	if err != nil {
		t.Fatalf("Failed to check capacity: %v", err)
	}
	if !hasCapacity {
		t.Error("Pharmacy should have capacity at 50% utilization")
	}

	// Test with high utilization (should not have capacity)
	err = service.SetCapacity(ctx, pharmacyID, 96, maxDailyRx, 0)
	if err != nil {
		t.Fatalf("Failed to set capacity: %v", err)
	}

	hasCapacity, err = service.HasCapacity(ctx, pharmacyID, 0.95)
	if err != nil {
		t.Fatalf("Failed to check capacity: %v", err)
	}
	if hasCapacity {
		t.Error("Pharmacy should not have capacity at 96% utilization")
	}
}

