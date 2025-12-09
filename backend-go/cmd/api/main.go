// Package main is the entry point for the Pharmonico API server
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pharmonico/backend-gogit/internal/config"
	"github.com/pharmonico/backend-gogit/internal/database"
)

func main() {
	log.Println("🚀 Starting Pharmonico API Server...")

	// Load configuration
	cfg := config.Load()
	log.Printf("📋 Configuration loaded (Environment: %s)", cfg.AppEnv)

	// Connect to MongoDB
	log.Println("🔌 Connecting to MongoDB...")
	mongoClient, err := database.ConnectMongo(cfg.MongoDBURI, "pharmonico")
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("⚠️  Error disconnecting MongoDB: %v", err)
		}
	}()
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
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}
	defer func() {
		if err := pgClient.Close(); err != nil {
			log.Printf("⚠️  Error closing PostgreSQL connection: %v", err)
		}
	}()
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
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("⚠️  Error closing Redis connection: %v", err)
		}
	}()
	log.Println("✅ Redis connected successfully")

	// TODO: Initialize router and start HTTP server
	log.Println("✅ Database connections established. API server ready to start...")

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("✅ API Server is running. Press Ctrl+C to stop.")
	<-quit
	log.Println("🛑 Shutting down API Server gracefully...")
}
