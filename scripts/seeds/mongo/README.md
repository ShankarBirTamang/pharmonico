# MongoDB Seed Scripts

This directory contains JSON seed data for MongoDB collections.

## Usage

The seed script has been moved to `backend-go/cmd/mongo-seed/main.go` to be part of the Go module.

### Run from project root:
```bash
go run ./backend-go/cmd/mongo-seed
```

### Run from backend-go directory:
```bash
cd backend-go
go run ./cmd/mongo-seed
```

## Environment Variables

- `MONGODB_URI` - MongoDB connection string (default: `mongodb://localhost:27017`)

## Seed Data Files

- `pharmacies.json` - Sample pharmacy data
- `prescribers.json` - Sample prescriber data  
- `patients.json` - Sample patient data

## Notes

- The script will skip collections that already contain data
- All collections are seeded into the `phil-my-meds` database

