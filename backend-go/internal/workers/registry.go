// Package workers provides worker handlers for processing Kafka events
package workers

import (
	"fmt"
	"sync"
)

// Registry manages worker handlers for different Kafka topics
type Registry struct {
	handlers map[string]Handler
	mu       sync.RWMutex
}

// NewRegistry creates a new worker registry
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]Handler),
	}
}

// Register adds a handler for a specific topic
// If a handler already exists for the topic, it will be replaced
func (r *Registry) Register(handler Handler) error {
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	topic := handler.Topic()
	if topic == "" {
		return fmt.Errorf("handler topic cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.handlers[topic] = handler
	return nil
}

// GetHandler returns the handler for a given topic, or nil if not found
func (r *Registry) GetHandler(topic string) Handler {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.handlers[topic]
}

// GetTopics returns all registered topics
func (r *Registry) GetTopics() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	topics := make([]string, 0, len(r.handlers))
	for topic := range r.handlers {
		topics = append(topics, topic)
	}
	return topics
}

// HasHandler checks if a handler exists for the given topic
func (r *Registry) HasHandler(topic string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.handlers[topic]
	return exists
}
