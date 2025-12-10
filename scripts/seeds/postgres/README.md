# PostgreSQL Seed Scripts

This directory contains JSON seed data for PostgreSQL tables.

## Usage

The seed script has been moved to `backend-go/cmd/pg-seed/main.go` to be part of the Go module.

### Run from project root:
```bash
go run ./backend-go/cmd/pg-seed
```

### Run from backend-go directory:
```bash
cd backend-go
go run ./cmd/pg-seed
```

### Using Makefile:
```bash
make pg-seed
```

## Environment Variables

- `POSTGRES_DSN` - PostgreSQL connection string (default: `postgres://postgres:postgres@localhost:5432/pharmonico?sslmode=disable`)

## Seed Data Files

- `audit_logs.json` - Sample audit log entries demonstrating various system events

## Notes

- The script will skip tables that already contain data
- All tables are seeded into the `pharmonico` database
- The audit_logs table must exist (created via migrations) before seeding

