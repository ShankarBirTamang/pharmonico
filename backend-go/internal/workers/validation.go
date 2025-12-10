// Package workers provides worker handlers for processing Kafka events
package workers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pharmonico/backend-gogit/internal/database"
	"github.com/pharmonico/backend-gogit/internal/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ValidationWorker handles prescription validation events
type ValidationWorker struct {
	mongoClient   *database.MongoClient
	kafkaProducer kafka.Producer
}

// NewValidationWorker creates a new validation worker
func NewValidationWorker(mongoClient *database.MongoClient, kafkaProducer kafka.Producer) *ValidationWorker {
	return &ValidationWorker{
		mongoClient:   mongoClient,
		kafkaProducer: kafkaProducer,
	}
}

// Topic returns the Kafka topic this handler consumes from
func (w *ValidationWorker) Topic() string {
	return kafka.TopicIntakeReceived
}

// Handle processes a prescription intake event and validates it
func (w *ValidationWorker) Handle(ctx context.Context, msg *kafka.Message) error {
	// Parse the event payload
	var event struct {
		EventID        string    `json:"event_id"`
		PrescriptionID string    `json:"prescription_id"`
		PatientID      string    `json:"patient_id,omitempty"`
		DrugNDC        string    `json:"drug_ndc,omitempty"`
		Timestamp      time.Time `json:"timestamp"`
	}

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("❌ Failed to unmarshal validation event: %v", err)
		return err
	}

	log.Printf("🔍 Processing validation for prescription: %s", event.PrescriptionID)

	// Fetch prescription from MongoDB
	prescriptionCollection := w.mongoClient.GetCollection("prescriptions")
	prescriptionID, err := primitive.ObjectIDFromHex(event.PrescriptionID)
	if err != nil {
		log.Printf("❌ Invalid prescription ID format: %s", event.PrescriptionID)
		return err
	}

	var prescription bson.M
	err = prescriptionCollection.FindOne(ctx, bson.M{"_id": prescriptionID}).Decode(&prescription)
	if err != nil {
		log.Printf("❌ Prescription not found: %s", event.PrescriptionID)
		return err
	}

	// Perform validation (simplified - actual validation logic will be implemented in task 1.2)
	validationErrors := []string{}
	isValid := true

	// Basic validation checks
	if prescription["patient_id"] == nil || prescription["patient_id"] == "" {
		validationErrors = append(validationErrors, "patient_id is required")
		isValid = false
	}

	if prescription["prescriber_id"] == nil || prescription["prescriber_id"] == "" {
		validationErrors = append(validationErrors, "prescriber_id is required")
		isValid = false
	}

	if prescription["medications"] == nil {
		validationErrors = append(validationErrors, "medications are required")
		isValid = false
	}

	// Update prescription status in MongoDB
	update := bson.M{
		"$set": bson.M{
			"status":            getValidationStatus(isValid),
			"validation_errors": validationErrors,
			"updated_at":        time.Now(),
		},
	}

	_, err = prescriptionCollection.UpdateOne(ctx, bson.M{"_id": prescriptionID}, update)
	if err != nil {
		log.Printf("❌ Failed to update prescription status: %v", err)
		return err
	}

	// If validation passed, emit validation completed event
	if isValid {
		validationEvent := map[string]interface{}{
			"event_id":        uuid.New().String(),
			"prescription_id": event.PrescriptionID,
			"patient_id":      event.PatientID,
			"validated_at":    time.Now().Format(time.RFC3339),
			"timestamp":       time.Now().Format(time.RFC3339),
		}

		eventBytes, err := json.Marshal(validationEvent)
		if err != nil {
			log.Printf("❌ Failed to marshal validation event: %v", err)
			return err
		}

		// Emit validation completed event
		if err := w.kafkaProducer.Publish(ctx, kafka.TopicValidationCompleted, event.PrescriptionID, eventBytes); err != nil {
			log.Printf("❌ Failed to publish validation completed event: %v", err)
			return err
		}

		log.Printf("✅ Validation completed for prescription: %s", event.PrescriptionID)
	} else {
		log.Printf("⚠️  Validation failed for prescription: %s - Errors: %v", event.PrescriptionID, validationErrors)
	}

	return nil
}

// getValidationStatus returns the status based on validation result
func getValidationStatus(isValid bool) string {
	if isValid {
		return "validated"
	}
	return "validation_failed"
}
