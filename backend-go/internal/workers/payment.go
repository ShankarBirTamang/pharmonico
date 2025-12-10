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

// PaymentWorker handles payment processing events
type PaymentWorker struct {
	mongoClient   *database.MongoClient
	kafkaProducer kafka.Producer
}

// NewPaymentWorker creates a new payment worker
func NewPaymentWorker(mongoClient *database.MongoClient, kafkaProducer kafka.Producer) *PaymentWorker {
	return &PaymentWorker{
		mongoClient:   mongoClient,
		kafkaProducer: kafkaProducer,
	}
}

// Topic returns the Kafka topic this handler consumes from
func (w *PaymentWorker) Topic() string {
	return kafka.TopicAdjudicationCompleted
}

// Handle processes an adjudication completed event and creates a payment link
func (w *PaymentWorker) Handle(ctx context.Context, msg *kafka.Message) error {
	// Parse the event payload
	var event struct {
		EventID            string                 `json:"event_id"`
		PrescriptionID     string                 `json:"prescription_id"`
		PatientID          string                 `json:"patient_id"`
		PharmacyID         string                 `json:"pharmacy_id"`
		AdjudicationResult map[string]interface{} `json:"adjudication_result"`
		AdjudicatedAt      string                 `json:"adjudicated_at"`
		Timestamp          string                 `json:"timestamp"`
	}

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("❌ Failed to unmarshal payment event: %v", err)
		return err
	}

	log.Printf("💵 Processing payment for prescription: %s", event.PrescriptionID)

	// Extract copay amount from adjudication result
	copayAmount := 0.0
	if result, ok := event.AdjudicationResult["copay_amount"].(float64); ok {
		copayAmount = result
	}

	// If copay is 0, skip payment and go directly to shipping
	if copayAmount == 0 {
		log.Printf("ℹ️  No copay required, skipping payment for prescription: %s", event.PrescriptionID)

		// Update prescription status
		prescriptionCollection := w.mongoClient.GetCollection("prescriptions")
		prescriptionID, err := primitive.ObjectIDFromHex(event.PrescriptionID)
		if err != nil {
			log.Printf("❌ Invalid prescription ID format: %s", event.PrescriptionID)
			return err
		}

		update := bson.M{
			"$set": bson.M{
				"status":     "payment_waived",
				"updated_at": time.Now(),
			},
		}

		_, err = prescriptionCollection.UpdateOne(ctx, bson.M{"_id": prescriptionID}, update)
		if err != nil {
			log.Printf("❌ Failed to update prescription status: %v", err)
			return err
		}

		// Emit payment completed event (waived)
		paymentEvent := map[string]interface{}{
			"event_id":        uuid.New().String(),
			"prescription_id": event.PrescriptionID,
			"patient_id":      event.PatientID,
			"amount":          0.0,
			"status":          "waived",
			"completed_at":    time.Now().Format(time.RFC3339),
			"timestamp":       time.Now().Format(time.RFC3339),
		}

		eventBytes, err := json.Marshal(paymentEvent)
		if err != nil {
			log.Printf("❌ Failed to marshal payment event: %v", err)
			return err
		}

		if err := w.kafkaProducer.Publish(ctx, kafka.TopicPaymentCompleted, event.PrescriptionID, eventBytes); err != nil {
			log.Printf("❌ Failed to publish payment completed event: %v", err)
			return err
		}

		return nil
	}

	// Create payment link (simplified - actual implementation will integrate with Stripe)
	paymentLinkID := uuid.New().String()
	paymentLinkURL := "https://pay.pharmonico.com/" + paymentLinkID

	// Store payment link in MongoDB
	paymentCollection := w.mongoClient.GetCollection("payments")
	paymentDoc := bson.M{
		"prescription_id":  event.PrescriptionID,
		"patient_id":       event.PatientID,
		"amount":           copayAmount,
		"payment_link_id":  paymentLinkID,
		"payment_link_url": paymentLinkURL,
		"status":           "pending",
		"created_at":       time.Now(),
		"updated_at":       time.Now(),
	}

	_, err := paymentCollection.InsertOne(ctx, paymentDoc)
	if err != nil {
		log.Printf("❌ Failed to create payment link: %v", err)
		return err
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
			"status":     "awaiting_payment",
			"updated_at": time.Now(),
		},
	}

	_, err = prescriptionCollection.UpdateOne(ctx, bson.M{"_id": prescriptionID}, update)
	if err != nil {
		log.Printf("❌ Failed to update prescription status: %v", err)
		return err
	}

	// Emit payment link created event
	paymentLinkEvent := map[string]interface{}{
		"event_id":         uuid.New().String(),
		"prescription_id":  event.PrescriptionID,
		"patient_id":       event.PatientID,
		"payment_link_id":  paymentLinkID,
		"payment_link_url": paymentLinkURL,
		"amount":           copayAmount,
		"created_at":       time.Now().Format(time.RFC3339),
		"timestamp":        time.Now().Format(time.RFC3339),
	}

	eventBytes, err := json.Marshal(paymentLinkEvent)
	if err != nil {
		log.Printf("❌ Failed to marshal payment link event: %v", err)
		return err
	}

	if err := w.kafkaProducer.Publish(ctx, kafka.TopicPaymentLinkCreated, event.PrescriptionID, eventBytes); err != nil {
		log.Printf("❌ Failed to publish payment link created event: %v", err)
		return err
	}

	log.Printf("✅ Payment link created for prescription: %s (Amount: $%.2f)", event.PrescriptionID, copayAmount)
	return nil
}
