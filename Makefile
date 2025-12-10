.PHONY: help seed mongo-seed pg-seed

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

seed: mongo-seed pg-seed ## Seed both MongoDB and PostgreSQL databases

mongo-seed: ## Seed MongoDB with initial data (pharmacies, prescribers, patients)
	@echo "🌱 Seeding MongoDB..."
	@cd backend-go && go run ./cmd/mongo-seed

pg-seed: ## Seed PostgreSQL with initial audit logs
	@echo "🌱 Seeding PostgreSQL..."
	@cd backend-go && go run ./cmd/pg-seed

