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

redis-cli: ## Open Redis CLI
	@echo "🔍 Opening Redis CLI..."
	@docker exec -it phil-my-meds-redis redis-cli

# ==========================================
# MongoDB Monitoring for the terminal window
# ==========================================

mongo-shell: ## Open MongoDB shell
	@echo "🔍 Opening MongoDB shell..."
	@docker exec -it phil-my-meds-mongodb mongosh phil-my-meds

mongo-monitor: ## Enable MongoDB profiler (logs all operations)
	@echo "🔍 Enabling MongoDB profiler..."
	@docker exec -it phil-my-meds-mongodb mongosh phil-my-meds --eval "db.setProfilingLevel(2); print('✅ Profiler enabled - all operations will be logged')"

mongo-monitor-slow: ## Enable profiler for slow operations only (>100ms)
	@echo "🔍 Enabling MongoDB profiler for slow operations (>100ms)..."
	@docker exec -it phil-my-meds-mongodb mongosh phil-my-meds --eval "db.setProfilingLevel(1, { slowms: 100 }); print('✅ Profiler enabled for slow operations')"

mongo-monitor-off: ## Disable MongoDB profiler
	@echo "🛑 Disabling MongoDB profiler..."
	@docker exec -it phil-my-meds-mongodb mongosh phil-my-meds --eval "db.setProfilingLevel(0); print('✅ Profiler disabled')"

mongo-watch: ## Watch profiler output in real-time (like Redis MONITOR)
	@echo "👀 Watching MongoDB operations (press Ctrl+C to stop)..."
	@docker exec -it phil-my-meds-mongodb mongosh phil-my-meds --eval "var cursor = db.system.profile.find().sort({ts: -1}).limit(1); var lastTs = null; while(true) { var doc = cursor.hasNext() ? cursor.next() : null; if(doc && doc.ts !== lastTs) { print(JSON.stringify(doc, null, 2)); lastTs = doc.ts; } sleep(500); cursor = db.system.profile.find().sort({ts: -1}).limit(1); }"

mongo-profiler-tail: ## Show recent profiler logs
	@echo "📋 Recent MongoDB operations:"
	@docker exec phil-my-meds-mongodb mongosh phil-my-meds --quiet --eval "db.system.profile.find().sort({ts: -1}).limit(10).forEach(function(op) { print('[' + op.ts + '] ' + op.op + ' on ' + op.ns + ' (' + (op.millis || 0) + 'ms)'); })"


mongo-stat: ## Show MongoDB server statistics (refreshes every 1s)
	@echo "📊 MongoDB server statistics (press Ctrl+C to stop):"
	@docker exec -it phil-my-meds-mongodb mongostat --host localhost:27017 1

mongo-top: ## Show MongoDB collection activity (refreshes every 1s)
	@echo "📊 MongoDB collection activity (press Ctrl+C to stop):"
	@docker exec -it phil-my-meds-mongodb mongotop --host localhost:27017 1


mongo-stats: ## Show MongoDB database statistics
	@echo "📊 MongoDB Database Statistics:"
	@docker exec phil-my-meds-mongodb mongosh phil-my-meds --quiet --eval "var stats = db.stats(1024*1024); print('Database: phil-my-meds'); print('Collections: ' + stats.collections); print('Data Size: ' + stats.dataSize.toFixed(2) + ' MB'); print('Storage Size: ' + stats.storageSize.toFixed(2) + ' MB'); print('Index Size: ' + stats.indexSize.toFixed(2) + ' MB');"

mongo-prescription-stats: ## Show prescription collection statistics
	@echo "📊 Prescription Collection Statistics:"
	@docker exec phil-my-meds-mongodb sh -c 'mongosh phil-my-meds --quiet --eval '"'"'print("Total Prescriptions: " + db.prescriptions.countDocuments()); print("\nBy Status:"); db.prescriptions.aggregate([{ "$$group": { "_id": "$$status", "count": { "$$sum": 1 } } }, { "$$sort": { "count": -1 } }]).forEach(function(doc) { print("  " + (doc._id || "null") + ": " + doc.count); });'"'"''


# ==========================================
# PostgreSQL Monitoring
# ==========================================
psql: ## Open PostgreSQL shell
	@echo "🔍 Opening PostgreSQL shell..."
	@docker exec -it phil-my-meds-postgres psql -U postgres -d phil-my-meds
