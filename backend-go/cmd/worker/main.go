// Package main is the entry point for the Pharmonico Worker service
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pharmonico/backend-gogit/internal/config"
)

func main() {
	log.Println("🚀 Starting Pharmonico Worker Service...")

	// Load configuration
	cfg := config.Load()
	log.Printf("📋 Configuration loaded (Environment: %s)", cfg.AppEnv)

	// Initialize worker with all dependencies
	worker, err := InitializeWorker(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to initialize worker: %v", err)
	}

	// Register worker handlers here
	// Example:
	// validationHandler := handlers.NewValidationHandler(worker.MongoClient, worker.KafkaProducer)
	// worker.Registry.Register(validationHandler)
	//
	// enrollmentHandler := handlers.NewEnrollmentHandler(worker.MongoClient, worker.KafkaProducer)
	// worker.Registry.Register(enrollmentHandler)
	//
	// ... register other handlers as needed

	registeredTopics := worker.Registry.GetTopics()
	if len(registeredTopics) == 0 {
		log.Println("⚠️  Warning: No worker handlers registered. Worker will not process any messages.")
		log.Println("💡 Register handlers in cmd/worker/main.go before starting the worker.")
	} else {
		log.Printf("📝 Registered handlers for topics: %v", registeredTopics)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start worker loop in a goroutine
	workerErr := make(chan error, 1)
	go func() {
		workerErr <- worker.Start(ctx)
	}()

	log.Println("✅ Worker Service is running. Press Ctrl+C to stop.")

	// Wait for shutdown signal or worker error
	select {
	case <-quit:
		log.Println("🛑 Shutting down Worker Service gracefully...")
		cancel() // Cancel context to stop worker loop
	case err := <-workerErr:
		if err != nil {
			log.Printf("❌ Worker error: %v", err)
		}
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := worker.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️  Error during shutdown: %v", err)
	}

	log.Println("✅ Worker Service stopped")
}
