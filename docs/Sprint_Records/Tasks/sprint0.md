
# ✅ **SPRINT 0 

---

# **TASK GROUP 1 — Monorepo Setup**

### **1.1 Create Monorepo Structure**

* 1.1.1 Create root folder structure
* 1.1.2 Add `backend-go/` directory
* 1.1.3 Add `frontend-react/` directory
* 1.1.4 Add `infra/` directory
* 1.1.5 Add baseline README
* 1.1.6 Add `.gitignore`

---

# **TASK GROUP 2 — Docker & Compose**

### **2.1 Dockerfile Creation**

* 2.1.1 Create Go API Dockerfile
* 2.1.2 Create Worker service Dockerfile
* 2.1.3 Create React Frontend Dockerfile

### **2.2 Compose Setup**

* 2.2.1 Add MongoDB, PostgreSQL, Redis
* 2.2.2 Add Kafka + Kafka UI or Redpanda
* 2.2.3 Add MinIO
* 2.2.4 Add API, Workers, Frontend
* 2.2.5 Configure shared bridge network

### **2.3 Developer Experience**

* 2.3.1 Add Air/Reflex for Go hot reload
* 2.3.2 Add Vite for React hot reload
* 2.3.3 Validate all containers run

---

# **TASK GROUP 3 — Database Layer**

## **3.1 MongoDB**

* 3.1.1 Add MongoDB container
* 3.1.2 Implement Mongo driver
* 3.1.3 Add indexes for collections
* 3.1.4 Add seed script for:

  * pharmacies
  * prescribers
  * patients

## **3.2 PostgreSQL (Audit Logging Only)**

* 3.2.1 Add PostgreSQL container
* 3.2.2 Create migration folder
* 3.2.3 Implement migration runner
* 3.2.4 Create table: `audit_logs`
* 3.2.5 Remove job queue tables entirely
* 3.2.6 Test connectivity from backend

---

# **TASK GROUP 4 — Redis Setup**

### **4.1 Redis Infrastructure**

* 4.1.1 Add Redis container
* 4.1.2 Implement Redis client wrapper

### **4.2 Redis Use Cases**

* 4.2.1 Magic link token store
* 4.2.2 Pharmacy capacity store
* 4.2.3 Basic rate limiter

### **4.3 Testing**

* 4.3.1 Write integration test for set/get
* 4.3.2 Validate Redis logs and TTL behavior

---

# **TASK GROUP 5 — Kafka Setup**

### **5.1 Kafka Infrastructure**

* 5.1.1 Add Kafka (or Redpanda) container
* 5.1.2 Add Kafka UI container
* 5.1.3 Configure storage volumes

### **5.2 Kafka Topics Setup**

Create topics:

* 5.2.1 `prescription.intake.received`
* 5.2.2 `prescription.validation.completed`
* 5.2.3 `patient.enrollment.completed`
* 5.2.4 `pharmacy.selected`
* 5.2.5 `insurance.adjudication.completed`
* 5.2.6 `payment.link.created`
* 5.2.7 `payment.completed`
* 5.2.8 `shipment.label.created`
* 5.2.9 `shipment.delivered`

### **5.3 Kafka Libraries**

* 5.3.1 Implement Kafka producer module
* 5.3.2 Implement Kafka consumer module
* 5.3.3 Configure consumer groups

### **5.4 Verification**

* 5.4.1 Produce test event
* 5.4.2 Consume test event
* 5.4.3 Validate event trace in Kafka UI

---

# **TASK GROUP 6 — MinIO Setup**

### **6.1 Infrastructure**

* 6.1.1 Add MinIO container
* 6.1.2 Configure root user & credentials

### **6.2 Bucket Structure**

* 6.2.1 Create bucket: `insurance-cards/`
* 6.2.2 Create bucket: `shipping-labels/`
* 6.2.3 Create bucket: `ncpdp-raw/`

### **6.3 Integration**

* 6.3.1 Implement upload utility
* 6.3.2 Implement signed URL generation
* 6.3.3 Validate upload, delete, list operations

---

# **TASK GROUP 7 — Go API Service**

### **7.1 Initialization**

* 7.1.1 Create `cmd/api` structure
* 7.1.2 Add modular folder structure
* 7.1.3 Add configuration loader

### **7.2 Routing**

* 7.2.1 Add router (chi/echo)
* 7.2.2 Add `/health` endpoint
* 7.2.3 Add `/api/v1` prefix

### **7.3 Middleware**

* 7.3.1 Logging middleware
* 7.3.2 Panic recovery
* 7.3.3 CORS
* 7.3.4 Correlation ID support

### **7.4 Integrations**

