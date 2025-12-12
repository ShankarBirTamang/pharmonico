// Package kafka provides Kafka producer and consumer helpers for event-driven architecture
package kafka

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

// Topic constants - all Kafka topics used in PhilMyMeds
const (
	// TopicIntakeReceived - published when a new prescription is received
	TopicIntakeReceived = "prescription.intake.received"

	// TopicValidationCompleted - published when prescription validation is completed
	TopicValidationCompleted = "prescription.validation.completed"

	// TopicEnrollmentCompleted - published when patient enrollment is done
	TopicEnrollmentCompleted = "patient.enrollment.completed"

	// TopicPharmacySelected - published when a pharmacy is selected for a prescription
	TopicPharmacySelected = "pharmacy.selected"

	// TopicAdjudicationCompleted - published when insurance adjudication is completed
	TopicAdjudicationCompleted = "insurance.adjudication.completed"

	// TopicPaymentLinkCreated - published when a payment link is created
	TopicPaymentLinkCreated = "payment.link.created"

	// TopicPaymentCompleted - published when payment is completed
	TopicPaymentCompleted = "payment.completed"

	// TopicShipmentLabelCreated - published when a shipment label is created
	TopicShipmentLabelCreated = "shipment.label.created"

	// TopicShipmentDelivered - published when shipment is delivered
	TopicShipmentDelivered = "shipment.delivered"

	// TopicDeadLetterQueue - dead letter queue for failed messages
	TopicDeadLetterQueue = "dead_letter_queue"
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
	// PollBatch retrieves up to batchSize messages (blocks up to timeoutMs)
	PollBatch(timeoutMs int, batchSize int) ([]*Message, error)
	// Commit marks the current message as processed
	Commit() error
	// Close gracefully shuts down the consumer
	Close() error
}

// Configuration Helpers
// NewConfig creates a new Kafka configuration
func NewConfig(brokers []string, consumerGroup, clientID string) *Config {
	return &Config{
		Brokers:       brokers,
		ConsumerGroup: consumerGroup,
		ClientID:      clientID,
	}
}

// NewConfigFromString creates a new Kafka configuration from a comma-separated broker string
func NewConfigFromString(brokersStr string, consumerGroup, clientID string) *Config {
	brokers := strings.Split(brokersStr, ",")
	// Trim whitespace from each broker
	for i, broker := range brokers {
		brokers[i] = strings.TrimSpace(broker)
	}
	return NewConfig(brokers, consumerGroup, clientID)
}

// LogMessage is a helper to log message details (useful for debugging)
func LogMessage(msg *Message) {
	log.Printf("[Kafka] Received: topic=%s, partition=%d, offset=%d, key=%s",
		msg.Topic, msg.Partition, msg.Offset, string(msg.Key))
}

// kafkaProducer implements the Producer interface using kafka-go
type kafkaProducer struct {
	writers map[string]*kafka.Writer
	config  *Config
}

// NewProducer creates a new Kafka producer instance
func NewProducer(config *Config) Producer {
	return &kafkaProducer{
		writers: make(map[string]*kafka.Writer),
		config:  config,
	}
}

// Publish sends a message to the specified topic
func (p *kafkaProducer) Publish(ctx context.Context, topic string, key string, value []byte) error {
	writer, exists := p.writers[topic]
	if !exists {
		writer = &kafka.Writer{
			Addr:         kafka.TCP(p.config.Brokers...),
			Topic:        topic,
			Balancer:     &kafka.Hash{}, // Use hash balancer for key-based partitioning
			WriteTimeout: 10 * time.Second,
			RequiredAcks: kafka.RequireOne, // Wait for leader acknowledgment
		}
		p.writers[topic] = writer
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}

	return writer.WriteMessages(ctx, msg)
}

