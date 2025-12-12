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

// RoutingWorker handles pharmacy routing events
type RoutingWorker struct {
	mongoClient   *database.MongoClient
	kafkaProducer kafka.Producer
}

// NewRoutingWorker creates a new routing worker
func NewRoutingWorker(mongoClient *database.MongoClient, kafkaProducer kafka.Producer) *RoutingWorker {
	return &RoutingWorker{
		mongoClient:   mongoClient,
		kafkaProducer: kafkaProducer,
	}
}

// Topic returns the Kafka topic this handler consumes from
func (w *RoutingWorker) Topic() string {
	return kafka.TopicEnrollmentCompleted
}

// Handle processes an enrollment completed event and selects a pharmacy
// 8.3.1: Consumes Kafka event (handled by worker loop)
// 8.3.2: Processes business logic (pharmacy selection)
// 8.3.3: Emits next Kafka event (pharmacy.selected)
func (w *RoutingWorker) Handle(ctx context.Context, msg *kafka.Message) error {
	// Extract correlation ID from message
	correlationID := ExtractCorrelationID(msg)

	// Parse the event payload
	var event struct {
		EventID        string `json:"event_id"`
		CorrelationID  string `json:"correlation_id,omitempty"`
		PrescriptionID string `json:"prescription_id"`
		PatientID      string `json:"patient_id"`
		EnrolledAt     string `json:"enrolled_at"`
		Timestamp      string `json:"timestamp"`
	}

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("❌ [correlation_id=%s] Failed to unmarshal routing event: %v", correlationID, err)
		return err
	}

	// Use correlation ID from event if available
	if event.CorrelationID != "" {
		correlationID = event.CorrelationID
	}

	log.Printf("📍 [correlation_id=%s] Processing pharmacy routing for prescription: %s", correlationID, event.PrescriptionID)

	// Fetch prescription to get patient location or preferences
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

	// Select a pharmacy (simplified - actual routing logic will be more sophisticated)
	pharmacyCollection := w.mongoClient.GetCollection("pharmacies")

	// Find an active pharmacy
	// Note: In production, this would use aggregation pipeline to compare
	// current_daily_count with max_prescriptions_per_day
	var pharmacy bson.M
	err = pharmacyCollection.FindOne(ctx, bson.M{"active": true}).Decode(&pharmacy)

	if err != nil {
		log.Printf("❌ No active pharmacy found: %v", err)
		return err
	}

	pharmacyID := pharmacy["_id"].(primitive.ObjectID).Hex()
	pharmacyNCPDPID := ""
	if ncpdpID, ok := pharmacy["ncpdp_id"].(string); ok {
		pharmacyNCPDPID = ncpdpID
	}

	log.Printf("🏥 Selected pharmacy: %s (NCPDP: %s)", pharmacyID, pharmacyNCPDPID)

	// Update prescription with selected pharmacy
	update := bson.M{
		"$set": bson.M{
			"pharmacy_id": pharmacyID,
			"status":      "pharmacy_selected",
			"updated_at":  time.Now(),
		},
	}

	_, err = prescriptionCollection.UpdateOne(ctx, bson.M{"_id": prescriptionID}, update)
	if err != nil {
		log.Printf("❌ Failed to update prescription with pharmacy: %v", err)
		return err
	}

	// 8.3.3: Emit pharmacy selected event
	routingEvent := CreateEvent(correlationID, event.PrescriptionID, map[string]interface{}{
		"patient_id":        event.PatientID,
		"pharmacy_id":       pharmacyID,
		"pharmacy_ncpdp_id": pharmacyNCPDPID,
		"selected_at":       time.Now().Format(time.RFC3339),
	})

	if err := PublishEvent(ctx, w.kafkaProducer, kafka.TopicPharmacySelected, event.PrescriptionID, routingEvent); err != nil {
		log.Printf("❌ [correlation_id=%s] Failed to publish pharmacy selected event: %v", correlationID, err)
		return err
	}

	log.Printf("✅ Pharmacy selected for prescription: %s", event.PrescriptionID)
	return nil
}
