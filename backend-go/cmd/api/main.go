// Package main is the entry point for the PhilMyMeds API server
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/phil-my-meds/backend-gogit/internal/config"
)

func main() {
	log.Println("🚀 Starting PhilMyMeds API Server...")

	// Load configuration
	cfg := config.Load()
	log.Printf("📋 Configuration loaded (Environment: %s)", cfg.AppEnv)

	// Initialize server with all dependencies
	server, err := InitializeServer(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to initialize server: %v", err)
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("❌ Failed to start HTTP server: %v", err)
		}
	}()

	log.Println("✅ API Server is running. Press Ctrl+C to stop.")
	<-quit
	log.Println("🛑 Shutting down API Server gracefully...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("⚠️  Error during shutdown: %v", err)
	}
}
