// Package config provides configuration loading from environment variables
package config

import (
	"os"
)

// Config holds all application configuration
type Config struct {
	// Environment
	AppEnv string

	// Database connections
	MongoDBURI  string
	PostgresDSN string
	RedisURL    string

	// Kafka
	KafkaBrokers string

	// MinIO
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string

	// SMTP
	SMTPHost string
	SMTPPort string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		MongoDBURI:     getEnv("MONGODB_URI", "mongodb://localhost:27017/pharmonico"),
		PostgresDSN:    getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/pharmonico?sslmode=disable"),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		KafkaBrokers:   getEnv("KAFKA_BROKERS", "localhost:29092"),
		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		SMTPHost:       getEnv("SMTP_HOST", "localhost"),
		SMTPPort:       getEnv("SMTP_PORT", "1025"),
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
