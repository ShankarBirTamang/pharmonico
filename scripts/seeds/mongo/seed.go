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
	mongoURI    = "mongodb://localhost:27017"
	dbName      = "pharmonico"
	collections = "pharmacies,prescribers,patients"
)

func main() {
	log.Println("🌱 Starting MongoDB seed script...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database(dbName)

	// Seed collections
	if err := seedCollection(ctx, db, "pharmacies", "pharmacies.json"); err != nil {
		log.Fatalf("Failed to seed pharmacies: %v", err)
	}

	if err := seedCollection(ctx, db, "prescribers", "prescribers.json"); err != nil {
		log.Fatalf("Failed to seed prescribers: %v", err)
	}

	if err := seedCollection(ctx, db, "patients", "patients.json"); err != nil {
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
	filePath := filepath.Join(filepath.Dir(os.Args[0]), jsonFile)
	file, err := os.Open(filePath)
	if err != nil {
		// Try relative path
		filePath = filepath.Join("scripts/seeds/mongo", jsonFile)
		file, err = os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", jsonFile, err)
		}
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

