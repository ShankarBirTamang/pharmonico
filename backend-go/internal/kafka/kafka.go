// Package kafka provides Kafka producer and consumer helpers for event-driven architecture
package kafka

import (
	"context"
	"log"
)

// Topic constants - all Kafka topics used in Pharmonico
const (
	// TopicIntakeReceived - published when a new prescription is received
	TopicIntakeReceived = "intake_received"

	// TopicValidatePrescription - published to trigger prescription validation
	TopicValidatePrescription = "validate_prescription"

	// TopicEnrollmentCompleted - published when patient enrollment is done
	TopicEnrollmentCompleted = "enrollment_completed"

	// TopicPharmacyRecommendation - published to request pharmacy routing
	TopicPharmacyRecommendation = "pharmacy_recommendation_requested"
)

// Config holds Kafka connection configuration
type Config struct {
	Brokers       []string // List of Kafka broker addresses
	ConsumerGroup string   // Consumer group ID for this service
	ClientID      string   // Unique client identifier
}

// Message represents a Kafka message
type Message struct {
	Topic     string // Topic name
	Key       []byte // Message key (used for partitioning)
	Value     []byte // Message payload (usually JSON)
	Partition int32  // Partition number
	Offset    int64  // Message offset in partition
}

// Producer interface for publishing messages to Kafka
type Producer interface {
	// Publish sends a message to the specified topic
	Publish(ctx context.Context, topic string, key string, value []byte) error
	// Close gracefully shuts down the producer
	Close() error
}

// Consumer interface for consuming messages from Kafka
type Consumer interface {
	// Subscribe registers interest in the given topics
	Subscribe(topics []string) error
	// Poll retrieves the next message (blocks up to timeoutMs)
	Poll(timeoutMs int) (*Message, error)
	// Commit marks the current message as processed
	Commit() error
	// Close gracefully shuts down the consumer
	Close() error
}

// NewConfig creates a new Kafka configuration
func NewConfig(brokers []string, consumerGroup, clientID string) *Config {
	return &Config{
		Brokers:       brokers,
		ConsumerGroup: consumerGroup,
		ClientID:      clientID,
	}
}

// LogMessage is a helper to log message details (useful for debugging)
func LogMessage(msg *Message) {
	log.Printf("[Kafka] Received: topic=%s, partition=%d, offset=%d, key=%s",
		msg.Topic, msg.Partition, msg.Offset, string(msg.Key))
}
