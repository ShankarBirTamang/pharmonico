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

// ShippingWorker handles shipping label creation events
type ShippingWorker struct {
	mongoClient   *database.MongoClient
	kafkaProducer kafka.Producer
}

// NewShippingWorker creates a new shipping worker
func NewShippingWorker(mongoClient *database.MongoClient, kafkaProducer kafka.Producer) *ShippingWorker {
	return &ShippingWorker{
		mongoClient:   mongoClient,
		kafkaProducer: kafkaProducer,
	}
}

// Topic returns the Kafka topic this handler consumes from
func (w *ShippingWorker) Topic() string {
	return kafka.TopicPaymentCompleted
}

// Handle processes a payment completed event and creates a shipping label
func (w *ShippingWorker) Handle(ctx context.Context, msg *kafka.Message) error {
	// Parse the event payload
	var event struct {
		EventID        string  `json:"event_id"`
		PrescriptionID string  `json:"prescription_id"`
		PatientID      string  `json:"patient_id"`
		Amount         float64 `json:"amount"`
		Status         string  `json:"status"`
		CompletedAt    string  `json:"completed_at"`
		Timestamp      string  `json:"timestamp"`
	}

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("❌ Failed to unmarshal shipping event: %v", err)
		return err
	}

	log.Printf("📦 Processing shipping for prescription: %s", event.PrescriptionID)

	// Fetch prescription to get pharmacy and patient details
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

	// Create shipping label (simplified - actual implementation will integrate with Shippo)
	trackingNumber := "TRK" + uuid.New().String()[:12]
	labelURL := "https://labels.pharmonico.com/" + trackingNumber

	// Store shipment in MongoDB
	shipmentCollection := w.mongoClient.GetCollection("shipments")
	shipmentDoc := bson.M{
		"prescription_id": prescriptionID,
		"patient_id":      event.PatientID,
		"tracking_number": trackingNumber,
		"label_url":       labelURL,
		"status":          "label_created",
		"created_at":      time.Now(),
		"updated_at":      time.Now(),
	}

	_, err = shipmentCollection.InsertOne(ctx, shipmentDoc)
	if err != nil {
		log.Printf("❌ Failed to create shipment: %v", err)
		return err
	}

	// Update prescription status
	update := bson.M{
		"$set": bson.M{
			"status":     "shipped",
			"updated_at": time.Now(),
		},
	}

	_, err = prescriptionCollection.UpdateOne(ctx, bson.M{"_id": prescriptionID}, update)
	if err != nil {
		log.Printf("❌ Failed to update prescription status: %v", err)
		return err
	}

	// Emit shipment label created event
	shippingEvent := map[string]interface{}{
		"event_id":        uuid.New().String(),
		"prescription_id": event.PrescriptionID,
		"patient_id":      event.PatientID,
		"tracking_number": trackingNumber,
		"label_url":       labelURL,
		"created_at":      time.Now().Format(time.RFC3339),
		"timestamp":       time.Now().Format(time.RFC3339),
	}

	eventBytes, err := json.Marshal(shippingEvent)
	if err != nil {
		log.Printf("❌ Failed to marshal shipping event: %v", err)
		return err
	}

	if err := w.kafkaProducer.Publish(ctx, kafka.TopicShipmentLabelCreated, event.PrescriptionID, eventBytes); err != nil {
		log.Printf("❌ Failed to publish shipment label created event: %v", err)
		return err
	}

	log.Printf("✅ Shipping label created for prescription: %s (Tracking: %s)", event.PrescriptionID, trackingNumber)
	return nil
}
