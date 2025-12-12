// Package handlers provides HTTP request handlers
package handlers

import (
	"github.com/phil-my-meds/backend-gogit/internal/database"
	"github.com/phil-my-meds/backend-gogit/internal/kafka"
)

// Dependencies holds all dependencies needed by handlers
type Dependencies struct {
	MongoClient   *database.MongoClient
	Postgres      *database.PostgresClient
	Redis         *database.RedisClient
	KafkaProducer kafka.Producer
}

// NewDependencies creates a new Dependencies struct
func NewDependencies(
	mongoClient *database.MongoClient,
	postgres *database.PostgresClient,
	redis *database.RedisClient,
	kafkaProducer kafka.Producer,
) *Dependencies {
	return &Dependencies{
		MongoClient:   mongoClient,
		Postgres:      postgres,
		Redis:         redis,
		KafkaProducer: kafkaProducer,
	}
}
