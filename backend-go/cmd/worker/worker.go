// Package main provides worker service initialization and setup
package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/pharmonico/backend-gogit/internal/config"
	"github.com/pharmonico/backend-gogit/internal/database"
	"github.com/pharmonico/backend-gogit/internal/kafka"
	"github.com/pharmonico/backend-gogit/internal/workers"
)

// Worker holds all the dependencies for the worker service
type Worker struct {
	Config        *config.Config
	MongoClient   *database.MongoClient
	Postgres      *database.PostgresClient
	Redis         *database.RedisClient
	KafkaConsumer kafka.Consumer
	KafkaProducer kafka.Producer
	Registry      *workers.Registry
}

// InitializeWorker sets up all database connections and returns a configured worker
func InitializeWorker(cfg *config.Config) (*Worker, error) {
	worker := &Worker{
		Config: cfg,
	}

	// Connect to MongoDB
	log.Println("🔌 Connecting to MongoDB...")
	mongoClient, err := database.ConnectMongo(cfg.MongoDBURI, "pharmonico")
	if err != nil {
		return nil, err
	}
	worker.MongoClient = mongoClient
	log.Println("✅ MongoDB connected successfully")

	// Connect to PostgreSQL
	log.Println("🔌 Connecting to PostgreSQL...")
	pgClient, err := database.ConnectPostgres(cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}
	worker.Postgres = pgClient
	log.Println("✅ PostgreSQL connected successfully")

	// Connect to Redis
	log.Println("🔌 Connecting to Redis...")
	redisClient, err := database.ConnectRedis(cfg.RedisURL)
	if err != nil {
		return nil, err
	}
	worker.Redis = redisClient
	log.Println("✅ Redis connected successfully")

	// Initialize Kafka consumer
	log.Println("🔌 Initializing Kafka consumer...")
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	for i, broker := range brokers {
		brokers[i] = strings.TrimSpace(broker)
	}
	kafkaConsumerConfig := kafka.NewConfig(brokers, "pharmonico-worker", "pharmonico-worker-consumer")
	kafkaConsumer := kafka.NewConsumer(kafkaConsumerConfig)
	worker.KafkaConsumer = kafkaConsumer
	log.Println("✅ Kafka consumer initialized successfully")

	// Initialize Kafka producer (for emitting events after processing)
	log.Println("🔌 Initializing Kafka producer...")
	kafkaProducerConfig := kafka.NewConfig(brokers, "pharmonico-worker", "pharmonico-worker-producer")
	kafkaProducer := kafka.NewProducer(kafkaProducerConfig)
	worker.KafkaProducer = kafkaProducer
	log.Println("✅ Kafka producer initialized successfully")

	// Initialize worker registry
	log.Println("🔧 Initializing worker registry...")
	worker.Registry = workers.NewRegistry()
	log.Println("✅ Worker registry initialized successfully")

	return worker, nil
}

// Start begins the worker loop that polls Kafka and processes messages
func (w *Worker) Start(ctx context.Context) error {
	// Get all registered topics
	topics := w.Registry.GetTopics()
	if len(topics) == 0 {
		log.Println("⚠️  Warning: No handlers registered. Worker will not process any messages.")
		return nil
	}

	log.Printf("📡 Subscribing to topics: %v", topics)
	if err := w.KafkaConsumer.Subscribe(topics); err != nil {
		return err
	}

	log.Println("🔄 Starting worker loop...")
	log.Println("✅ Worker is running. Press Ctrl+C to stop.")

	// Main worker loop
	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Worker loop stopped (context cancelled)")
			return nil
		default:
			// Poll for messages with a 1 second timeout
			msg, err := w.KafkaConsumer.Poll(1000)
			if err != nil {
				log.Printf("❌ Error polling Kafka: %v", err)
				// Continue loop to retry
				time.Sleep(1 * time.Second)
				continue
			}

			// No message available (timeout)
			if msg == nil {
				continue
			}

			// Process the message
			w.processMessage(ctx, msg)
		}
	}
}

// processMessage handles a single Kafka message by routing it to the appropriate handler
func (w *Worker) processMessage(ctx context.Context, msg *kafka.Message) {
	log.Printf("📨 Received message: topic=%s, partition=%d, offset=%d, key=%s",
		msg.Topic, msg.Partition, msg.Offset, string(msg.Key))

	// Get handler for this topic
	handler := w.Registry.GetHandler(msg.Topic)
	if handler == nil {
		log.Printf("⚠️  Warning: No handler registered for topic: %s. Skipping message.", msg.Topic)
		// Note: In production, you might want to send to dead letter queue
		return
	}

	// Process the message
	if err := handler.Handle(ctx, msg); err != nil {
		log.Printf("❌ Error processing message (topic=%s, offset=%d): %v", msg.Topic, msg.Offset, err)
		// Note: In production, you might want to implement retry logic or send to dead letter queue
		return
	}

	log.Printf("✅ Successfully processed message: topic=%s, offset=%d", msg.Topic, msg.Offset)
}

// Shutdown gracefully closes all connections
func (w *Worker) Shutdown(ctx context.Context) error {
	log.Println("🛑 Shutting down worker connections...")

	if w.KafkaConsumer != nil {
		if err := w.KafkaConsumer.Close(); err != nil {
			log.Printf("⚠️  Error closing Kafka consumer: %v", err)
		}
	}

	if w.KafkaProducer != nil {
		if err := w.KafkaProducer.Close(); err != nil {
			log.Printf("⚠️  Error closing Kafka producer: %v", err)
		}
	}

	if w.MongoClient != nil {
		if err := w.MongoClient.Disconnect(ctx); err != nil {
			log.Printf("⚠️  Error disconnecting MongoDB: %v", err)
		}
	}

	if w.Postgres != nil {
		if err := w.Postgres.Close(); err != nil {
			log.Printf("⚠️  Error closing PostgreSQL connection: %v", err)
		}
	}

	if w.Redis != nil {
		if err := w.Redis.Close(); err != nil {
			log.Printf("⚠️  Error closing Redis connection: %v", err)
		}
	}

	log.Println("✅ All connections closed")
	return nil
}
