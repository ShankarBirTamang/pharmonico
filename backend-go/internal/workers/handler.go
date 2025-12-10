// Package workers provides worker handlers for processing Kafka events
package workers

import (
	"context"

	"github.com/pharmonico/backend-gogit/internal/kafka"
)

// Handler defines the interface for processing Kafka messages
type Handler interface {
	// Handle processes a Kafka message and returns an error if processing fails
	// The handler should handle its own business logic and optionally emit new events
	Handle(ctx context.Context, msg *kafka.Message) error

	// Topic returns the Kafka topic this handler is registered for
	Topic() string
}
