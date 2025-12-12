// Package workers provides worker handlers for processing Kafka events
package workers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/phil-my-meds/backend-gogit/internal/database"
	"github.com/phil-my-meds/backend-gogit/internal/kafka"
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
// 8.3.1: Consumes Kafka event (handled by worker loop)
// 8.3.2: Processes business logic (validation)
// 8.3.3: Emits next Kafka event (validation.completed)
func (w *ValidationWorker) Handle(ctx context.Context, msg *kafka.Message) error {
	// Extract correlation ID from message
	correlationID := ExtractCorrelationID(msg)

	// Parse the event payload
	var event struct {
		EventID        string    `json:"event_id"`
		CorrelationID  string    `json:"correlation_id,omitempty"`
		PrescriptionID string    `json:"prescription_id"`
		PatientID      string    `json:"patient_id,omitempty"`
		DrugNDC        string    `json:"drug_ndc,omitempty"`
		Timestamp      time.Time `json:"timestamp"`
	}

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("❌ [correlation_id=%s] Failed to unmarshal validation event: %v", correlationID, err)
		return err
	}

	// Use correlation ID from event if available, otherwise use extracted one
	if event.CorrelationID != "" {
		correlationID = event.CorrelationID
	}

	log.Printf("🔍 [correlation_id=%s] Processing validation for prescription: %s", correlationID, event.PrescriptionID)

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

	// 8.3.3: Emit next Kafka event if validation passed
	if isValid {
		validationEvent := CreateEvent(correlationID, event.PrescriptionID, map[string]interface{}{
			"patient_id":   event.PatientID,
			"validated_at": time.Now().Format(time.RFC3339),
		})

		if err := PublishEvent(ctx, w.kafkaProducer, kafka.TopicValidationCompleted, event.PrescriptionID, validationEvent); err != nil {
			log.Printf("❌ [correlation_id=%s] Failed to publish validation completed event: %v", correlationID, err)
			return err
		}

		log.Printf("✅ [correlation_id=%s] Validation completed for prescription: %s", correlationID, event.PrescriptionID)
	} else {
		log.Printf("⚠️  [correlation_id=%s] Validation failed for prescription: %s - Errors: %v", correlationID, event.PrescriptionID, validationErrors)
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