// Close gracefully shuts down the producer
func (p *kafkaProducer) Close() error {
	var errs []string
	for topic, writer := range p.writers {
		if err := writer.Close(); err != nil {
			errs = append(errs, topic+": "+err.Error())
		}
	}
	if len(errs) > 0 {
		return &ProducerCloseError{Errors: errs}
	}
	return nil
}

// ProducerCloseError represents errors that occurred during producer shutdown
type ProducerCloseError struct {
	Errors []string
}

func (e *ProducerCloseError) Error() string {
	return "errors closing producers: " + strings.Join(e.Errors, "; ")
}

// kafkaConsumer implements the Consumer interface using kafka-go
type kafkaConsumer struct {
	reader *kafka.Reader
	config *Config
}

// NewConsumer creates a new Kafka consumer instance
func NewConsumer(config *Config) Consumer {
	return &kafkaConsumer{
		config: config,
	}
}

// Subscribe registers interest in the given topics
func (c *kafkaConsumer) Subscribe(topics []string) error {
	if c.reader != nil {
		if err := c.reader.Close(); err != nil {
			log.Printf("[Kafka] Warning: error closing existing reader: %v", err)
		}
	}

	// Create a reader that subscribes to multiple topics
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        c.config.Brokers,
		GroupID:        c.config.ConsumerGroup,
		GroupTopics:    topics,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        1 * time.Second,
		CommitInterval: time.Second,      // Commit offsets every second
		StartOffset:    kafka.LastOffset, // Start from the latest offset by default
	})

	log.Printf("[Kafka] Consumer subscribed to topics: %v (group: %s)", topics, c.config.ConsumerGroup)
	return nil
}

// Poll retrieves the next message (blocks up to timeoutMs)
func (c *kafkaConsumer) Poll(timeoutMs int) (*Message, error) {
	if c.reader == nil {
		return nil, &ConsumerNotSubscribedError{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMs)*time.Millisecond)
	defer cancel()

	kafkaMsg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, nil // Timeout is not an error, just no message available
		}
		return nil, err
	}

	return &Message{
		Topic:     kafkaMsg.Topic,
		Key:       kafkaMsg.Key,
		Value:     kafkaMsg.Value,
		Partition: int32(kafkaMsg.Partition),
		Offset:    kafkaMsg.Offset,
	}, nil
}

// PollBatch retrieves up to batchSize messages (blocks up to timeoutMs)
// Task 1.2.3: Batch size = 10 jobs
func (c *kafkaConsumer) PollBatch(timeoutMs int, batchSize int) ([]*Message, error) {
	if c.reader == nil {
		return nil, &ConsumerNotSubscribedError{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMs)*time.Millisecond)
	defer cancel()

	var messages []*Message
	for i := 0; i < batchSize; i++ {
		kafkaMsg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if err == context.DeadlineExceeded {
				// Timeout is not an error, just return what we have
				break
			}
			return nil, err
		}

		messages = append(messages, &Message{
			Topic:     kafkaMsg.Topic,
			Key:       kafkaMsg.Key,
			Value:     kafkaMsg.Value,
			Partition: int32(kafkaMsg.Partition),
			Offset:    kafkaMsg.Offset,
		})
	}

	return messages, nil
}

// Commit marks the current message as processed
// Note: With kafka-go, commits are automatic when using GroupID, but this method
// can be used for explicit commits if needed
func (c *kafkaConsumer) Commit() error {
	if c.reader == nil {
		return &ConsumerNotSubscribedError{}
	}
	// kafka-go automatically commits offsets when using GroupID
	// This method is kept for interface compatibility
	// If explicit commit is needed, we can implement it here
	return nil
}

// Close gracefully shuts down the consumer
func (c *kafkaConsumer) Close() error {
	if c.reader == nil {
		return nil
	}
	return c.reader.Close()
}

// ConsumerNotSubscribedError indicates that Subscribe was not called before Poll
type ConsumerNotSubscribedError struct{}

func (e *ConsumerNotSubscribedError) Error() string {
	return "consumer not subscribed to any topics; call Subscribe() first"
}