* 7.4.1 Inject MongoDB
* 7.4.2 Inject PostgreSQL
* 7.4.3 Inject Redis
* 7.4.4 Inject Kafka producer

---

# **TASK GROUP 8 — Worker Service (Kafka-Based)**

### **8.1 Worker Initialization**

* 8.1.1 Create `cmd/worker`
* 8.1.2 Add worker registry
* 8.1.3 Set up base worker loop

### **8.2 Worker Handlers**

* 8.2.1 Validation worker
* 8.2.2 Enrollment worker
* 8.2.3 Routing worker
* 8.2.4 Adjudication worker
* 8.2.5 Payment worker
* 8.2.6 Shipping worker
* 8.2.7 Delivery tracking worker

### **8.3 Event Flow**

* 8.3.1 Worker consumes Kafka event
* 8.3.2 Process business logic
* 8.3.3 Emit next Kafka event

### **8.4 Remove Old Logic**

* 8.4.1 Delete PostgreSQL polling loops
* 8.4.2 Remove job queue dependencies

---

# **TASK GROUP 9 — Logging & Observability**

### **9.1 Logging**

* 9.1.1 Integrate Zerolog/Zap
* 9.1.2 Standardize log format
* 9.1.3 Propagate correlation IDs

### **9.2 Observability**

* 9.2.1 Add Prometheus metric placeholders
* 9.2.2 Add request duration counters
* 9.2.3 Add worker processing metrics

---

# **TASK GROUP 10 — Authentication**

### **10.1 JWT Setup**

* 10.1.1 Implement token generator
* 10.1.2 Implement token validator

### **10.2 Middleware**

* 10.2.1 JWT middleware
* 10.2.2 Role-based middleware (ops_agent, ops_manager)

### **10.3 Testing**

* 10.3.1 Protect a sample route
* 10.3.2 Validate token expiration handling

---

# **TASK GROUP 11 — GitHub Actions (CI/CD)**

### **11.1 Backend CI**

* 11.1.1 Go lint
* 11.1.2 Go unit tests
* 11.1.3 Build API + Worker images

### **11.2 Frontend CI**

* 11.2.1 React lint
* 11.2.2 React build

### **11.3 PR Checks**

* 11.3.1 Add combined status checks
* 11.3.2 Enable required checks for main branch

---

# **TASK GROUP 12 — Seed Scripts**

### **12.1 MongoDB Seeds**

* 12.1.1 Generate pharmacy seed data
* 12.1.2 Add prescriber seed data
* 12.1.3 Add patient seed data

### **12.2 PostgreSQL Seed**

* 12.2.1 Add initial audit logs

### **12.3 Makefile**

* 12.3.1 `make seed`
* 12.3.2 `make mongo-seed`
* 12.3.3 `make pg-seed`

---

# **TASK GROUP 13 — Frontend Setup**

### **13.1 Initialization**

* 13.1.1 Create Vite React project
* 13.1.2 Configure TailwindCSS

### **13.2 Routing**

* 13.2.1 Add React Router
* 13.2.2 Create page: `/ops/dashboard`
* 13.2.3 Create page: `/enroll/:token`

### **13.3 Integrations**

* 13.3.1 Call backend `/health`
* 13.3.2 Add basic API client wrapper

---

# **TASK GROUP 14 — Developer Tooling**

### **14.1 Makefile**

* 14.1.1 `make dev` (start backend + workers + FE)
* 14.1.2 `make dev-down`
* 14.1.3 `make dev-logs`
* 14.1.4 `make mongo-shell`
* 14.1.5 `make psql`

### **14.2 Environment Files**

* 14.2.1 Add `.env.example`
* 14.2.2 Add environment loader

### **14.3 Scripts**

* 14.3.1 Add folder `scripts/`
* 14.3.2 Add utility scripts for seeding

---

# **TASK GROUP 15 — Architecture Documentation**

### **15.1 Core Docs**

* 15.1.1 Create `docs/architecture.md`
* 15.1.2 Document Kafka-only workflow
* 15.1.3 Document MinIO usage
* 15.1.4 Document Redis usage

### **15.2 Event Flows**

* 15.2.1 Kafka event flow diagrams
* 15.2.2 Worker processing flow

### **15.3 ADRs**

* 15.3.1 Why Kafka instead of job queue
* 15.3.2 Why PostgreSQL only for audits
* 15.3.3 Why MongoDB for domain entities

---

# ✔️ **End of Sprint 0**

Once these tasks are completed, the system has:

* Complete local development environment
* Full infrastructure (Kafka, Redis, MongoDB, PostgreSQL, MinIO)
* API and Worker skeletons
* Working CI/CD
* Documented architecture
