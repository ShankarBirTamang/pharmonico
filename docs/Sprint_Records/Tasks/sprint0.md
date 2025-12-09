# ✅ **SPRINT PLAN 0 — TASK LIST (Improved Architecture)**

**Goal:** Set up full development infrastructure, monorepo, core services, and baseline architecture.

---

# **TASK GROUP 1 — Monorepo & Base Project Setup**

### **Tasks**

1. Create monorepo folder structure
2. Add `backend-go/` folder
3. Add `frontend-react/` folder
4. Add `infra/` folder (docker, nginx, migrations)
5. Add initial README
6. Configure `.gitignore`

---

# **TASK GROUP 2 — Docker & Compose Setup**

### **Tasks**

1. Create Dockerfile for Go API
2. Create Dockerfile for Worker service
3. Create Dockerfile for React Frontend
4. Create `docker-compose.yml` with all services (Mongo, Postgres, Redis, Kafka, MinIO, Frontend, API)
5. Add hot reload for Go (Air/Reflex)
6. Add hot reload for React (Vite dev server)
7. Test all containers start successfully

---

# **TASK GROUP 3 — Database Setup (MongoDB + PostgreSQL)**

### **MongoDB Tasks**

1. Add MongoDB container with volume
2. Create Mongo connection driver
3. Add seed data script (test pharmacies, prescribers, sample patient)

### **PostgreSQL Tasks**

1. Add Postgres container with volume
2. Create migrations folder
3. Create base tables:

   * `audit_logs`
   * `validation_jobs`
   * `enrollment_jobs`
   * `routing_jobs`
   * `adjudication_jobs`
   * `payment_jobs`
   * `shipping_jobs`
   * `tracking_jobs`
4. Build migration runner
5. Test DB connection in Go

---

# **TASK GROUP 4 — Redis Setup**

### **Tasks**

1. Add Redis container
2. Implement Redis client wrapper
3. Implement caching placeholders:

   * magic link token store
   * pharmacy capacity store
   * rate limiter
4. Write integration test: set/get key

---

# **TASK GROUP 5 — Kafka Setup**

### **Tasks**

1. Add Kafka + UI containers (or Redpanda + console)
2. Create Kafka topics:

   * `prescription.intake.received`
   * `prescription.validation.completed`
   * `patient.enrollment.completed`
   * `pharmacy.selected`
   * `insurance.adjudication.completed`
   * `payment.link.created`
   * `payment.completed`
   * `shipment.label.created`
   * `shipment.delivered`
3. Build Kafka producer module in Go
4. Build Kafka consumer module in Go
5. Write test producer → consumer
6. Verify events appear in Kafka UI

---

# **TASK GROUP 6 — MinIO Setup**

### **Tasks**

1. Add MinIO container
2. Create buckets:

   * `insurance-cards/`
   * `shipping-labels/`
   * `ncpdp-raw/`
3. Implement MinIO upload utility
4. Test file upload + signed URL generation

---

# **TASK GROUP 7 — Go API Skeleton**

### **Tasks**

1. Initialize API service under `cmd/api`
2. Add router (chi/echo)
3. Add middleware (logging, recovery, CORS)
4. Add `/health` endpoint
5. Add `/api/v1` route group
6. Implement dependency injection
7. Load env variables
8. Connect DB, Redis, Kafka

---

# **TASK GROUP 8 — Worker Service Setup**

### **Tasks**

1. Initialize worker service under `cmd/worker`
2. Implement base worker loop structure
3. Add PostgreSQL row-level locking for job queue
4. Create placeholders for:

   * validation worker
   * enrollment worker
   * routing worker
   * adjudication worker
   * payment monitor worker
   * shipping worker
   * delivery tracking worker
5. Log job polling and idle status

---

# **TASK GROUP 9 — Logging & Observability**

### **Tasks**

1. Add structured logging (Zerolog/Zap)
2. Add correlation ID middleware
3. Implement global error handler
4. Create log format standard
5. Add Prometheus-ready counter placeholders

---

# **TASK GROUP 10 — Authentication Setup**

### **Tasks**

1. Add JWT utility functions
2. Create middleware for JWT validation
3. Add role-based access control placeholder
4. Test JWT-protected route

---

# **TASK GROUP 11 — GitHub Actions (CI/CD) Setup**

### **Tasks**

1. Create workflow for Go lint + build
2. Create workflow for React lint + build
3. Add Go unit tests runner
4. Build Docker images in CI
5. Add PR checks

---

# **TASK GROUP 12 — Seed Scripts**

### **Tasks**

1. MongoDB seed script (sample entities)
2. PostgreSQL seed script (audit logs & job tables)
3. Add Makefile commands:

   * `make seed`
   * `make mongo-seed`
   * `make pg-seed`

---

# **TASK GROUP 13 — Frontend Bootstrap**

### **Tasks**

1. Initialize React project with Vite
2. Install TailwindCSS
3. Create base layout
4. Add React Router
5. Add placeholder pages:

   * `/ops/dashboard`
   * `/enroll/:token`
6. Connect frontend → backend health endpoint

---

# **TASK GROUP 14 — Developer Tooling**

### **Tasks**

1. Add Makefile commands:

   * `make dev`
   * `make dev-down`
   * `make dev-logs`
   * `make mongo-shell`
   * `make psql`
2. Create `.env.example`
3. Add scripts under `scripts/`
4. Update README with dev instructions

---

# **TASK GROUP 15 — Architecture Documentation**

### **Tasks**

1. Create `docs/architecture.md`
2. Add system diagram (updated architecture)
3. Add Kafka event flow documentation
4. Add Redis caching strategy docs
5. Add job queue flow docs
6. Add ADRs (architecture decision records)

---

# ✔️ **End Result of Sprint 0**

By the end of Sprint 0, the system will have:

* Fully working local environment
* All infra running on Docker
* API + Workers skeleton
* Kafka, Redis, MinIO integrated
* Job queue tables ready
* Frontend bootstrapped
* CI pipeline working
* Clear architecture documentation

