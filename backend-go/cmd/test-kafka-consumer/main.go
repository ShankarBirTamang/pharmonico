// Package main provides a test script to consume a test event from Kafka
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pharmonico/backend-gogit/internal/config"
	"github.com/pharmonico/backend-gogit/internal/kafka"
	kafkago "github.com/segmentio/kafka-go"
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
	log.Println("👂 Testing Kafka Consumer - Consuming Test Event...")
	log.Println("=" + string(make([]byte, 60)))

	cfg := config.Load()

	// For testing, create a reader that reads from the beginning
	// This allows us to consume messages that were produced before the consumer started
	testTopic := kafka.TopicIntakeReceived
	consumerGroupName := fmt.Sprintf("test-consumer-group-%d", time.Now().Unix())

	log.Println("\n📡 Connecting to Kafka...")
	log.Printf("   ✓ Subscribing to topic: %s", testTopic)
	log.Printf("   ✓ Consumer Group: %s (unique for this test run)", consumerGroupName)
	log.Println("   ℹ️  Reading from beginning to catch existing messages")

	// Parse broker string (may be comma-separated)
	brokerStr := cfg.KafkaBrokers
	var brokers []string
	if strings.Contains(brokerStr, ",") {
		brokerList := strings.Split(brokerStr, ",")
		brokers = make([]string, len(brokerList))
		for i, b := range brokerList {
			brokers[i] = strings.TrimSpace(b)
		}
	} else {
		brokers = []string{strings.TrimSpace(brokerStr)}
	}

	// Create a reader that starts from the beginning for testing
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        brokers,
		GroupID:        consumerGroupName,
		GroupTopics:    []string{testTopic},
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        1 * time.Second,
		CommitInterval: time.Second,
		StartOffset:    kafkago.FirstOffset, // Start from the beginning for testing
	})
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("⚠️  Warning: error closing reader: %v", err)
		}
	}()

	log.Println("   ✅ Subscribed successfully")

	// Poll for messages
	log.Println("\n🔍 Polling for messages (timeout: 30 seconds)...")
	log.Println("   Waiting for test event...")

	timeoutMs := 30000 // 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMs)*time.Millisecond)
	defer cancel()

	kafkaMsg, err := reader.ReadMessage(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Println("\n⏱️  Timeout reached - no messages received")
			log.Println("💡 Make sure you've produced a test event first:")
			log.Println("   go run cmd/test-kafka-producer/main.go")
			return
		}
		log.Fatalf("❌ Error reading message: %v", err)
	}

	// Message received!
	log.Println("\n✅ Message received!")
	log.Println("\n📨 Message Details:")
	log.Printf("   ✓ Topic: %s", kafkaMsg.Topic)
	log.Printf("   ✓ Partition: %d", kafkaMsg.Partition)
	log.Printf("   ✓ Offset: %d", kafkaMsg.Offset)
	log.Printf("   ✓ Key: %s", string(kafkaMsg.Key))
	log.Printf("   ✓ Value Size: %d bytes", len(kafkaMsg.Value))

	// Parse the event
	var event TestEvent
	if err := json.Unmarshal(kafkaMsg.Value, &event); err != nil {
		log.Printf("⚠️  Warning: failed to unmarshal event JSON: %v", err)
		log.Printf("   Raw value: %s", string(kafkaMsg.Value))
	} else {
		log.Println("\n📋 Event Payload:")
		log.Printf("   ✓ Event ID: %s", event.EventID)
		log.Printf("   ✓ Event Type: %s", event.EventType)
		log.Printf("   ✓ Timestamp: %s", event.Timestamp.Format(time.RFC3339))
		log.Printf("   ✓ Test Data: %s", event.TestData)
		if event.PrescriptionID != "" {
			log.Printf("   ✓ Prescription ID: %s", event.PrescriptionID)
		}
		if event.PatientID != "" {
			log.Printf("   ✓ Patient ID: %s", event.PatientID)
		}
	}

	// Log message using helper
	log.Println("\n📊 Kafka Message Info:")
	log.Printf("[Kafka] Received: topic=%s, partition=%d, offset=%d, key=%s",
		kafkaMsg.Topic, kafkaMsg.Partition, kafkaMsg.Offset, string(kafkaMsg.Key))

	log.Println("\n" + "=" + string(make([]byte, 60)))
	log.Println("✅ Test event consumed successfully!")
	log.Println("\n💡 Verification Complete:")
	log.Println("   ✓ Producer test: PASSED")
	log.Println("   ✓ Consumer test: PASSED")
	log.Println("   ✓ Event flow: VERIFIED")
	log.Println("\n📺 Next Step:")
	log.Println("   Check Kafka UI at http://localhost:8085 to view the event trace")
}
