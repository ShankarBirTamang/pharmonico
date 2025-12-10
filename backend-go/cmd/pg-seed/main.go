// Package main provides PostgreSQL seed script for initial data
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

const (
	defaultPostgresDSN = "postgres://postgres:postgres@localhost:5432/pharmonico?sslmode=disable"
)

func main() {
	log.Println("🌱 Starting PostgreSQL seed script...")

	// Get PostgreSQL DSN from environment or use default
	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		postgresDSN = defaultPostgresDSN
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	// Get the directory where seed JSON files are located
	// Try multiple possible locations relative to project root
	seedDirs := []string{
		"../../scripts/seeds/postgres", // From backend-go/cmd/pg-seed
		"../scripts/seeds/postgres",    // From backend-go/cmd
		"scripts/seeds/postgres",       // From project root
		filepath.Join(filepath.Dir(os.Args[0]), "../../../scripts/seeds/postgres"), // From compiled binary
	}

	var seedDir string
	for _, dir := range seedDirs {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			// Check if JSON files exist in this directory
			if _, err := os.Stat(filepath.Join(dir, "audit_logs.json")); err == nil {
				seedDir = dir
				break
			}
		}
	}

	if seedDir == "" {
		log.Fatalf("Could not find seed data directory. Tried: %v\nPlease run from project root or ensure scripts/seeds/postgres/ exists", seedDirs)
	}

	log.Printf("📁 Using seed data directory: %s", seedDir)

	// Seed audit_logs table
	if err := seedAuditLogs(ctx, db, filepath.Join(seedDir, "audit_logs.json")); err != nil {
		log.Fatalf("Failed to seed audit_logs: %v", err)
	}

	log.Println("✅ PostgreSQL seeding completed successfully!")
}

func seedAuditLogs(ctx context.Context, db *sql.DB, jsonFile string) error {
	log.Printf("📦 Seeding collection: audit_logs")

	// Check if table already has data
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_logs").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count documents: %w", err)
	}

	if count > 0 {
		log.Printf("⚠️  Table audit_logs already has %d records. Skipping...", count)
		return nil
	}

	// Read JSON file
	file, err := os.Open(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", jsonFile, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON array
	var records []map[string]interface{}
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Prepare insert statement (matching actual table schema)
	stmt := `
		INSERT INTO audit_logs (
			entity_type,
			entity_id,
			action,
			actor_id,
			actor_type,
			changes,
			metadata,
			ip_address,
			user_agent,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7::jsonb, $8::inet, $9, $10)
	`

	inserted := 0
	for _, record := range records {
		// Extract fields matching actual table schema
		entityType, _ := record["entity_type"].(string)
		entityID, _ := record["entity_id"].(string)
		action, _ := record["action"].(string)
		
		// Map user_id to actor_id, determine actor_type
		actorID, _ := record["user_id"].(string)
		if actorID == "" {
			actorID, _ = record["actor_id"].(string)
		}
		actorType := "user"
		if actorID == "system" || actorID == "validation_worker" || actorID == "enrollment_worker" || 
		   actorID == "routing_worker" || actorID == "adjudication_worker" || actorID == "payment_worker" ||
		   actorID == "shipping_worker" || actorID == "delivery_worker" {
			actorType = "system"
		}
		
		ipAddress, _ := record["ip_address"].(string)
		userAgent, _ := record["user_agent"].(string)

		// Convert details to changes JSONB
		changesJSON, err := json.Marshal(record["details"])
		if err != nil {
			log.Printf("⚠️  Warning: Failed to marshal details for record: %v", err)
			continue
		}

		// Create metadata JSONB (can include event_type and other metadata)
		metadata := map[string]interface{}{}
		if eventType, ok := record["event_type"].(string); ok {
			metadata["event_type"] = eventType
		}
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			log.Printf("⚠️  Warning: Failed to marshal metadata for record: %v", err)
			metadataJSON = []byte("{}")
		}

		// Parse created_at timestamp
		createdAtStr, _ := record["created_at"].(string)
		var createdAt time.Time
		if createdAtStr != "" {
			createdAt, err = time.Parse(time.RFC3339, createdAtStr)
			if err != nil {
				log.Printf("⚠️  Warning: Failed to parse created_at '%s': %v, using current time", createdAtStr, err)
				createdAt = time.Now()
			}
		} else {
			createdAt = time.Now()
		}

		// Insert record
		_, err = db.ExecContext(ctx, stmt,
			entityType,
			entityID,
			action,
			actorID,
			actorType,
			string(changesJSON),
			string(metadataJSON),
			ipAddress,
			userAgent,
			createdAt,
		)
		if err != nil {
			log.Printf("⚠️  Warning: Failed to insert record: %v", err)
			continue
		}
		inserted++
	}

	log.Printf("✅ Inserted %d records into audit_logs", inserted)
	return nil
}

