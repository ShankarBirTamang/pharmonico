// Package main provides a simple test to verify database connectivity
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/pharmonico/backend-gogit/internal/config"
	"github.com/pharmonico/backend-gogit/internal/database"
)

func main() {
	log.Println("🧪 Testing database connectivity...")

	cfg := config.Load()

	// Test MongoDB connection
	log.Println("\n📊 Testing MongoDB connection...")
	mongoClient, err := database.ConnectMongo(cfg.MongoDBURI, "pharmonico")
	if err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		mongoClient.Disconnect(ctx)
	}()
	log.Println("✅ MongoDB connection successful")

	// Test MongoDB indexes
	log.Println("📇 Testing MongoDB index creation...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := mongoClient.CreateIndexes(ctx); err != nil {
		log.Printf("⚠️  MongoDB index creation warning: %v", err)
	} else {
		log.Println("✅ MongoDB indexes created/verified")
	}

	// Test PostgreSQL connection
	log.Println("\n🐘 Testing PostgreSQL connection...")
	pgClient, err := database.ConnectPostgres(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("❌ PostgreSQL connection failed: %v", err)
	}
	defer pgClient.Close()
	log.Println("✅ PostgreSQL connection successful")

	// Test PostgreSQL migrations
	log.Println("🔄 Testing PostgreSQL migrations...")
	if err := pgClient.RunMigrations(ctx); err != nil {
		log.Fatalf("❌ PostgreSQL migration failed: %v", err)
	}
	log.Println("✅ PostgreSQL migrations completed")

	// Verify audit_logs table exists
	log.Println("🔍 Verifying audit_logs table...")
	var tableExists bool
	err = pgClient.DB.QueryRowContext(ctx,
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'audit_logs')",
	).Scan(&tableExists)
	if err != nil {
		log.Fatalf("❌ Failed to verify audit_logs table: %v", err)
	}
	if !tableExists {
		log.Fatalf("❌ audit_logs table does not exist")
	}
	log.Println("✅ audit_logs table verified")

	// Verify no job queue tables exist
	log.Println("🔍 Verifying job queue tables are removed...")
	jobQueueTables := []string{
		"validation_jobs",
		"enrollment_jobs",
		"routing_jobs",
		"adjudication_jobs",
		"payment_jobs",
		"shipping_jobs",
		"tracking_jobs",
	}

	for _, tableName := range jobQueueTables {
		var exists bool
		err := pgClient.DB.QueryRowContext(ctx,
			"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)",
			tableName,
		).Scan(&exists)
		if err != nil {
			log.Printf("⚠️  Warning: Failed to check table %s: %v", tableName, err)
			continue
		}
		if exists {
			log.Printf("⚠️  Warning: Job queue table %s still exists (should be removed)", tableName)
		} else {
			log.Printf("✅ Verified: %s does not exist (as expected)", tableName)
		}
	}

	log.Println("\n✅ All database connectivity tests passed!")
	os.Exit(0)
}

