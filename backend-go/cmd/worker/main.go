// Package main is the entry point for the PhilMyMeds Worker service
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/phil-my-meds/backend-gogit/internal/config"
	"github.com/phil-my-meds/backend-gogit/internal/workers"
)

func main() {
	log.Println("🚀 Starting PhilMyMeds Worker Service...")

	// Load configuration
	cfg := config.Load()
	log.Printf("📋 Configuration loaded (Environment: %s)", cfg.AppEnv)

	// Initialize worker with all dependencies
	worker, err := InitializeWorker(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to initialize worker: %v", err)
	}

	// Register worker handlers
	log.Println("📝 Registering worker handlers...")

	// 1. Validation worker - processes prescription intake
	validationHandler := workers.NewValidationWorker(worker.MongoClient, worker.KafkaProducer)
	worker.Registry.Register(validationHandler)

	// 2. Enrollment worker - handles patient enrollment
	enrollmentHandler := workers.NewEnrollmentWorker(worker.MongoClient, worker.KafkaProducer)
	worker.Registry.Register(enrollmentHandler)

	// 3. Routing worker - selects pharmacy for prescription
	routingHandler := workers.NewRoutingWorker(worker.MongoClient, worker.KafkaProducer)
	worker.Registry.Register(routingHandler)

	// 4. Adjudication worker - processes insurance adjudication
	adjudicationHandler := workers.NewAdjudicationWorker(worker.MongoClient, worker.KafkaProducer)
	worker.Registry.Register(adjudicationHandler)

	// 5. Payment worker - creates payment links
	paymentHandler := workers.NewPaymentWorker(worker.MongoClient, worker.KafkaProducer)
	worker.Registry.Register(paymentHandler)

	// 6. Shipping worker - creates shipping labels
	shippingHandler := workers.NewShippingWorker(worker.MongoClient, worker.KafkaProducer)
	worker.Registry.Register(shippingHandler)

	// 7. Delivery worker - tracks delivery status
	deliveryHandler := workers.NewDeliveryWorker(worker.MongoClient, worker.KafkaProducer)
	worker.Registry.Register(deliveryHandler)

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

	// Start Kafka-based worker loop in a goroutine
	// Task 1.2: Worker uses Kafka with 10-second ticker and batch processing
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
		cancel() // Cancel context to stop worker loop
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := worker.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️  Error during shutdown: %v", err)
	}

	log.Println("✅ Worker Service stopped")
}
