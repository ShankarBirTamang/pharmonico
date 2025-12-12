// Package workers provides worker handlers for processing Kafka events
package workers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/phil-my-meds/backend-gogit/internal/kafka"
)

// EventMetadata represents common metadata for all events
type EventMetadata struct {
	EventID        string    `json:"event_id"`
	CorrelationID  string    `json:"correlation_id,omitempty"`
	PrescriptionID string    `json:"prescription_id,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
}

// ExtractCorrelationID extracts correlation ID from a Kafka message
// It looks for correlation_id in the message value, or generates a new one
func ExtractCorrelationID(msg *kafka.Message) string {
	var eventData map[string]interface{}
	if err := json.Unmarshal(msg.Value, &eventData); err == nil {
		if corrID, ok := eventData["correlation_id"].(string); ok && corrID != "" {
			return corrID
		}
	}
	// Generate new correlation ID if not found
	return uuid.New().String()
}

// CreateEvent creates a new event with proper metadata
func CreateEvent(correlationID, prescriptionID string, additionalData map[string]interface{}) map[string]interface{} {
	event := map[string]interface{}{
		"event_id":        uuid.New().String(),
		"correlation_id":  correlationID,
		"prescription_id": prescriptionID,
		"timestamp":       time.Now().Format(time.RFC3339),
	}

	// Merge additional data
	for k, v := range additionalData {
		event[k] = v
	}

	return event
}

// PublishEvent publishes an event to Kafka with proper error handling
func PublishEvent(ctx context.Context, producer kafka.Producer, topic, prescriptionID string, event map[string]interface{}) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("❌ Failed to marshal event: %v", err)
		return err
	}

	if err := producer.Publish(ctx, topic, prescriptionID, eventBytes); err != nil {
		log.Printf("❌ Failed to publish event to topic %s: %v", topic, err)
		return err
	}

	correlationID := ""
	if id, ok := event["correlation_id"].(string); ok {
		correlationID = id
	}

	log.Printf("✅ Published event: topic=%s, prescription_id=%s, correlation_id=%s", topic, prescriptionID, correlationID)
	return nil
}

// PublishToDeadLetterQueue publishes a failed message to the dead letter queue
func PublishToDeadLetterQueue(ctx context.Context, producer kafka.Producer, originalMsg *kafka.Message, errorMsg string) error {
	dlqEvent := map[string]interface{}{
		"event_id":       uuid.New().String(),
		"original_topic": originalMsg.Topic,
		"original_key":   string(originalMsg.Key),
		"original_value": string(originalMsg.Value),
		"error":          errorMsg,
		"failed_at":      time.Now().Format(time.RFC3339),
		"partition":      originalMsg.Partition,
		"offset":         originalMsg.Offset,
	}

	return PublishEvent(ctx, producer, kafka.TopicDeadLetterQueue, string(originalMsg.Key), dlqEvent)
}
