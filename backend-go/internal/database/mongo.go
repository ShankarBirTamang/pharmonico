// Package database provides database connection and client management
package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoClient wraps the MongoDB client and database
type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// ConnectMongo establishes a connection to MongoDB
func ConnectMongo(uri string, dbName string) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(dbName)

	return &MongoClient{
		Client:   client,
		Database: db,
	}, nil
}

// Disconnect closes the MongoDB connection
func (mc *MongoClient) Disconnect(ctx context.Context) error {
	if mc.Client != nil {
		return mc.Client.Disconnect(ctx)
	}
	return nil
}

// GetCollection returns a MongoDB collection
func (mc *MongoClient) GetCollection(name string) *mongo.Collection {
	return mc.Database.Collection(name)
}
