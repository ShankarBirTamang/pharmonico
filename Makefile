.PHONY: help seed mongo-seed pg-seed dev dev-frontend dev-backend dev-full dev-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ==========================================
# Database Seeding
# ==========================================

seed: mongo-seed pg-seed ## Seed both MongoDB and PostgreSQL databases

mongo-seed: ## Seed MongoDB with initial data (pharmacies, prescribers, patients)
	@echo "🌱 Seeding MongoDB..."
	@cd backend-go && go run ./cmd/mongo-seed

pg-seed: ## Seed PostgreSQL with initial audit logs
	@echo "🌱 Seeding PostgreSQL..."
	@cd backend-go && go run ./cmd/pg-seed

# ==========================================
# Development Commands
# ==========================================

dev: ## Start hybrid: Backend (Docker) + Frontend (Native) - Recommended for daily development
	@echo "🚀 Starting hybrid development environment..."
	@echo "   Backend services: Docker Compose"
	@echo "   Frontend: Native (npm run dev)"
	@echo ""
	@docker compose up -d api worker mongodb postgres redis kafka minio kafka-ui maildev
	@echo ""
	@echo "✅ Backend services started. Starting frontend..."
	@echo "   Frontend will run on http://localhost:5173"
	@echo "   Press Ctrl+C to stop frontend, then run 'make dev-down' to stop backend"
	@echo ""
	@cd frontend-react && npm run dev

dev-frontend: ## Start frontend only (native) - Fast iteration for UI work
	@echo "🚀 Starting frontend (native)..."
	@echo "   Make sure backend services are running (use 'make dev-backend')"
	@cd frontend-react && npm run dev

dev-backend: ## Start backend services only (Docker) - For backend development
	@echo "🐳 Starting backend services (Docker)..."
	@docker compose up api worker mongodb postgres redis kafka minio kafka-ui maildev

dev-full: ## Start everything in Docker - For integration testing
	@echo "🐳 Starting full stack (Docker)..."
	@docker compose up 

dev-down: ## Stop all Docker services
	@echo "🛑 Stopping Docker services..."
	@docker compose down

backend-logs: ## View backend logs
	@echo "🔍 Viewing backend logs..."
	@docker compose logs -f api

worker-logs: ## View worker logs
	@echo "🔍 Viewing worker logs..."
	@docker compose logs -f worker

mongodb-logs: ## View MongoDB logs
	@echo "🔍 Viewing MongoDB logs..."
	@docker compose logs -f mongodb

postgres-logs: ## View PostgreSQL logs
	@echo "🔍 Viewing PostgreSQL logs..."
	@docker compose logs -f postgres

redis-logs: ## View Redis logs
	@echo "🔍 Viewing Redis logs..."
	@docker compose logs -f redis