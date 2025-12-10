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
)

func main() {
	log.Println("🚀 Starting Pharmonico API Server...")

	// Load configuration
	cfg := config.Load()
	log.Printf("📋 Configuration loaded (Environment: %s)", cfg.AppEnv)

	// Initialize server with all dependencies
	server, err := InitializeServer(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to initialize server: %v", err)
	}

	// Setup graceful shutdown
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("⚠️  Error during shutdown: %v", err)
		}
	}()

	// TODO: Initialize router and start HTTP server (Task 7.2)
	log.Println("✅ Database connections established. API server ready to start...")

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("✅ API Server is running. Press Ctrl+C to stop.")
	<-quit
	log.Println("🛑 Shutting down API Server gracefully...")
}
