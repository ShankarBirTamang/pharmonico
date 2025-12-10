// Package workers provides worker handlers for processing Kafka events
package workers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/pharmonico/backend-gogit/internal/database"
	"github.com/pharmonico/backend-gogit/internal/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EnrollmentWorker handles patient enrollment events
type EnrollmentWorker struct {
	mongoClient   *database.MongoClient
	kafkaProducer kafka.Producer
}

// NewEnrollmentWorker creates a new enrollment worker
func NewEnrollmentWorker(mongoClient *database.MongoClient, kafkaProducer kafka.Producer) *EnrollmentWorker {
	return &EnrollmentWorker{
		mongoClient:   mongoClient,
		kafkaProducer: kafkaProducer,
	}
}

// Topic returns the Kafka topic this handler consumes from
func (w *EnrollmentWorker) Topic() string {
	return kafka.TopicValidationCompleted
}

// Handle processes a validation completed event and handles patient enrollment
// 8.3.1: Consumes Kafka event (handled by worker loop)
// 8.3.2: Processes business logic (enrollment)
// 8.3.3: Emits next Kafka event (enrollment.completed)
func (w *EnrollmentWorker) Handle(ctx context.Context, msg *kafka.Message) error {
	// Extract correlation ID from message
	correlationID := ExtractCorrelationID(msg)

	// Parse the event payload
	var event struct {
		EventID        string `json:"event_id"`
		CorrelationID  string `json:"correlation_id,omitempty"`
		PrescriptionID string `json:"prescription_id"`
		PatientID      string `json:"patient_id"`
		ValidatedAt    string `json:"validated_at"`
		Timestamp      string `json:"timestamp"`
	}

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("❌ [correlation_id=%s] Failed to unmarshal enrollment event: %v", correlationID, err)
		return err
	}

	// Use correlation ID from event if available
	if event.CorrelationID != "" {
		correlationID = event.CorrelationID
	}

	log.Printf("👤 [correlation_id=%s] Processing enrollment for patient: %s (prescription: %s)", correlationID, event.PatientID, event.PrescriptionID)

	// Check if patient is already enrolled
	patientCollection := w.mongoClient.GetCollection("patients")
	patientID, err := primitive.ObjectIDFromHex(event.PatientID)
	if err != nil {
		log.Printf("❌ Invalid patient ID format: %s", event.PatientID)
		return err
	}

	var patient bson.M
	err = patientCollection.FindOne(ctx, bson.M{"_id": patientID}).Decode(&patient)
	if err != nil {
		log.Printf("❌ Patient not found: %s", event.PatientID)
		return err
	}

	// Check enrollment status
	isEnrolled := false
	if enrolled, ok := patient["enrolled"].(bool); ok {
		isEnrolled = enrolled
	}

	// If not enrolled, mark as enrolled
	if !isEnrolled {
		update := bson.M{
			"$set": bson.M{
				"enrolled":    true,
				"enrolled_at": time.Now(),
				"updated_at":  time.Now(),
			},
		}

		_, err = patientCollection.UpdateOne(ctx, bson.M{"_id": patientID}, update)
		if err != nil {
			log.Printf("❌ Failed to update patient enrollment status: %v", err)
			return err
		}

		log.Printf("✅ Patient enrolled: %s", event.PatientID)
	} else {
		log.Printf("ℹ️  Patient already enrolled: %s", event.PatientID)
	}

	// Update prescription status
	prescriptionCollection := w.mongoClient.GetCollection("prescriptions")
	prescriptionID, err := primitive.ObjectIDFromHex(event.PrescriptionID)
	if err != nil {
		log.Printf("❌ Invalid prescription ID format: %s", event.PrescriptionID)
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":     "enrolled",
			"updated_at": time.Now(),
		},
	}

	_, err = prescriptionCollection.UpdateOne(ctx, bson.M{"_id": prescriptionID}, update)
	if err != nil {
		log.Printf("❌ Failed to update prescription status: %v", err)
		return err
	}

	// 8.3.3: Emit enrollment completed event
	enrollmentEvent := CreateEvent(correlationID, event.PrescriptionID, map[string]interface{}{
		"patient_id":  event.PatientID,
		"enrolled_at": time.Now().Format(time.RFC3339),
	})

	if err := PublishEvent(ctx, w.kafkaProducer, kafka.TopicEnrollmentCompleted, event.PrescriptionID, enrollmentEvent); err != nil {
		log.Printf("❌ [correlation_id=%s] Failed to publish enrollment completed event: %v", correlationID, err)
		return err
	}

	log.Printf("✅ [correlation_id=%s] Enrollment completed for patient: %s", correlationID, event.PatientID)
	return nil
}
