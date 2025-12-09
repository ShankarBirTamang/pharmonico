// Package main provides a test script to produce a test event to Kafka
package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/pharmonico/backend-gogit/internal/config"
	"github.com/pharmonico/backend-gogit/internal/kafka"
)

// TestEvent represents a test event payload
type TestEvent struct {
	EventID        string    `json:"event_id"`
	EventType      string    `json:"event_type"`
	Timestamp      time.Time `json:"timestamp"`
	TestData       string    `json:"test_data"`
	PrescriptionID string    `json:"prescription_id,omitempty"`
	PatientID      string    `json:"patient_id,omitempty"`
}

func main() {
	log.Println("🚀 Testing Kafka Producer - Producing Test Event...")
	log.Println("=" + string(make([]byte, 60)))

	cfg := config.Load()

	// Create Kafka configuration
	log.Println("\n📡 Connecting to Kafka...")
	kafkaConfig := kafka.NewConfigFromString(cfg.KafkaBrokers, "test-producer-group", "test-producer-client")
	producer := kafka.NewProducer(kafkaConfig)
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("⚠️  Warning: error closing producer: %v", err)
		}
	}()

	// Test topic - using one of the existing topics
	testTopic := kafka.TopicIntakeReceived
	log.Printf("   ✓ Using topic: %s", testTopic)

	// Create test event
	testEvent := TestEvent{
		EventID:        "test-event-" + time.Now().Format("20060102-150405"),
		EventType:      "test.verification",
		Timestamp:      time.Now(),
		TestData:       "This is a test event for Kafka verification (Task 5.4.1)",
		PrescriptionID: "test-rx-12345",
		PatientID:      "test-patient-67890",
	}

	// Serialize event to JSON
	eventJSON, err := json.Marshal(testEvent)
	if err != nil {
		log.Fatalf("❌ Failed to marshal event: %v", err)
	}

	// Produce the event
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("\n📤 Producing test event...")
	log.Printf("   Event ID: %s", testEvent.EventID)
	log.Printf("   Event Type: %s", testEvent.EventType)
	log.Printf("   Topic: %s", testTopic)
	log.Printf("   Payload: %s", string(eventJSON))

	// Use prescription ID as the message key for partitioning
	messageKey := testEvent.PrescriptionID
	if err := producer.Publish(ctx, testTopic, messageKey, eventJSON); err != nil {
		log.Fatalf("❌ Failed to produce event: %v", err)
	}

	log.Println("\n✅ Test event produced successfully!")
	log.Println("\n📊 Event Details:")
	log.Printf("   ✓ Topic: %s", testTopic)
	log.Printf("   ✓ Message Key: %s", messageKey)
	log.Printf("   ✓ Event ID: %s", testEvent.EventID)
	log.Printf("   ✓ Timestamp: %s", testEvent.Timestamp.Format(time.RFC3339))
	log.Printf("   ✓ Payload Size: %d bytes", len(eventJSON))

	log.Println("\n💡 Next Steps:")
	log.Println("   1. Check Kafka UI at http://localhost:8085")
	log.Println("   2. Navigate to the topic: " + testTopic)
	log.Println("   3. Verify the message appears in the topic")
	log.Println("   4. Run the consumer test: go run cmd/test-kafka-consumer/main.go")
	log.Println("\n" + "=" + string(make([]byte, 60)))
}
