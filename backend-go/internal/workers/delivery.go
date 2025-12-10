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

// DeliveryWorker handles delivery tracking events
type DeliveryWorker struct {
	mongoClient   *database.MongoClient
	kafkaProducer kafka.Producer
}

// NewDeliveryWorker creates a new delivery tracking worker
func NewDeliveryWorker(mongoClient *database.MongoClient, kafkaProducer kafka.Producer) *DeliveryWorker {
	return &DeliveryWorker{
		mongoClient:   mongoClient,
		kafkaProducer: kafkaProducer,
	}
}

// Topic returns the Kafka topic this handler consumes from
func (w *DeliveryWorker) Topic() string {
	return kafka.TopicShipmentLabelCreated
}

// Handle processes a shipment label created event and tracks delivery
func (w *DeliveryWorker) Handle(ctx context.Context, msg *kafka.Message) error {
	// Parse the event payload
	var event struct {
		EventID        string `json:"event_id"`
		PrescriptionID string `json:"prescription_id"`
		PatientID      string `json:"patient_id"`
		TrackingNumber string `json:"tracking_number"`
		LabelURL       string `json:"label_url"`
		CreatedAt      string `json:"created_at"`
		Timestamp      string `json:"timestamp"`
	}

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("❌ Failed to unmarshal delivery event: %v", err)
		return err
	}

	log.Printf("🚚 Tracking delivery for prescription: %s (Tracking: %s)", event.PrescriptionID, event.TrackingNumber)

	// In a real implementation, this worker would:
	// 1. Poll shipping carrier API for delivery status
	// 2. Update shipment status in MongoDB
	// 3. When delivered, emit delivery completed event

	// For now, we'll simulate tracking and mark as delivered after a delay
	// In production, this would be handled by webhooks from the shipping carrier

	// Update shipment status to "in_transit"
	shipmentCollection := w.mongoClient.GetCollection("shipments")
	prescriptionID, err := primitive.ObjectIDFromHex(event.PrescriptionID)
	if err != nil {
		log.Printf("❌ Invalid prescription ID format: %s", event.PrescriptionID)
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":     "in_transit",
			"updated_at": time.Now(),
		},
	}

	_, err = shipmentCollection.UpdateOne(
		ctx,
		bson.M{"prescription_id": prescriptionID},
		update,
	)
	if err != nil {
		log.Printf("⚠️  Failed to update shipment status: %v (may not exist yet)", err)
	}

	// Update prescription status
	prescriptionCollection := w.mongoClient.GetCollection("prescriptions")
	update = bson.M{
		"$set": bson.M{
			"status":     "in_transit",
			"updated_at": time.Now(),
		},
	}

	_, err = prescriptionCollection.UpdateOne(ctx, bson.M{"_id": prescriptionID}, update)
	if err != nil {
		log.Printf("❌ Failed to update prescription status: %v", err)
		return err
	}

	// Note: In production, delivery would be triggered by a webhook from the carrier
	// For now, we'll just log that tracking has started
	// The actual delivery event would be emitted when the carrier confirms delivery

	log.Printf("✅ Delivery tracking started for prescription: %s", event.PrescriptionID)

	// In a real scenario, you would:
	// 1. Set up webhook listener for carrier delivery notifications
	// 2. When delivery is confirmed, emit the delivery completed event:
	//
	// deliveryEvent := map[string]interface{}{
	// 	"event_id":        uuid.New().String(),
	// 	"prescription_id": event.PrescriptionID,
	// 	"patient_id":      event.PatientID,
	// 	"tracking_number": event.TrackingNumber,
	// 	"delivered_at":    time.Now().Format(time.RFC3339),
	// 	"timestamp":       time.Now().Format(time.RFC3339),
	// }
	//
	// eventBytes, _ := json.Marshal(deliveryEvent)
	// w.kafkaProducer.Publish(ctx, kafka.TopicShipmentDelivered, event.PrescriptionID, eventBytes)

	return nil
}
