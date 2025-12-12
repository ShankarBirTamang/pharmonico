// Package workers provides worker handlers for processing Kafka events
package workers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/phil-my-meds/backend-gogit/internal/database"
	"github.com/phil-my-meds/backend-gogit/internal/kafka"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdjudicationWorker handles insurance adjudication events
type AdjudicationWorker struct {
	mongoClient   *database.MongoClient
	kafkaProducer kafka.Producer
}

// NewAdjudicationWorker creates a new adjudication worker
func NewAdjudicationWorker(mongoClient *database.MongoClient, kafkaProducer kafka.Producer) *AdjudicationWorker {
	return &AdjudicationWorker{
		mongoClient:   mongoClient,
		kafkaProducer: kafkaProducer,
	}
}

// Topic returns the Kafka topic this handler consumes from
func (w *AdjudicationWorker) Topic() string {
	return kafka.TopicPharmacySelected
}

// Handle processes a pharmacy selected event and performs insurance adjudication
func (w *AdjudicationWorker) Handle(ctx context.Context, msg *kafka.Message) error {
	// Parse the event payload
	var event struct {
		EventID         string `json:"event_id"`
		PrescriptionID  string `json:"prescription_id"`
		PatientID       string `json:"patient_id"`
		PharmacyID      string `json:"pharmacy_id"`
		PharmacyNCPDPID string `json:"pharmacy_ncpdp_id"`
		SelectedAt      string `json:"selected_at"`
		Timestamp       string `json:"timestamp"`
	}

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("❌ Failed to unmarshal adjudication event: %v", err)
		return err
	}

	log.Printf("💳 Processing insurance adjudication for prescription: %s", event.PrescriptionID)

	// Fetch prescription
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

	// Perform adjudication (simplified - actual adjudication will call insurance APIs)
	// For now, we'll simulate a successful adjudication
	adjudicationResult := map[string]interface{}{
		"status":         "approved",
		"copay_amount":   25.00,
		"insurance_pays": 175.00,
		"total_cost":     200.00,
		"adjudicated_at": time.Now(),
	}

	// Store adjudication result in MongoDB
	adjudicationCollection := w.mongoClient.GetCollection("adjudications")
	adjudicationDoc := bson.M{
		"prescription_id": prescriptionID,
		"patient_id":      event.PatientID,
		"pharmacy_id":     event.PharmacyID,
		"result":          adjudicationResult,
		"created_at":      time.Now(),
		"updated_at":      time.Now(),
	}

	_, err = adjudicationCollection.InsertOne(ctx, adjudicationDoc)
	if err != nil {
		log.Printf("❌ Failed to store adjudication result: %v", err)
		return err
	}

	// Update prescription status
	update := bson.M{
		"$set": bson.M{
			"status":     "adjudicated",
			"updated_at": time.Now(),
		},
	}

	_, err = prescriptionCollection.UpdateOne(ctx, bson.M{"_id": prescriptionID}, update)
	if err != nil {
		log.Printf("❌ Failed to update prescription status: %v", err)
		return err
	}

	// Emit adjudication completed event
	adjudicationEvent := map[string]interface{}{
		"event_id":            uuid.New().String(),
		"prescription_id":     event.PrescriptionID,
		"patient_id":          event.PatientID,
		"pharmacy_id":         event.PharmacyID,
		"adjudication_result": adjudicationResult,
		"adjudicated_at":      time.Now().Format(time.RFC3339),
		"timestamp":           time.Now().Format(time.RFC3339),
	}

	eventBytes, err := json.Marshal(adjudicationEvent)
	if err != nil {
		log.Printf("❌ Failed to marshal adjudication event: %v", err)
		return err
	}

	if err := w.kafkaProducer.Publish(ctx, kafka.TopicAdjudicationCompleted, event.PrescriptionID, eventBytes); err != nil {
		log.Printf("❌ Failed to publish adjudication completed event: %v", err)
		return err
	}

	log.Printf("✅ Adjudication completed for prescription: %s", event.PrescriptionID)
	return nil
}
