// Package services provides service layer tests
package services

import (
	"context"
	"testing"
	"time"
)

// TestMagicLinkService_GenerateToken tests token generation
func TestMagicLinkService_GenerateToken(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	service := NewMagicLinkService(redisClient)
	ctx := context.Background()

	token := "test-token-123"
	prescriptionID := "rx_abc123"
	patientID := "pat_def456"

	// Generate token
	err := service.GenerateToken(ctx, token, prescriptionID, patientID, 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Cleanup
	defer service.DeleteToken(ctx, token)
}

// TestMagicLinkService_ValidateToken tests token validation
func TestMagicLinkService_ValidateToken(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	service := NewMagicLinkService(redisClient)
	ctx := context.Background()

	token := "test-token-validate"
	prescriptionID := "rx_abc123"
	patientID := "pat_def456"

	// Generate token
	err := service.GenerateToken(ctx, token, prescriptionID, patientID, 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	defer service.DeleteToken(ctx, token)

	// Validate token
	data, err := service.ValidateToken(ctx, token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if data.PrescriptionID != prescriptionID {
		t.Errorf("Expected prescription ID %s, got %s", prescriptionID, data.PrescriptionID)
	}

	if data.PatientID != patientID {
		t.Errorf("Expected patient ID %s, got %s", patientID, data.PatientID)
	}

	if data.Used {
		t.Error("Token should not be marked as used")
	}
}

// TestMagicLinkService_ValidateToken_Expired tests expired token validation
func TestMagicLinkService_ValidateToken_Expired(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	service := NewMagicLinkService(redisClient)
	ctx := context.Background()

	token := "test-token-expired"
	prescriptionID := "rx_abc123"
	patientID := "pat_def456"

	// Generate token with short TTL
	err := service.GenerateToken(ctx, token, prescriptionID, patientID, 1*time.Second)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// Try to validate expired token
	_, err = service.ValidateToken(ctx, token)
	if err == nil {
		t.Fatal("Should fail to validate expired token")
	}
}

// TestMagicLinkService_MarkAsUsed tests marking token as used
func TestMagicLinkService_MarkAsUsed(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	redisClient := getTestRedisClient(t)
	defer redisClient.Close()

	service := NewMagicLinkService(redisClient)
	ctx := context.Background()

	token := "test-token-used"
	prescriptionID := "rx_abc123"
	patientID := "pat_def456"

	// Generate token
	err := service.GenerateToken(ctx, token, prescriptionID, patientID, 1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	defer service.DeleteToken(ctx, token)

	// Mark as used
	err = service.MarkAsUsed(ctx, token)
	if err != nil {
		t.Fatalf("Failed to mark token as used: %v", err)
	}

	// Try to validate used token (should fail)
	_, err = service.ValidateToken(ctx, token)
	if err == nil {
		t.Fatal("Should fail to validate used token")
	}
}
