// Package main provides server initialization and setup
package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pharmonico/backend-gogit/internal/config"
	"github.com/pharmonico/backend-gogit/internal/database"
	"github.com/pharmonico/backend-gogit/internal/kafka"
)

// Server holds all the dependencies for the API server
type Server struct {
	Config        *config.Config
	MongoClient   *database.MongoClient
	Postgres      *database.PostgresClient
	Redis         *database.RedisClient
	KafkaProducer kafka.Producer
	Router        *http.Server
}

// InitializeServer sets up all database connections and returns a configured server
func InitializeServer(cfg *config.Config) (*Server, error) {
	server := &Server{
		Config: cfg,
	}

	// Connect to MongoDB
	log.Println("🔌 Connecting to MongoDB...")
	mongoClient, err := database.ConnectMongo(cfg.MongoDBURI, "pharmonico")
	if err != nil {
		return nil, err
	}
	server.MongoClient = mongoClient
	log.Println("✅ MongoDB connected successfully")

	// Create MongoDB indexes
	log.Println("📇 Creating MongoDB indexes...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := mongoClient.CreateIndexes(ctx); err != nil {
		log.Printf("⚠️  Warning: Failed to create MongoDB indexes: %v", err)
	} else {
		log.Println("✅ MongoDB indexes created successfully")
	}

	// Connect to PostgreSQL
	log.Println("🔌 Connecting to PostgreSQL...")
	pgClient, err := database.ConnectPostgres(cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}
	server.Postgres = pgClient
	log.Println("✅ PostgreSQL connected successfully")

	// Run migrations
	log.Println("🔄 Running PostgreSQL migrations...")
	if err := pgClient.RunMigrations(ctx); err != nil {
		log.Printf("⚠️  Warning: Failed to run migrations: %v", err)
	} else {
		log.Println("✅ PostgreSQL migrations completed successfully")
	}

	// Connect to Redis
	log.Println("🔌 Connecting to Redis...")
	redisClient, err := database.ConnectRedis(cfg.RedisURL)
	if err != nil {
		return nil, err
	}
	server.Redis = redisClient
	log.Println("✅ Redis connected successfully")

	// Initialize Kafka producer
	log.Println("🔌 Initializing Kafka producer...")
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	for i, broker := range brokers {
		brokers[i] = strings.TrimSpace(broker)
	}
	kafkaConfig := kafka.NewConfig(brokers, "pharmonico-api", "pharmonico-api-producer")
	kafkaProducer := kafka.NewProducer(kafkaConfig)
	server.KafkaProducer = kafkaProducer
	log.Println("✅ Kafka producer initialized successfully")

	// Setup router
	log.Println("🔧 Setting up router...")
	router := server.setupRouter()
	server.Router = &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}
	log.Println("✅ Router configured successfully")

	return server, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("🌐 Starting HTTP server on port %s...", s.Config.Port)
	if err := s.Router.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully closes all database connections and stops the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("🛑 Shutting down server connections...")

	// Shutdown HTTP server
	if s.Router != nil {
		if err := s.Router.Shutdown(ctx); err != nil {
			log.Printf("⚠️  Error shutting down HTTP server: %v", err)
		}
	}

	if s.MongoClient != nil {
		if err := s.MongoClient.Disconnect(ctx); err != nil {
			log.Printf("⚠️  Error disconnecting MongoDB: %v", err)
		}
	}

	if s.Postgres != nil {
		if err := s.Postgres.Close(); err != nil {
			log.Printf("⚠️  Error closing PostgreSQL connection: %v", err)
		}
	}

	if s.Redis != nil {
		if err := s.Redis.Close(); err != nil {
			log.Printf("⚠️  Error closing Redis connection: %v", err)
		}
	}

	if s.KafkaProducer != nil {
		if err := s.KafkaProducer.Close(); err != nil {
			log.Printf("⚠️  Error closing Kafka producer: %v", err)
		}
	}

	log.Println("✅ All connections closed")
	return nil
}
