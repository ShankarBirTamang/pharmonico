// Package database provides PostgreSQL connection and migration management
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// PostgresClient wraps the PostgreSQL database connection
type PostgresClient struct {
	DB *sql.DB
}

// ConnectPostgres establishes a connection to PostgreSQL
func ConnectPostgres(dsn string) (*PostgresClient, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	return &PostgresClient{DB: db}, nil
}

// Close closes the PostgreSQL connection
func (pc *PostgresClient) Close() error {
	if pc.DB != nil {
		return pc.DB.Close()
	}
	return nil
}

