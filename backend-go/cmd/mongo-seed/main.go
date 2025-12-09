// Package main provides MongoDB seed script for initial data
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI = "mongodb://localhost:27017"
	dbName   = "pharmonico"
)

func main() {
	log.Println("🌱 Starting MongoDB seed script...")

	// Get MongoDB URI from environment or use default
	mongoURIEnv := os.Getenv("MONGODB_URI")
	if mongoURIEnv == "" {
		mongoURIEnv = mongoURI
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURIEnv))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database(dbName)

	// Get the directory where seed JSON files are located
	// Try multiple possible locations relative to project root
	seedDirs := []string{
		"../../scripts/seeds/mongo", // From backend-go/cmd/mongo-seed
		"../scripts/seeds/mongo",    // From backend-go/cmd
		"scripts/seeds/mongo",       // From project root
		filepath.Join(filepath.Dir(os.Args[0]), "../../../scripts/seeds/mongo"), // From compiled binary
	}

	var seedDir string
	for _, dir := range seedDirs {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			// Check if JSON files exist in this directory
			if _, err := os.Stat(filepath.Join(dir, "pharmacies.json")); err == nil {
				seedDir = dir
				break
			}
		}
	}

	if seedDir == "" {
		log.Fatalf("Could not find seed data directory. Tried: %v\nPlease run from project root or ensure scripts/seeds/mongo/ exists", seedDirs)
	}

	log.Printf("📁 Using seed data directory: %s", seedDir)

	// Seed collections
	if err := seedCollection(ctx, db, "pharmacies", filepath.Join(seedDir, "pharmacies.json")); err != nil {
		log.Fatalf("Failed to seed pharmacies: %v", err)
	}

	if err := seedCollection(ctx, db, "prescribers", filepath.Join(seedDir, "prescribers.json")); err != nil {
		log.Fatalf("Failed to seed prescribers: %v", err)
	}

	if err := seedCollection(ctx, db, "patients", filepath.Join(seedDir, "patients.json")); err != nil {
		log.Fatalf("Failed to seed patients: %v", err)
	}

	log.Println("✅ MongoDB seeding completed successfully!")
}

func seedCollection(ctx context.Context, db *mongo.Database, collectionName, jsonFile string) error {
	log.Printf("📦 Seeding collection: %s", collectionName)

	collection := db.Collection(collectionName)

	// Check if collection already has data
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to count documents: %w", err)
	}

	if count > 0 {
		log.Printf("⚠️  Collection %s already has %d documents. Skipping...", collectionName, count)
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
	var documents []interface{}
	if err := json.Unmarshal(data, &documents); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Insert documents
	if len(documents) > 0 {
		result, err := collection.InsertMany(ctx, documents)
		if err != nil {
			return fmt.Errorf("failed to insert documents: %w", err)
		}
		log.Printf("✅ Inserted %d documents into %s", len(result.InsertedIDs), collectionName)
	}

	return nil
}
