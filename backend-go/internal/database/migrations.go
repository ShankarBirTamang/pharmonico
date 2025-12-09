// Package database provides migration management for PostgreSQL
package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migration represents a database migration
type Migration struct {
	Version     int
	Name        string
	SQL         string
	AppliedAt   *time.Time
}

// RunMigrations executes all pending migrations
func (pc *PostgresClient) RunMigrations(ctx context.Context) error {
	// Create migrations table if it doesn't exist
	if err := pc.createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Load all migration files
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	applied, err := pc.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if applied[migration.Version] {
			continue
		}

		if err := pc.applyMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
		}
	}

	return nil
}

// createMigrationsTable creates the migrations tracking table
func (pc *PostgresClient) createMigrationsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`
	_, err := pc.DB.ExecContext(ctx, query)
	return err
}

// getAppliedMigrations returns a map of applied migration versions
func (pc *PostgresClient) getAppliedMigrations(ctx context.Context) (map[int]bool, error) {
	query := `SELECT version FROM schema_migrations ORDER BY version`
	rows, err := pc.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// applyMigration applies a single migration
func (pc *PostgresClient) applyMigration(ctx context.Context, migration *Migration) error {
	tx, err := pc.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.ExecContext(ctx, migration.SQL); err != nil {
		return fmt.Errorf("migration SQL failed: %w", err)
	}

	// Record migration
	recordQuery := `INSERT INTO schema_migrations (version, name) VALUES ($1, $2)`
	if _, err := tx.ExecContext(ctx, recordQuery, migration.Version, migration.Name); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}

// loadMigrations loads all migration files from the embedded filesystem
func loadMigrations() ([]*Migration, error) {
	var migrations []*Migration

	err := fs.WalkDir(migrationsFS, "migrations", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		// Parse version from filename (e.g., "001_create_audit_logs.sql")
		filename := filepath.Base(path)
		parts := strings.SplitN(filename, "_", 2)
		if len(parts) < 2 {
			return fmt.Errorf("invalid migration filename format: %s", filename)
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid migration version in filename %s: %w", filename, err)
		}

		// Read migration SQL
		sql, err := migrationsFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", path, err)
		}

		name := strings.TrimSuffix(parts[1], ".sql")
		migrations = append(migrations, &Migration{
			Version: version,
			Name:    name,
			SQL:     string(sql),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

