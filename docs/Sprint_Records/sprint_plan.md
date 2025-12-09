# ✅ **Pharmonico — Improved Sprint Plan (16-Week Development Cycle)**

This plan reflects the **corrected architecture** where **pharmacy handles manufacturer program adjudication**, not our enrollment system.

**Timeline**: 16 weeks (4 months)  
**Team Size**: 1-2 developers  
**Approach**: Learning-focused, production-ready code

---

## 📋 **Sprint Overview**

| Sprint | Duration | Focus | Key Deliverables |
|--------|----------|-------|------------------|
| **Sprint 0** | Week 1-2 | Infrastructure & Foundation | Docker env, monorepo, seed data, CI |
| **Sprint 1** | Week 3-4 | Intake & Validation | NCPDP intake API, validation worker, basic dashboard |
| **Sprint 2** | Week 5-6 | Enrollment Flow | Magic links, Redis tokens, patient portal, HIPAA consent |
| **Sprint 3** | Week 7-9 | Pharmacy Routing & Program Management | Scoring engine, program catalog, ops selection UI |
| **Sprint 4** | Week 10-12 | Adjudication & Payment | Pharmacy integration, Stripe, program lookup API |
| **Sprint 5** | Week 13-15 | Shipping & Completion | Shippo integration, tracking, notifications, audit logs |
| **Sprint 6** | Week 16 | Testing & Polish | E2E tests, bug fixes, documentation |

---

# 🚀 **SPRINT 0 — Infrastructure & Foundation (Week 1-2)**

## **Goal**: Set up development environment and core infrastructure

---

### **TASK-0.1: Monorepo Structure**
**Estimate**: 4 hours

Create project structure:
```
pharmonico/
├── backend/
│   ├── cmd/
│   │   ├── api/           # Main API server
│   │   ├── worker/        # Background workers
│   │   └── migrator/      # Database migrations
│   ├── internal/
│   │   ├── models/        # Go structs
│   │   ├── handlers/      # HTTP handlers
│   │   ├── services/      # Business logic
│   │   ├── repositories/  # Database access
│   │   ├── workers/       # Worker implementations
│   │   ├── kafka/         # Kafka producers/consumers
│   │   └── middleware/    # Auth, logging, etc.
│   ├── pkg/
│   │   ├── ncpdp/         # NCPDP parser
│   │   ├── validation/    # Validation rules
│   │   └── utils/         # Utilities
│   ├── configs/           # Config files
│   └── go.mod
│
├── frontend/
│   ├── src/
│   │   ├── pages/
│   │   │   ├── dashboard/     # Ops dashboard
│   │   │   └── enrollment/    # Patient portal
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── services/          # API clients
│   │   ├── contexts/          # React contexts
│   │   └── utils/
│   ├── public/
│   └── package.json
│
├── infra/
│   ├── docker/
│   │   ├── docker-compose.yml
│   │   ├── docker-compose.dev.yml
│   │   └── docker-compose.prod.yml
│   ├── postgres/
│   │   ├── migrations/
│   │   └── seeds/
│   ├── mongodb/
│   │   └── seeds/
│   └── nginx/
│       └── nginx.conf
│
├── scripts/
│   ├── dev.sh              # Start dev environment
│   ├── seed.sh             # Seed databases
│   ├── test.sh             # Run tests
│   └── deploy.sh           # Deployment script
│
├── docs/
│   ├── api/                # API documentation
│   ├── architecture/       # Architecture diagrams
│   └── workflows/          # Workflow documentation
│
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── cd.yml
│
├── .editorconfig
├── .gitignore
├── Makefile
└── README.md
```

**Files:**
- `.editorconfig` - Editor consistency
- `.gitignore` - Git exclusions
- `Makefile` - Development commands
- `README.md` - Project documentation

---

### **TASK-0.2: Docker Compose Environment**
**Estimate**: 8 hours

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  # API Server
  api:
    build: ./backend
    container_name: pharmonico-api
    ports:
      - "8080:8080"
    environment:
      - MONGODB_URI=mongodb://mongodb:27017/pharmonico
      - POSTGRES_URI=postgres://postgres:postgres@postgres:5432/pharmonico
      - REDIS_URI=redis://redis:6379
      - KAFKA_BROKERS=kafka:9092
      - MINIO_ENDPOINT=minio:9000
      - STRIPE_SECRET_KEY=${STRIPE_SECRET_KEY}
      - SHIPPO_API_KEY=${SHIPPO_API_KEY}
    volumes:
      - ./backend:/app
    depends_on:
      - mongodb
      - postgres
      - redis
      - kafka
      - minio
    command: air # Hot reload with Air
    networks:
      - pharmonico-network

  # Worker Service
  worker:
    build: ./backend
    container_name: pharmonico-worker
    environment:
      - MONGODB_URI=mongodb://mongodb:27017/pharmonico
      - POSTGRES_URI=postgres://postgres:postgres@postgres:5432/pharmonico
      - REDIS_URI=redis://redis:6379
      - KAFKA_BROKERS=kafka:9092
    volumes:
      - ./backend:/app
    depends_on:
      - mongodb
      - postgres
      - redis
      - kafka
    command: go run cmd/worker/main.go
    networks:
      - pharmonico-network

  # React Frontend
  frontend:
    build: ./frontend
    container_name: pharmonico-frontend
    ports:
      - "3000:3000"
    environment:
      - REACT_APP_API_URL=http://localhost:8080
    volumes:
      - ./frontend:/app
      - /app/node_modules
    command: npm start
    networks:
      - pharmonico-network

  # MongoDB
  mongodb:
    image: mongo:7
    container_name: pharmonico-mongodb
    ports:
      - "27017:27017"
    environment:
      - MONGO_INITDB_DATABASE=pharmonico
    volumes:
      - mongodb-data:/data/db
      - ./infra/mongodb/seeds:/docker-entrypoint-initdb.d
    networks:
      - pharmonico-network

  # PostgreSQL
  postgres:
    image: postgres:16
    container_name: pharmonico-postgres
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_DB=pharmonico
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./infra/postgres/migrations:/docker-entrypoint-initdb.d
    networks:
      - pharmonico-network

  # Redis
  redis:
    image: redis:7-alpine
    container_name: pharmonico-redis
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes
    volumes:
      - redis-data:/data
    networks:
      - pharmonico-network

  # Kafka + Zookeeper
  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    container_name: pharmonico-zookeeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    networks:
      - pharmonico-network

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    container_name: pharmonico-kafka
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    networks:
      - pharmonico-network

  # MinIO (S3-compatible storage)
  minio:
    image: minio/minio:latest
    container_name: pharmonico-minio
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      - MINIO_ROOT_USER=minioadmin
      - MINIO_ROOT_PASSWORD=minioadmin
    volumes:
      - minio-data:/data
    command: server /data --console-address ":9001"
    networks:
      - pharmonico-network

  # Maildev (Email testing)
  maildev:
    image: maildev/maildev:latest
    container_name: pharmonico-maildev
    ports:
      - "1080:1080"  # Web UI
      - "1025:1025"  # SMTP
    networks:
      - pharmonico-network

  # Nginx (Optional reverse proxy)
  nginx:
    image: nginx:alpine
    container_name: pharmonico-nginx
    ports:
      - "80:80"
    volumes:
      - ./infra/nginx/nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - api
      - frontend
    networks:
      - pharmonico-network

volumes:
  mongodb-data:
  postgres-data:
  redis-data:
  minio-data:

networks:
  pharmonico-network:
    driver: bridge
```

**Health Checks:**
- Add health check endpoints for all services
- Implement graceful shutdown
- Add restart policies

---

### **TASK-0.3: Database Schemas & Seed Data**
**Estimate**: 6 hours

**PostgreSQL Migrations:**

`001_create_job_tables.sql`:
```sql
-- Validation Jobs
CREATE TABLE validation_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_id VARCHAR(50) UNIQUE NOT NULL,
    prescription_id VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    error_message TEXT,
    locked_at TIMESTAMP,
    locked_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

CREATE INDEX idx_validation_jobs_status ON validation_jobs(status);
CREATE INDEX idx_validation_jobs_locked ON validation_jobs(locked_at);

-- Enrollment Jobs
CREATE TABLE enrollment_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_id VARCHAR(50) UNIQUE NOT NULL,
    prescription_id VARCHAR(50) NOT NULL,
    enrollment_id VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    check_type VARCHAR(50) NOT NULL,
    next_check_at TIMESTAMP NOT NULL,
    retry_count INT DEFAULT 0,
    locked_at TIMESTAMP,
    locked_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Routing Jobs
CREATE TABLE routing_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_id VARCHAR(50) UNIQUE NOT NULL,
    prescription_id VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    retry_count INT DEFAULT 0,
    locked_at TIMESTAMP,
    locked_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Adjudication Jobs
CREATE TABLE adjudication_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_id VARCHAR(50) UNIQUE NOT NULL,
    prescription_id VARCHAR(50) NOT NULL,
    pharmacy_id VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    retry_count INT DEFAULT 0,
    locked_at TIMESTAMP,
    locked_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Payment Jobs
CREATE TABLE payment_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_id VARCHAR(50) UNIQUE NOT NULL,
    prescription_id VARCHAR(50) NOT NULL,
    payment_id VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    check_type VARCHAR(50) NOT NULL,
    next_check_at TIMESTAMP NOT NULL,
    retry_count INT DEFAULT 0,
    locked_at TIMESTAMP,
    locked_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Shipping Jobs
CREATE TABLE shipping_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_id VARCHAR(50) UNIQUE NOT NULL,
    prescription_id VARCHAR(50) NOT NULL,
    pharmacy_id VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    retry_count INT DEFAULT 0,
    locked_at TIMESTAMP,
    locked_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Tracking Jobs
CREATE TABLE tracking_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_id VARCHAR(50) UNIQUE NOT NULL,
    prescription_id VARCHAR(50) NOT NULL,
    shipment_id VARCHAR(50) NOT NULL,
    tracking_number VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    last_checked_at TIMESTAMP,
    next_check_at TIMESTAMP NOT NULL,
    retry_count INT DEFAULT 0,
    locked_at TIMESTAMP,
    locked_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Audit Logs
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    event_type VARCHAR(50) NOT NULL,
    prescription_id VARCHAR(50),
    user_id VARCHAR(50),
    action VARCHAR(100) NOT NULL,
    details JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_prescription ON audit_logs(prescription_id);
CREATE INDEX idx_audit_logs_event_type ON audit_logs(event_type);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
```

**MongoDB Seed Data:**

`seed_pharmacies.js`:
```javascript
db.pharmacies.insertMany([
  {
    _id: ObjectId(),
    name: "SpecialtyRx Pharmacy",
    license_number: "PH123456",
    npi: "9876543210",
    address: {
      line1: "123 Main St",
      city: "Boston",
      state: "MA",
      zip: "02101",
      country: "US",
      coordinates: { lat: 42.3601, lng: -71.0589 }
    },
    phone: "+1-617-555-0100",
    email: "info@specialtyrx.com",
    capabilities: {
      specialty_pharmacy: true,
      cold_chain: true,
      controlled_substances: true,
      compounding: false,
      specialty_drugs: ["00002-7510-02", "00006-0228-54"],
      max_days_supply: 90
    },
    insurance_contracts: [
      {
        payer: "Blue Cross Blue Shield",
        network_tier: "preferred",
        contract_start: new Date("2024-01-01"),
        contract_end: new Date("2024-12-31"),
        is_active: true
      }
    ],
    capacity: {
      max_daily_rx: 100,
      current_daily_rx: 0,
      max_concurrent_fills: 50,
      current_concurrent: 0
    },
    performance: {
      avg_fulfillment_time_hours: 24,
      fill_accuracy_rate: 0.998,
      customer_satisfaction: 4.8,
      on_time_delivery_rate: 0.95
    },
    status: "active",
    has_realtime_api: false,
    created_at: new Date(),
    updated_at: new Date()
  }
  // Add 4-5 more pharmacies with different characteristics
]);
```

`seed_manufacturer_programs.js`:
```javascript
db.manufacturer_programs.insertMany([
  {
    _id: ObjectId(),
    program_code: "HUMIRA_2024",
    program_name: "Humira Complete Savings Card",
    manufacturer: {
      id: "mfr_abbvie",
      name: "AbbVie Inc."
    },
    drug: {
      ndc: "00002-7510-02",
      brand_name: "Humira",
      generic_name: "adalimumab",
      generic_available: false
    },
    program_credentials: {
      bin: "004682",
      pcn: "CNRX",
      group_id: "HUMIRA"
    },
    program_type: "copay_card",
    max_annual_benefit: 16000.00,
    max_per_prescription: 2000.00,
    copay_reduction_method: "reduce_to_amount",
    target_copay: 5.00,
    eligibility_rules: {
      insurance_types_allowed: ["commercial", "marketplace"],
      insurance_types_excluded: ["medicare", "medicaid", "tricare", "va"],
      requires_commercial_primary: true,
      age_minimum: 18,
      prior_auth_acceptable: true
    },
    adjudication_endpoint: {
      type: "ncpdp_telecom",
      host: "claims.abbviecopay.com",
      port: 8080
    },
    status: "active",
    effective_dates: {
      start: new Date("2024-01-01"),
      end: new Date("2024-12-31")
    },
    terms_url: "https://humira.com/savings",
    created_at: new Date(),
    updated_at: new Date()
  }
  // Add 3-4 more programs for different drugs
]);
```

---

### **TASK-0.4: Development Scripts**
**Estimate**: 2 hours

Create `Makefile`:
```makefile
.PHONY: dev dev-down seed test clean

dev:
	docker-compose -f infra/docker/docker-compose.yml up -d
	@echo "✅ Development environment started"
	@echo "API: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@echo "Maildev: http://localhost:1080"
	@echo "MinIO: http://localhost:9001"

dev-down:
	docker-compose -f infra/docker/docker-compose.yml down

seed:
	@echo "Seeding databases..."
	docker exec pharmonico-mongodb mongosh pharmonico /docker-entrypoint-initdb.d/seed.js
	docker exec pharmonico-postgres psql -U postgres -d pharmonico -f /docker-entrypoint-initdb.d/002_seed_data.sql
	@echo "✅ Databases seeded"

test:
	cd backend && go test ./... -v -cover

clean:
	docker-compose -f infra/docker/docker-compose.yml down -v
	@echo "✅ All volumes cleaned"

logs:
	docker-compose -f infra/docker/docker-compose.yml logs -f

restart:
	docker-compose -f infra/docker/docker-compose.yml restart
```

---

### **TASK-0.5: CI/CD Pipeline**
**Estimate**: 4 hours

`.github/workflows/ci.yml`:
```yaml
name: CI Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
      
      - name: Install dependencies
        run: cd backend && go mod download
      
      - name: Run tests
        run: cd backend && go test ./... -v -cover
      
      - name: Run linter
        uses: golangci/golangci-lint-action@v3
        with:
          working-directory: backend

  test-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Cache node modules
        uses: actions/cache@v3
        with:
          path: frontend/node_modules
          key: ${{ runner.os }}-node-${{ hashFiles('**/package-lock.json') }}
      
      - name: Install dependencies
        run: cd frontend && npm ci
      
      - name: Run tests
        run: cd frontend && npm test
      
      - name: Run linter
        run: cd frontend && npm run lint

  build:
    runs-on: ubuntu-latest
    needs: [test-backend, test-frontend]
    steps:
      - uses: actions/checkout@v3
      
      - name: Build Docker images
        run: docker-compose -f infra/docker/docker-compose.yml build
```

---

## **Sprint 0 Deliverables:**

✅ Complete monorepo structure  
✅ Docker Compose environment with all services  
✅ PostgreSQL job queue tables  
✅ MongoDB collections with seed data  
✅ Development scripts (make dev, make seed)  
✅ CI/CD pipeline configured  
✅ Documentation in README

---

# 🚀 **SPRINT 1 — Intake & Validation (Week 3-4)**

## **Goal**: Implement prescription intake and automated validation

---

### **TASK-1.1: NCPDP Intake API**
**Estimate**: 12 hours

**Endpoint**: `POST /api/v1/prescriptions/intake`

**Implementation:**

1. **NCPDP Parser Package** (`pkg/ncpdp/parser.go`):
```go
package ncpdp

type NCPDPParser struct{}

func (p *NCPDPParser) Parse(xmlData []byte) (*models.Prescription, error) {
    // Parse NCPDP SCRIPT XML
    // Extract patient, prescriber, medication, insurance
    // Return structured Prescription
}
```

2. **Intake Handler** (`internal/handlers/intake.go`):
```go
func (h *Handler) IntakePrescription(c *gin.Context) {
    // 1. Parse request body (NCPDP XML or JSON)
    // 2. Check Redis for duplicate (hash of key fields)
    // 3. Parse NCPDP data
    // 4. Create prescription in MongoDB (status: "received")
    // 5. Cache in Redis (5-min TTL)
    // 6. Publish Kafka event: prescription.intake.received
    // 7. Return prescription ID
}
```

3. **Kafka Producer Integration**:
```go
func (k *KafkaProducer) PublishIntakeReceived(rx *models.Prescription) error {
    event := map[string]interface{}{
        "event_id":        uuid.New().String(),
        "prescription_id": rx.ID.Hex(),
        "patient_id":      rx.PatientID.Hex(),
        "drug_ndc":        rx.Medication.NDC,
        "timestamp":       time.Now(),
    }
    return k.Publish("prescription.intake.received", event)
}
```

**Test Cases:**
- Valid NCPDP XML intake
- Duplicate detection
- Invalid XML format
- Missing required fields
- Kafka event publishing

---

### **TASK-1.2: Validation Worker**
**Estimate**: 16 hours

**Worker Loop** (`internal/workers/validation.go`):

```go
func (w *ValidationWorker) Start(ctx context.Context) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            w.processBatch(ctx)
        case <-ctx.Done():
            return
        }
    }
}

func (w *ValidationWorker) processBatch(ctx context.Context) {
    // 1. Fetch pending jobs with row-level locking
    jobs := w.getPendingJobs(ctx, 10)
    
    // 2. Process each job concurrently
    var wg sync.WaitGroup
    for _, job := range jobs {
        wg.Add(1)
        go func(j Job) {
            defer wg.Done()
            w.processJob(ctx, j)
        }(job)
    }
    wg.Wait()
}
```

**Validation Rules** (`pkg/validation/rules.go`):

```go
type Validator interface {
    Validate(rx *models.Prescription) []ValidationError
}

// NPI Validator
type NPIValidator struct {
    registry NPIRegistry
}

func (v *NPIValidator) Validate(rx *models.Prescription) []ValidationError {
    // Check NPI format (10 digits)
    // Verify against NPI registry API
    // Return errors if invalid
}

// DEA Validator
type DEAValidator struct{}

func (v *DEAValidator) Validate(rx *models.Prescription) []ValidationError {
    // Check DEA format
    // Verify checksum algorithm
    // Check if medication requires DEA
}

// NDC Validator
type NDCValidator struct {
    fdaAPI FDAAPI
}

func (v *NDCValidator) Validate(rx *models.Prescription) []ValidationError {
    // Check NDC format
    // Verify against FDA database
    // Check drug status (active/discontinued)
}
```

**Outcome Handling:**
```go
func (w *ValidationWorker) handleValidationResult(job Job, result ValidationResult) {
    if result.IsValid {
        // Update MongoDB: status = "validated"
        // Mark job as completed
        // Publish Kafka event
        w.kafka.Publish("prescription.validation.completed", event)
    } else {
        // Update MongoDB: status = "validation_failed", errors = result.Errors
        // Mark job as failed
        // Send ops notification
    }
    
    // Log to audit_logs
    w.auditLog(job, result)
}
```

---

### **TASK-1.3: Basic Ops Dashboard**
**Estimate**: 12 hours

**React Components:**

1. **Dashboard Layout** (`src/pages/dashboard/Dashboard.tsx`):
```typescript
export const Dashboard = () => {
  return (
    <Layout>
      <Header />
      <Tabs>
        <TabPanel label="Intake" value="intake">
          <IntakeTab />
        </TabPanel>
        <TabPanel label="Validation" value="validation">
          <ValidationTab />
        </TabPanel>
        {/* More tabs in later sprints */}
      </Tabs>
    </Layout>
  );
};
```

2. **Intake Tab** (`src/pages/dashboard/IntakeTab.tsx`):
```typescript
export const IntakeTab = () => {
  const { prescriptions, loading } = usePrescriptions({ status: 'received' });
  
  return (
    <Grid>
      {prescriptions.map(rx => (
        <PrescriptionCard
          key={rx.id}
          prescription={rx}
          actions={
            <Button onClick={() => triggerValidation(rx.id)}>
              Validate
            </Button>
          }
        />
      ))}
    </Grid>
  );
};
```

3. **Validation Tab**:
```typescript
export const ValidationTab = () => {
  const { prescriptions } = usePrescriptions({ 
    status: ['validated', 'validation_failed'] 
  });
  
  return (
    <Grid>
      {prescriptions.map(rx => (
        <PrescriptionCard
          key={rx.id}
          prescription={rx}
          showValidationErrors={rx.status === 'validation_failed'}
        />
      ))}
    </Grid>
  );
};
```

4. **Real-Time Updates via WebSocket**:
```typescript
// Subscribe to Kafka events via WebSocket
const useRealtimeUpdates = () => {
  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8080/ws');
    
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      // Update local state when events arrive
      queryClient.invalidateQueries(['prescriptions']);
    };
    
    return () => ws.close();
  }, []);
};
```

---

## **Sprint 1 Deliverables:**

✅ NCPDP intake API endpoint  
✅ Validation worker with all rule checks  
✅ Kafka event flow working  
✅ Basic ops dashboard with Intake and Validation tabs  
✅ Real-time updates via WebSocket  
✅ Unit tests for validation rules  
✅ Integration test for complete flow

---

# 🚀 **SPRINT 2 — Enrollment Flow (Week 5-6)**

## **Goal**: Implement patient enrollment with magic links and HIPAA consent

---

### **TASK-2.1: Magic Link System**
**Estimate**: 8 hours

**Implementation:**

1. **Generate Magic Link** (`POST /api/v1/enrollment/initiate`):
```go
func (h *Handler) InitiateEnrollment(c *gin.Context) {
    var req struct {
        PrescriptionID string `json:"prescription_id"`
    }
    c.BindJSON(&req)
    
    // 1. Fetch prescription
    rx := h.repo.GetPrescription(req.PrescriptionID)
    
    // 2. Generate unique token
    token := uuid.New().String()
    
    // 3. Store in Redis with 48-hour TTL
    enrollmentData := map[string]interface{}{
        "prescription_id": rx.ID.Hex(),
        "patient_id":      rx.PatientID.Hex(),
        "expires_at":      time.Now().Add(48 * time.Hour),
        "used":            false,
    }
    h.redis.Set(
        ctx,
        fmt.Sprintf("magic_link:%s", token),
        enrollmentData,
        48*time.Hour,
    )
    
    // 4. Create enrollment record
    enrollment := &models.Enrollment{
        PrescriptionID: rx.ID,
        PatientID:      rx.PatientID,
        MagicLinkToken: token,
        TokenExpiresAt: time.Now().Add(48 * time.Hour),
        Status:         "pending",
    }
    h.repo.CreateEnrollment(enrollment)
    
    // 5. Send email/SMS
    magicLink := fmt.Sprintf("https://enroll.pharmonico.com/enroll/%s", token)
    h.notificationService.SendEnrollmentInvite(rx.Patient, magicLink)
    
    c.JSON(200, gin.H{"enrollment_id": enrollment.ID, "magic_link": magicLink})
}
```

2. **Validate Token** (`GET /api/v1/enrollment/validate/:token`):
```go
func (h *Handler) ValidateEnrollmentToken(c *gin.Context) {
    token := c.Param("token")
    
    // 1. Check Redis
    key := fmt.Sprintf("magic_link:%s", token)
    data, err := h.redis.Get(ctx, key).Result()
    if err == redis.Nil {
        c.JSON(404, gin.H{"error": "Token not found or expired"})
        return
    }
    
    var enrollmentData map[string]interface{}
    json.Unmarshal([]byte(data), &enrollmentData)
    
    // 2. Check if used
    if enrollmentData["used"].(bool) {
        c.JSON(400, gin.H{"error": "Token already used"})
        return
    }
    
    // 3. Check expiration
    expiresAt, _ := time.Parse(time.RFC3339, enrollmentData["expires_at"].(string))
    if time.Now().After(expiresAt) {
        c.JSON(400, gin.H{"error": "Token expired"})
        return
    }
    
    // 4. Fetch prescription data
    prescriptionID := enrollmentData["prescription_id"].(string)
    rx := h.repo.GetPrescription(prescriptionID)
    
    c.JSON(200, gin.H{
        "valid": true,
        "prescription": rx,
    })
}
```

---

### **TASK-2.2: Patient Enrollment Portal**
**Estimate**: 16 hours

**React Components:**

1. **Enrollment Flow** (`src/pages/enrollment/EnrollmentFlow.tsx`):
```typescript
export const EnrollmentFlow = () => {
  const { token } = useParams();
  const [step, setStep] = useState(1);
  
  const steps = [
    { id: 1, title: 'Insurance Information', component: <InsuranceStep /> },
    { id: 2, title: 'HIPAA Consent', component: <HIPAAConsentStep /> },
    { id: 3, title: 'Confirmation', component: <ConfirmationStep /> },
  ];
  
  return (
    <Container>
      <Stepper activeStep={step - 1}>
        {steps.map(s => (
          <Step key={s.id}>
            <StepLabel>{s.title}</StepLabel>
          </Step>
        ))}
      </Stepper>
      
      {steps[step - 1].component}
      
      <NavigationButtons
        onNext={() => setStep(step + 1)}
        onBack={() => setStep(step - 1)}
        canGoNext={validateCurrentStep()}
      />
    </Container>
  );
};
```

2. **Insurance Information Step**:
```typescript
export const InsuranceStep = () => {
  const { control, handleSubmit } = useForm<InsuranceForm>();
  
  const onSubmit = (data: InsuranceForm) => {
    // Upload insurance card images to MinIO
    const frontUrl = await uploadToMinIO(data.cardFront);
    const backUrl = await uploadToMinIO(data.cardBack);
    
    // Save to context for final submission
    updateEnrollmentData({
      insurance: {
        ...data,
        card_front_url: frontUrl,
        card_back_url: backUrl,
      },
    });
    
    nextStep();
  };
  
  return (
    <Form onSubmit={handleSubmit(onSubmit)}>
      <TextField
        name="payer_name"
        label="Insurance Company"
        control={control}
        rules={{ required: true }}
      />
      
      <TextField
        name="member_id"
        label="Member ID"
        control={control}
        rules={{ required: true }}
      />
      
      <TextField
        name="bin"
        label="BIN (6 digits)"
        control={control}
        rules={{ required: true, pattern: /^\d{6}$/ }}
      />
      
      <FileUpload
        name="cardFront"
        label="Insurance Card - Front"
        accept="image/*"
        control={control}
      />
      
      <FileUpload
        name="cardBack"
        label="Insurance Card - Back"
        accept="image/*"
        control={control}
      />
    </Form>
  );
};
```

3. **HIPAA Consent Step**:
```typescript
export const HIPAAConsentStep = () => {
  const [signature, setSignature] = useState<string>('');
  const signatureRef = useRef<SignatureCanvas>(null);
  
  const handleSubmit = () => {
    const signatureDataURL = signatureRef.current?.toDataURL();
    
    updateEnrollmentData({
      hipaa_consent: {
        authorization_text: HIPAA_TEXT,
        signature: signatureDataURL,
        signature_name: signatureName,
        signature_date: new Date().toISOString(),
        ip_address: await getClientIP(),
      },
    });
    
    nextStep();
  };
  
  return (
    <>
      <HIPAAAuthorizationText />
      
      <Checkbox
        label="I have read and agree to the HIPAA authorization"
        onChange={(checked) => setConsented(checked)}
      />
      
      <SignatureCanvas
        ref={signatureRef}
        penColor="black"
        canvasProps={{
          width: 500,
          height: 200,
          className: 'signature-canvas'
        }}
      />
      
      <Button onClick={() => signatureRef.current?.clear()}>
        Clear Signature
      </Button>
      
      <TextField
        label="Print Full Name"
        value={signatureName}
        onChange={(e) => setSignatureName(e.target.value)}
      />
      
      <Button 
        onClick={handleSubmit}
        disabled={!consented || !signature || !signatureName}
      >
        Submit Enrollment
      </Button>
    </>
  );
};
```

4. **Final Submission**:
```typescript
const submitEnrollment = async () => {
  const response = await api.post('/api/v1/enrollment/submit', {
    token,
    ...enrollmentData,
  });
  
  if (response.success) {
    showSuccess('Enrollment completed successfully!');
    navigate('/enrollment/confirmation');
  }
};
```

---

### **TASK-2.3: MinIO File Upload Integration**
**Estimate**: 4 hours

**Backend Handler**:
```go
func (h *Handler) UploadInsuranceCard(c *gin.Context) {
    file, _ := c.FormFile("file")
    prescriptionID := c.PostForm("prescription_id")
    side := c.PostForm("side") // "front" or "back"
    
    // 1. Open file
    src, _ := file.Open()
    defer src.Close()
    
    // 2. Encrypt file (HIPAA requirement)
    encryptedData := h.encryption.Encrypt(src)
    
    // 3. Upload to MinIO
    objectName := fmt.Sprintf(
        "insurance_cards/%s/%s_%s.jpg.enc",
        prescriptionID,
        side,
        uuid.New().String(),
    )
    
    info, err := h.minio.PutObject(
        ctx,
        "pharmonico-phi",
        objectName,
        encryptedData,
        -1,
        minio.PutObjectOptions{
            ContentType: "application/octet-stream",
        },
    )
    
    // 4. Store file metadata in MongoDB
    fileAsset := &models.FileAsset{
        Type:           fmt.Sprintf("insurance_card_%s", side),
        Filename:       file.Filename,
        ContentType:    file.Header.Get("Content-Type"),
        FileSize:       file.Size,
        StorageURL:     objectName,
        IsEncrypted:    true,
        EncryptionAlgo: "aes-256-gcm",
        PrescriptionID: &prescriptionID,
    }
    h.repo.CreateFileAsset(fileAsset)
    
    c.JSON(200, gin.H{"url": objectName})
}
```

---

## **Sprint 2 Deliverables:**

✅ Magic link generation and validation  
✅ Redis token storage with TTL  
✅ Patient enrollment portal (insurance + HIPAA)  
✅ Insurance card upload to MinIO  
✅ Electronic signature capture  
✅ Enrollment submission API  
✅ Email/SMS notification integration  
✅ Enrollment completed Kafka event

---

# 🚀 **SPRINT 3 — Pharmacy Routing & Program Management (Week 7-9)**

## **Goal**: Implement pharmacy scoring, program catalog, and ops selection UI

---

### **TASK-3.1: Manufacturer Program Management**
**Estimate**: 12 hours

**Admin API for Program Management:**

```go
// Create Program
POST /api/v1/admin/programs

// Update Program
PUT /api/v1/admin/programs/:id

// List Programs
GET /api/v1/admin/programs?ndc=00002-7510-02

// Activate/Deactivate Program
PATCH /api/v1/admin/programs/:id/status
```

**Program Lookup API for Pharmacy:**

```go
func (h *Handler) LookupPrograms(c *gin.Context) {
    ndc := c.Query("ndc")
    insuranceType := c.Query("insurance_type") // "commercial" or "government"
    
    // 1. Check Redis cache first
    cacheKey := fmt.Sprintf("programs:ndc:%s:%s", ndc, insuranceType)
    if cached := h.redis.Get(ctx, cacheKey).Val(); cached != "" {
        var programs []models.ManufacturerProgram
        json.Unmarshal([]byte(cached), &programs)
        c.JSON(200, programs)
        return
    }
    
    // 2. Query MongoDB
    filter := bson.M{
        "drug.ndc": ndc,
        "status": "active",
        "effective_dates.start": bson.M{"$lte": time.Now()},
        "effective_dates.end": bson.M{"$gte": time.Now()},
    }
    
    // 3. Filter by insurance type
    if insuranceType == "government" {
        // No programs for government insurance
        c.JSON(200, []models.ManufacturerProgram{})
        return
    }
    
    programs := h.repo.FindPrograms(filter)
    
    // 4. Cache for 1 hour
    programsJSON, _ := json.Marshal(programs)
    h.redis.Set(ctx, cacheKey, programsJSON, 1*time.Hour)
    
    c.JSON(200, programs)
}
```

---

### **TASK-3.2: Pharmacy Routing Worker**
**Estimate**: 16 hours

**Routing Algorithm** (`internal/services/routing.go`):

```go
func (s *RoutingService) GenerateRecommendations(rx *models.Prescription) ([]PharmacyRecommendation, error) {
    // 1. Fetch all active pharmacies
    pharmacies := s.repo.GetActivePharmacies()
    
    // 2. Apply filters
    filtered := s.filterPharmacies(pharmacies, rx)
    
    // 3. Score each pharmacy
    scored := s.scorePharmacies(filtered, rx)
    
    // 4. Sort by score (highest first)
    sort.Slice(scored, func(i, j int) bool {
        return scored[i].Score.TotalScore > scored[j].Score.TotalScore
    })
    
    // 5. Return top 5
    if len(scored) > 5 {
        scored = scored[:5]
    }
    
    return scored, nil
}

func (s *RoutingService) filterPharmacies(pharmacies []models.Pharmacy, rx *models.Prescription) []models.Pharmacy {
    filtered := []models.Pharmacy{}
    
    for _, pharmacy := range pharmacies {
        // Filter 1: Status active
        if pharmacy.Status != "active" {
            continue
        }
        
        // Filter 2: Can handle the drug
        if !contains(pharmacy.Capabilities.SpecialtyDrugs, rx.Medication.NDC) {
            continue
        }
        
        // Filter 3: Accepts patient's insurance
        if !s.acceptsInsurance(pharmacy, rx.Insurance) {
            continue
        }
        
        // Filter 4: Has capacity
        capacityUtil := float64(pharmacy.Capacity.CurrentDailyRx) / float64(pharmacy.Capacity.MaxDailyRx)
        if capacityUtil > 0.95 {
            continue
        }
        
        // Filter 5: Cold chain if needed
        if rx.Medication.RequiresColdChain && !pharmacy.Capabilities.ColdChain {
            continue
        }
        
        // Filter 6: Controlled substances if needed
        if rx.Medication.IsControlled && !pharmacy.Capabilities.ControlledSubstances {
            continue
        }
        
        filtered = append(filtered, pharmacy)
    }
    
    return filtered
}

func (s *RoutingService) scorePharmacies(pharmacies []models.Pharmacy, rx *models.Prescription) []PharmacyRecommendation {
    recommendations := []PharmacyRecommendation{}
    
    patientLat := rx.Patient.Address.Coordinates.Latitude
    patientLon := rx.Patient.Address.Coordinates.Longitude
    
    for _, pharmacy := range pharmacies {
        score := PharmacyScore{
            Breakdown: make(map[string]float64),
        }
        
        // Score 1: Geographic Distance (30%)
        distance := calculateDistance(
            patientLat, patientLon,
            pharmacy.Address.Coordinates.Latitude,
            pharmacy.Address.Coordinates.Longitude,
        )
        distanceScore := 1.0 - (distance / 100.0) // Max 100 miles
        if distanceScore < 0 {
            distanceScore = 0
        }
        score.Breakdown["distance"] = distanceScore * 0.30
        
        // Score 2: Insurance Network Tier (25%)
        networkTier := s.getNetworkTier(pharmacy, rx.Insurance)
        networkScore := 0.0
        switch networkTier {
        case "preferred":
            networkScore = 1.0
        case "standard":
            networkScore = 0.7
        case "out_of_network":
            networkScore = 0.3
        }
        score.Breakdown["network"] = networkScore * 0.25
        
        // Score 3: Current Capacity (20%)
        capacityUtil := float64(pharmacy.Capacity.CurrentDailyRx) / float64(pharmacy.Capacity.MaxDailyRx)
        capacityScore := 1.0 - capacityUtil
        score.Breakdown["capacity"] = capacityScore * 0.20
        
        // Score 4: Performance Metrics (15%)
        performanceScore := (pharmacy.Performance.FillAccuracyRate * 0.5) +
            (pharmacy.Performance.CustomerSatisfaction / 5.0 * 0.5)
        score.Breakdown["performance"] = performanceScore * 0.15
        
        // Score 5: Fulfillment Speed (10%)
        speedScore := 1.0 - (float64(pharmacy.Performance.AvgFulfillmentTimeHours) / 72.0)
        if speedScore < 0 {
            speedScore = 0
        }
        score.Breakdown["speed"] = speedScore * 0.10
        
        // Calculate total
        totalScore := 0.0
        for _, s := range score.Breakdown {
            totalScore += s
        }
        score.TotalScore = totalScore
        
        recommendations = append(recommendations, PharmacyRecommendation{
            PharmacyID:        pharmacy.ID,
            PharmacyName:      pharmacy.Name,
            Location:          pharmacy.Address,
            DistanceMiles:     distance,
            InsuranceNetwork:  networkTier,
            Score:             score,
            EstimatedFillTime: fmt.Sprintf("%d hours", int(pharmacy.Performance.AvgFulfillmentTimeHours)),
            CapacityAvailable: capacityUtil < 0.9,
        })
    }
    
    return recommendations
}
```

---

### **TASK-3.3: Pharmacy Selection UI**
**Estimate**: 12 hours

**React Component** (`src/pages/dashboard/PharmacyRoutingTab.tsx`):

```typescript
export const PharmacyRoutingTab = () => {
  const { prescriptions } = usePrescriptions({ status: 'awaiting_pharmacy_selection' });
  
  return (
    <Grid>
      {prescriptions.map(rx => (
        <PharmacySelectionCard
          key={rx.id}
          prescription={rx}
          recommendations={rx.pharmacy_recommendations}
        />
      ))}
    </Grid>
  );
};

const PharmacySelectionCard = ({ prescription, recommendations }) => {
  const [selectedPharmacy, setSelectedPharmacy] = useState(null);
  
  const handleSelect = async (pharmacyId: string) => {
    await api.post(`/api/v1/prescriptions/${prescription.id}/select-pharmacy`, {
      pharmacy_id: pharmacyId,
      selected_by: getCurrentUser().id,
      selection_reason: 'Best overall match based on scoring',
    });
    
    showSuccess('Pharmacy selected successfully!');
    refetch();
  };
  
  return (
    <Card>
      <CardHeader>
        <Typography variant="h6">
          {prescription.patient.name} - {prescription.medication.drug_name}
        </Typography>
      </CardHeader>
      
      <CardContent>
        <Typography variant="subtitle2" gutterBottom>
          Top 5 Pharmacy Recommendations
        </Typography>
        
        {recommendations.map((rec, index) => (
          <PharmacyOption
            key={rec.pharmacy_id}
            recommendation={rec}
            rank={index + 1}
            onSelect={() => handleSelect(rec.pharmacy_id)}
            isRecommended={index === 0}
          />
        ))}
      </CardContent>
    </Card>
  );
};

const PharmacyOption = ({ recommendation, rank, onSelect, isRecommended }) => {
  return (
    <Box
      sx={{
        border: isRecommended ? '2px solid #4CAF50' : '1px solid #ddd',
        borderRadius: 2,
        p: 2,
        mb: 2,
        backgroundColor: isRecommended ? '#E8F5E9' : '#fff',
      }}
    >
      <Grid container spacing={2} alignItems="center">
        <Grid item xs={1}>
          <Chip
            label={`#${rank}`}
            color={isRecommended ? 'success' : 'default'}
          />
        </Grid>
        
        <Grid item xs={5}>
          <Typography variant="h6">{recommendation.pharmacy_name}</Typography>
          <Typography variant="body2" color="textSecondary">
            {recommendation.location.city}, {recommendation.location.state} • 
            {recommendation.distance_miles.toFixed(1)} miles away
          </Typography>
        </Grid>
        
        <Grid item xs={2}>
          <Chip
            label={recommendation.insurance_network}
            color={
              recommendation.insurance_network === 'preferred'
                ? 'success'
                : 'default'
            }
          />
          {recommendation.insurance_network === 'preferred' && (
            <Typography variant="caption" display="block">
              Lower copay
            </Typography>
          )}
        </Grid>
        
        <Grid item xs={2}>
          <Typography variant="body2">
            <strong>Score: {(recommendation.score.total_score * 100).toFixed(0)}%</strong>
          </Typography>
          <Typography variant="caption">
            Est. fill: {recommendation.estimated_fill_time}
          </Typography>
        </Grid>
        
        <Grid item xs={2}>
          <Button
            variant={isRecommended ? 'contained' : 'outlined'}
            color="primary"
            fullWidth
            onClick={onSelect}
          >
            {isRecommended ? 'Select Best Match' : 'Select'}
          </Button>
        </Grid>
      </Grid>
      
      {/* Score Breakdown */}
      <Accordion sx={{ mt: 1 }}>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography variant="body2">View Score Breakdown</Typography>
        </AccordionSummary>
        <AccordionDetails>
          <Grid container spacing={1}>
            {Object.entries(recommendation.score.breakdown).map(([key, value]) => (
              <Grid item xs={12} key={key}>
                <Box display="flex" alignItems="center">
                  <Typography variant="caption" sx={{ minWidth: 100 }}>
                    {key.charAt(0).toUpperCase() + key.slice(1)}:
                  </Typography>
                  <LinearProgress
                    variant="determinate"
                    value={value * 100}
                    sx={{ flexGrow: 1, mx: 1 }}
                  />
                  <Typography variant="caption">
                    {(value * 100).toFixed(0)}%
                  </Typography>
                </Box>
              </Grid>
            ))}
          </Grid>
        </AccordionDetails>
      </Accordion>
    </Box>
  );
};
```

---

## **Sprint 3 Deliverables:**

✅ Manufacturer program management APIs  
✅ Program lookup API with Redis caching  
✅ Pharmacy routing worker with scoring algorithm  
✅ Redis capacity tracking  
✅ Pharmacy selection UI with ranked recommendations  
✅ Score breakdown visualization  
✅ Kafka event: pharmacy.selected

---

# 🚀 **SPRINT 4 — Adjudication & Payment (Week 10-12)**

## **Goal**: Coordinate pharmacy adjudication, integrate Stripe payments, handle prior authorizations

---

### **TASK-4.1: Pharmacy Adjudication API Integration**
**Estimate**: 16 hours

**Adjudication Coordinator** (`internal/services/adjudication.go`):

```go
type AdjudicationService struct {
    repo          Repository
    pharmacyAPI   PharmacyAPI
    programSvc    ProgramService
    kafka         KafkaProducer
}

func (s *AdjudicationService) ProcessAdjudication(prescriptionID string) error {
    // 1. Fetch prescription with all related data
    rx := s.repo.GetPrescriptionWithDetails(prescriptionID)
    pharmacy := s.repo.GetPharmacy(rx.SelectedPharmacyID)
    
    // 2. Prepare adjudication request
    request := s.prepareAdjudicationRequest(rx, pharmacy)
    
    // 3. Send to pharmacy API
    result, err := s.pharmacyAPI.SubmitAdjudication(pharmacy.APIEndpoint, request)
    if err != nil {
        return err
    }
    
    // 4. Handle prior authorization if needed
    if result.PriorAuthRequired {
        return s.handlePriorAuth(rx, result)
    }
    
    // 5. Store adjudication results
    adjudication := s.storeAdjudicationResults(rx, result)
    
    // 6. Publish event
    s.kafka.Publish("insurance.adjudication.completed", Event{
        PrescriptionID: rx.ID.Hex(),
        AdjudicationID: adjudication.ID.Hex(),
        FinalCopay:     adjudication.CostBreakdown.FinalPatientCopay,
        Timestamp:      time.Now(),
    })
    
    return nil
}

func (s *AdjudicationService) prepareAdjudicationRequest(rx *models.Prescription, pharmacy *models.Pharmacy) AdjudicationRequest {
    return AdjudicationRequest{
        PrescriptionID: rx.PrescriptionNumber,
        Patient: PatientData{
            FirstName: rx.Patient.FirstName,
            LastName:  rx.Patient.LastName,
            DOB:       rx.Patient.DateOfBirth,
            Gender:    rx.Patient.Sex,
        },
        Insurance: InsuranceData{
            PayerName:  rx.Insurance.PayerName,
            MemberID:   rx.Insurance.MemberID,
            GroupID:    rx.Insurance.GroupNumber,
            BIN:        rx.Insurance.BIN,
            PCN:        rx.Insurance.PCN,
        },
        Medication: MedicationData{
            NDC:        rx.Medication.NDC,
            Quantity:   rx.Medication.Quantity,
            DaysSupply: rx.Medication.DaysSupply,
            DAW:        rx.Medication.DAW,
        },
        Prescriber: PrescriberData{
            NPI:  rx.Prescriber.NPI,
            DEA:  rx.Prescriber.DEA,
            Name: fmt.Sprintf("%s %s", rx.Prescriber.FirstName, rx.Prescriber.LastName),
        },
    }
}
```

**Pharmacy API Client** (`internal/clients/pharmacy_api.go`):

```go
type PharmacyAPIClient struct {
    httpClient *http.Client
}

type AdjudicationResult struct {
    PrimaryInsurance struct {
        ClaimID         string  `json:"claim_id"`
        Status          string  `json:"status"`
        DrugCost        float64 `json:"drug_cost"`
        InsurancePaid   float64 `json:"insurance_paid"`
        PatientCopay    float64 `json:"patient_copay"`
        RejectionReason string  `json:"rejection_reason,omitempty"`
    } `json:"primary_insurance"`
    
    PriorAuthRequired bool `json:"prior_auth_required"`
    
    ManufacturerPrograms []struct {
        ProgramID       string  `json:"program_id"`
        ProgramName     string  `json:"program_name"`
        Status          string  `json:"status"`
        DiscountAmount  float64 `json:"discount_amount"`
        ReducedCopay    float64 `json:"reduced_copay"`
        RejectionReason string  `json:"rejection_reason,omitempty"`
    } `json:"manufacturer_programs"`
    
    CostBreakdown struct {
        TotalDrugCost        float64  `json:"total_drug_cost"`
        InsuranceCovered     float64  `json:"insurance_covered"`
        InitialCopay         float64  `json:"initial_copay"`
        ManufacturerDiscount float64  `json:"manufacturer_discount"`
        FinalPatientCopay    float64  `json:"final_patient_copay"`
        PatientSavings       float64  `json:"patient_savings"`
        ProgramsApplied      []string `json:"programs_applied"`
    } `json:"cost_breakdown"`
}

func (c *PharmacyAPIClient) SubmitAdjudication(endpoint string, request AdjudicationRequest) (*AdjudicationResult, error) {
    // In production, this would call the pharmacy's real API
    // For development, we'll use a mock response
    
    resp, err := c.httpClient.Post(
        fmt.Sprintf("%s/adjudication", endpoint),
        "application/json",
        bytes.NewBuffer(mustMarshal(request)),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result AdjudicationResult
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}
```

**Mock Pharmacy API for Development** (`cmd/mock-pharmacy/main.go`):

```go
// This simulates a pharmacy's adjudication system
func main() {
    r := gin.Default()
    
    r.POST("/adjudication", func(c *gin.Context) {
        var req AdjudicationRequest
        c.BindJSON(&req)
        
        // Simulate processing time
        time.Sleep(2 * time.Second)
        
        // Step 1: Primary Insurance Adjudication
        drugCost := calculateDrugCost(req.Medication.NDC)
        insurancePaid := drugCost * 0.77  // Insurance covers 77%
        patientCopay := drugCost - insurancePaid
        
        // Step 2: Check for manufacturer programs
        programs := lookupManufacturerPrograms(req.Medication.NDC, req.Insurance)
        
        var finalCopay float64
        var manufacturerDiscount float64
        var programsApplied []string
        
        if len(programs) > 0 {
            // Simulate secondary claim to manufacturer
            program := programs[0]
            manufacturerDiscount = patientCopay - 5.00 // Reduce to $5
            finalCopay = 5.00
            programsApplied = append(programsApplied, program.Name)
        } else {
            finalCopay = patientCopay
            manufacturerDiscount = 0
        }
        
        // Return results
        result := AdjudicationResult{
            PrimaryInsurance: struct{...}{
                ClaimID:       fmt.Sprintf("CLM_%s", uuid.New().String()[:8]),
                Status:        "approved",
                DrugCost:      drugCost,
                InsurancePaid: insurancePaid,
                PatientCopay:  patientCopay,
            },
            PriorAuthRequired: false,
            ManufacturerPrograms: []struct{...}{
                {
                    ProgramID:      program.ID,
                    ProgramName:    program.Name,
                    Status:         "approved",
                    DiscountAmount: manufacturerDiscount,
                    ReducedCopay:   finalCopay,
                },
            },
            CostBreakdown: struct{...}{
                TotalDrugCost:        drugCost,
                InsuranceCovered:     insurancePaid,
                InitialCopay:         patientCopay,
                ManufacturerDiscount: manufacturerDiscount,
                FinalPatientCopay:    finalCopay,
                PatientSavings:       manufacturerDiscount,
                ProgramsApplied:      programsApplied,
            },
        }
        
        c.JSON(200, result)
    })
    
    r.Run(":8090")
}
```

---

### **TASK-4.2: Prior Authorization Workflow**
**Estimate**: 12 hours

**Prior Authorization Handler** (`internal/services/prior_auth.go`):

```go
func (s *AdjudicationService) handlePriorAuth(rx *models.Prescription, result *AdjudicationResult) error {
    // 1. Create PA record
    pa := &models.PriorAuthorization{
        PrescriptionID: rx.ID,
        PayerName:      rx.Insurance.PayerName,
        RequestedAt:    time.Now(),
        RequestedBy:    "pharmacy",
        RequestReason:  result.PriorAuthRequirement.Reason,
        DiagnosisCodes: []string{}, // Would come from prescriber
        Status:         "pending",
    }
    
    paID := s.repo.CreatePriorAuth(pa)
    
    // 2. Update prescription status
    s.repo.UpdatePrescriptionStatus(rx.ID, "prior_auth_required")
    
    // 3. Notify ops team
    s.notificationSvc.SendOpsAlert(Alert{
        Type:    "prior_auth_required",
        Message: fmt.Sprintf("Prior authorization needed for Rx %s", rx.PrescriptionNumber),
        Priority: "high",
        Data: map[string]interface{}{
            "prescription_id": rx.ID.Hex(),
            "pa_id":          paID.Hex(),
            "payer":          rx.Insurance.PayerName,
            "reason":         result.PriorAuthRequirement.Reason,
        },
    })
    
    // 4. Create follow-up job
    s.createPAMonitorJob(paID)
    
    return nil
}
```

**Ops Dashboard PA Management**:

```typescript
// src/pages/dashboard/PriorAuthTab.tsx
export const PriorAuthTab = () => {
  const { priorAuths } = usePriorAuths({ status: 'pending' });
  
  return (
    <Grid>
      {priorAuths.map(pa => (
        <PriorAuthCard
          key={pa.id}
          priorAuth={pa}
          prescription={pa.prescription}
        />
      ))}
    </Grid>
  );
};

const PriorAuthCard = ({ priorAuth, prescription }) => {
  const [showDialog, setShowDialog] = useState(false);
  
  const handleSubmitPA = async (formData) => {
    await api.post(`/api/v1/prior-auth/${priorAuth.id}/submit`, {
      diagnosis_codes: formData.diagnosisCodes,
      clinical_notes: formData.clinicalNotes,
      supporting_docs: formData.documentIds,
    });
    
    showSuccess('Prior authorization submitted to insurance');
    setShowDialog(false);
  };
  
  const handleMarkApproved = async (approvalCode) => {
    await api.patch(`/api/v1/prior-auth/${priorAuth.id}`, {
      status: 'approved',
      approval_code: approvalCode,
    });
    
    // Trigger re-adjudication
    await api.post(`/api/v1/prescriptions/${prescription.id}/re-adjudicate`);
  };
  
  return (
    <Card>
      <CardHeader>
        <Typography variant="h6">
          PA Required: {prescription.medication.drug_name}
        </Typography>
        <Chip
          label={priorAuth.status.toUpperCase()}
          color={priorAuth.status === 'approved' ? 'success' : 'warning'}
        />
      </CardHeader>
      
      <CardContent>
        <Typography><strong>Patient:</strong> {prescription.patient.name}</Typography>
        <Typography><strong>Payer:</strong> {priorAuth.payer_name}</Typography>
        <Typography><strong>Reason:</strong> {priorAuth.request_reason}</Typography>
        <Typography><strong>Requested:</strong> {formatDate(priorAuth.requested_at)}</Typography>
        
        {priorAuth.status === 'pending' && (
          <Box mt={2}>
            <Button
              variant="contained"
              onClick={() => setShowDialog(true)}
            >
              Submit PA Request
            </Button>
          </Box>
        )}
        
        {priorAuth.status === 'submitted' && (
          <Box mt={2}>
            <TextField
              label="Approval Code (if received)"
              fullWidth
              onBlur={(e) => handleMarkApproved(e.target.value)}
            />
          </Box>
        )}
      </CardContent>
      
      <PASubmissionDialog
        open={showDialog}
        onClose={() => setShowDialog(false)}
        onSubmit={handleSubmitPA}
      />
    </Card>
  );
};
```

---

### **TASK-4.3: Stripe Payment Integration**
**Estimate**: 12 hours

**Payment Link Creation** (`internal/services/payment.go`):

```go
func (s *PaymentService) CreatePaymentLink(prescriptionID string) (*models.Payment, error) {
    // 1. Fetch prescription with adjudication
    rx := s.repo.GetPrescriptionWithAdjudication(prescriptionID)
    
    if rx.Adjudication == nil {
        return nil, errors.New("prescription not adjudicated")
    }
    
    finalCopay := rx.Adjudication.CostBreakdown.FinalPatientCopay
    
    // 2. Create Stripe Checkout Session
    params := &stripe.CheckoutSessionParams{
        PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
        Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
        SuccessURL:        stripe.String("https://pharmonico.com/payment/success?session_id={CHECKOUT_SESSION_ID}"),
        CancelURL:         stripe.String("https://pharmonico.com/payment/cancel"),
        CustomerEmail:     stripe.String(rx.Patient.Email),
        LineItems: []*stripe.CheckoutSessionLineItemParams{
            {
                PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
                    Currency: stripe.String("usd"),
                    ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
                        Name:        stripe.String(fmt.Sprintf("Prescription: %s", rx.Medication.DrugName)),
                        Description: stripe.String(fmt.Sprintf("Copay for %s", rx.PrescriptionNumber)),
                    },
                    UnitAmount: stripe.Int64(int64(finalCopay * 100)), // Convert to cents
                },
                Quantity: stripe.Int64(1),
            },
        },
        Metadata: map[string]string{
            "prescription_id": rx.ID.Hex(),
            "patient_id":      rx.PatientID.Hex(),
        },
    }
    
    session, err := stripeSession.New(params)
    if err != nil {
        return nil, err
    }
    
    // 3. Store payment record
    payment := &models.Payment{
        PrescriptionID:        rx.ID,
        PatientID:             rx.PatientID,
        Amount:                finalCopay,
        Currency:              "USD",
        StripePaymentLink:     session.URL,
        StripeSessionID:       session.ID,
        StripePaymentIntentID: session.PaymentIntent.ID,
        Status:                "pending",
        LinkCreatedAt:         time.Now(),
        LinkExpiresAt:         time.Now().Add(24 * time.Hour),
    }
    
    s.repo.CreatePayment(payment)
    
    // 4. Update prescription status
    s.repo.UpdatePrescriptionStatus(rx.ID, "awaiting_payment")
    
    // 5. Publish event
    s.kafka.Publish("payment.link.created", Event{
        PrescriptionID: rx.ID.Hex(),
        PaymentID:      payment.ID.Hex(),
        Amount:         finalCopay,
        Timestamp:      time.Now(),
    })
    
    return payment, nil
}
```

**Stripe Webhook Handler** (`internal/handlers/webhooks.go`):

```go
func (h *WebhookHandler) HandleStripeWebhook(c *gin.Context) {
    // 1. Verify webhook signature
    payload, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid payload"})
        return
    }
    
    signature := c.GetHeader("Stripe-Signature")
    
    event, err := webhook.ConstructEvent(
        payload,
        signature,
        os.Getenv("STRIPE_WEBHOOK_SECRET"),
    )
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid signature"})
        return
    }
    
    // 2. Handle event type
    switch event.Type {
    case "checkout.session.completed":
        h.handleCheckoutCompleted(event)
    case "payment_intent.succeeded":
        h.handlePaymentSucceeded(event)
    case "payment_intent.payment_failed":
        h.handlePaymentFailed(event)
    }
    
    c.JSON(200, gin.H{"received": true})
}

func (h *WebhookHandler) handleCheckoutCompleted(event stripe.Event) {
    var session stripe.CheckoutSession
    err := json.Unmarshal(event.Data.Raw, &session)
    if err != nil {
        log.Error("Error parsing webhook", err)
        return
    }
    
    // Get prescription ID from metadata
    prescriptionID := session.Metadata["prescription_id"]
    
    // Update payment record
    payment := h.repo.GetPaymentByStripeSessionID(session.ID)
    payment.Status = "paid"
    paidAt := time.Now()
    payment.PaidAt = &paidAt
    payment.PaymentMethod = string(session.PaymentMethodTypes[0])
    payment.ReceiptURL = session.ReceiptURL
    h.repo.UpdatePayment(payment)
    
    // Update prescription status
    h.repo.UpdatePrescriptionStatus(prescriptionID, "paid")
    
    // Cache receipt in Redis
    h.redis.Set(
        context.Background(),
        fmt.Sprintf("payment_receipt:%s", payment.ID.Hex()),
        payment,
        24*time.Hour,
    )
    
    // Log to audit
    h.auditLog(AuditEntry{
        EventType:      "payment_completed",
        PrescriptionID: prescriptionID,
        UserID:         "patient",
        Action:         "payment_successful",
        Details: map[string]interface{}{
            "amount":         payment.Amount,
            "payment_method": payment.PaymentMethod,
            "stripe_session": session.ID,
        },
    })
    
    // Publish event
    h.kafka.Publish("payment.completed", Event{
        PrescriptionID: prescriptionID,
        PaymentID:      payment.ID.Hex(),
        AmountPaid:     payment.Amount,
        Timestamp:      time.Now(),
    })
}
```

---

### **TASK-4.4: Payment Notification System**
**Estimate**: 8 hours

**Notification Service** (`internal/services/notification.go`):

```go
func (s *NotificationService) SendPaymentLink(rx *models.Prescription, payment *models.Payment) error {
    // Prepare cost breakdown for email
    adjudication := rx.Adjudication
    
    // Email notification
    emailData := map[string]interface{}{
        "patient_name":          rx.Patient.FirstName,
        "drug_name":             rx.Medication.DrugName,
        "original_drug_cost":    adjudication.CostBreakdown.TotalDrugCost,
        "insurance_covered":     adjudication.CostBreakdown.InsuranceCovered,
        "initial_copay":         adjudication.CostBreakdown.InitialCopay,
        "manufacturer_discount": adjudication.CostBreakdown.ManufacturerDiscount,
        "final_copay":           adjudication.CostBreakdown.FinalPatientCopay,
        "patient_savings":       adjudication.CostBreakdown.PatientSavings,
        "programs_applied":      adjudication.CostBreakdown.ProgramsApplied,
        "payment_link":          payment.StripePaymentLink,
        "expires_in":            "24 hours",
    }
    
    err := s.sendgrid.Send(Email{
        To:           rx.Patient.Email,
        From:         "noreply@pharmonico.com",
        Subject:      fmt.Sprintf("Complete Your Payment - $%.2f Copay", payment.Amount),
        TemplateName: "payment_request",
        Data:         emailData,
    })
    
    // SMS notification
    smsMessage := fmt.Sprintf(
        "Pharmonico: Your %s Rx is ready. Final copay: $%.2f (saved $%.0f!). Pay now: %s",
        rx.Medication.DrugName,
        payment.Amount,
        adjudication.CostBreakdown.PatientSavings,
        shortenURL(payment.StripePaymentLink),
    )
    
    s.twilio.Send(SMS{
        To:   rx.Patient.Phone,
        From: "+1234567890",
        Body: smsMessage,
    })
    
    // Store notification record
    notification := &models.Notification{
        PrescriptionID: rx.ID,
        PatientID:      rx.PatientID,
        Type:           "payment_request",
        Channel:        "email",
        EmailTo:        rx.Patient.Email,
        EmailSubject:   emailData["subject"].(string),
        Status:         "sent",
        SentAt:         ptr(time.Now()),
    }
    s.repo.CreateNotification(notification)
    
    return err
}
```

**Email Template** (SendGrid):

```html
<!-- templates/payment_request.html -->
<!DOCTYPE html>
<html>
<head>
    <style>
        .container { max-width: 600px; margin: 0 auto; font-family: Arial; }
        .header { background: #4CAF50; color: white; padding: 20px; }
        .cost-breakdown { background: #f5f5f5; padding: 20px; margin: 20px 0; }
        .cost-row { display: flex; justify-content: space-between; padding: 10px 0; }
        .total { font-size: 24px; font-weight: bold; color: #4CAF50; }
        .button { background: #4CAF50; color: white; padding: 15px 30px; text-decoration: none; display: inline-block; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Your Prescription is Ready!</h1>
        </div>
        
        <p>Hi {{patient_name}},</p>
        
        <p>Great news! Your prescription for <strong>{{drug_name}}</strong> has been processed and is ready for fulfillment.</p>
        
        <div class="cost-breakdown">
            <h3>Cost Breakdown:</h3>
            
            <div class="cost-row">
                <span>Original Drug Cost:</span>
                <span>${{original_drug_cost}}</span>
            </div>
            
            <div class="cost-row">
                <span>Insurance Coverage:</span>
                <span>-${{insurance_covered}}</span>
            </div>
            
            <div class="cost-row">
                <span>Initial Copay:</span>
                <span>${{initial_copay}}</span>
            </div>
            
            {{#if manufacturer_discount}}
            <div class="cost-row" style="color: #4CAF50;">
                <span>Manufacturer Discount ({{programs_applied}}):</span>
                <span>-${{manufacturer_discount}}</span>
            </div>
            {{/if}}
            
            <hr>
            
            <div class="cost-row total">
                <span>Your Final Copay:</span>
                <span>${{final_copay}}</span>
            </div>
            
            <div class="cost-row" style="color: #4CAF50; font-weight: bold;">
                <span>You Saved:</span>
                <span>${{patient_savings}}</span>
            </div>
        </div>
        
        <p style="text-align: center; margin: 30px 0;">
            <a href="{{payment_link}}" class="button">Pay ${{final_copay}} Now</a>
        </p>
        
        <p><small>This payment link expires in {{expires_in}}.</small></p>
        
        <p>Questions? Contact us at support@pharmonico.com or (555) 123-4567.</p>
    </div>
</body>
</html>
```

---

### **TASK-4.5: Payment Monitoring Worker**
**Estimate**: 6 hours

**Payment Timeout Handler**:

```go
func (w *PaymentMonitorWorker) Start(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            w.checkExpiredPayments(ctx)
            w.sendReminders(ctx)
        case <-ctx.Done():
            return
        }
    }
}

func (w *PaymentMonitorWorker) checkExpiredPayments(ctx context.Context) {
    // Find payments that expired
    expiredPayments := w.repo.FindPayments(bson.M{
        "status":           "pending",
        "link_expires_at": bson.M{"$lt": time.Now()},
    })
    
    for _, payment := range expiredPayments {
        // Update status
        payment.Status = "expired"
        w.repo.UpdatePayment(payment)
        
        // Update prescription
        w.repo.UpdatePrescriptionStatus(payment.PrescriptionID, "payment_timeout")
        
        // Notify ops team
        w.notificationSvc.SendOpsAlert(Alert{
            Type:    "payment_timeout",
            Message: fmt.Sprintf("Payment expired for Rx %s", payment.PrescriptionID.Hex()),
            Priority: "medium",
        })
    }
}

func (w *PaymentMonitorWorker) sendReminders(ctx context.Context) {
    // Find payments expiring in 2 hours
    twoHoursFromNow := time.Now().Add(2 * time.Hour)
    
    paymentsNeedingReminder := w.repo.FindPayments(bson.M{
        "status":           "pending",
        "link_expires_at": bson.M{
            "$lt": twoHoursFromNow,
            "$gt": time.Now(),
        },
        "reminder_sent": false,
    })
    
    for _, payment := range paymentsNeedingReminder {
        rx := w.repo.GetPrescription(payment.PrescriptionID)
        
        // Send reminder
        w.notificationSvc.SendPaymentReminder(rx, payment)
        
        // Mark reminder sent
        payment.ReminderSent = true
        w.repo.UpdatePayment(payment)
    }
}
```

---

## **Sprint 4 Deliverables:**

✅ Pharmacy adjudication API integration  
✅ Mock pharmacy API for development  
✅ Prior authorization workflow  
✅ PA management in ops dashboard  
✅ Stripe checkout session creation  
✅ Stripe webhook handler  
✅ Payment link generation  
✅ Payment notification system (email/SMS)  
✅ Beautiful email template with cost breakdown  
✅ Payment monitoring worker  
✅ Payment timeout handling  
✅ Kafka event: payment.completed

---

# 🚀 **SPRINT 5 — Shipping & Completion (Week 13-15)**

## **Goal**: Integrate Shippo, implement delivery tracking, complete notification system, finalize audit logging

---

### **TASK-5.1: Shippo Shipping Integration**
**Estimate**: 12 hours

**Shipping Service** (`internal/services/shipping.go`):

```go
type ShippingService struct {
    repo       Repository
    shippoAPI  ShippoClient
    kafka      KafkaProducer
    notifySvc  NotificationService
}

func (s *ShippingService) CreateShipment(prescriptionID string) (*models.Shipment, error) {
    // 1. Fetch prescription with all details
    rx := s.repo.GetPrescriptionWithDetails(prescriptionID)
    pharmacy := s.repo.GetPharmacy(rx.SelectedPharmacyID)
    
    // 2. Prepare shipment request
    shipmentRequest := s.prepareShippoRequest(rx, pharmacy)
    
    // 3. Create shipment via Shippo
    shipment, err := s.shippoAPI.CreateShipment(shipmentRequest)
    if err != nil {
        return nil, err
    }
    
    // 4. Select best rate (usually cheapest with reasonable delivery time)
    rate := s.selectBestRate(shipment.Rates, rx.Medication.IsControlled)
    
    // 5. Purchase shipping label
    transaction, err := s.shippoAPI.PurchaseLabel(rate.ObjectID)
    if err != nil {
        return nil, err
    }
    
    // 6. Store shipment in MongoDB
    shipmentRecord := &models.Shipment{
        PrescriptionID:      rx.ID,
        PharmacyID:          pharmacy.ID,
        ShippoShipmentID:    shipment.ObjectID,
        ShippoRateID:        rate.ObjectID,
        ShippoTransactionID: transaction.ObjectID,
        Carrier:             rate.Provider,
        ServiceLevel:        rate.ServiceLevel.Name,
        TrackingNumber:      transaction.TrackingNumber,
        TrackingURL:         transaction.TrackingURLProvider,
        LabelURL:            transaction.LabelURL,
        ShippingAddress:     rx.Patient.Address,
        ShippingCost:        rate.Amount,
        Status:              "label_created",
        RequiresSignature:   rx.Medication.IsControlled,
        RequiresColdChain:   rx.Medication.RequiresColdChain,
        TrackingHistory:     []models.ShipmentEvent{},
    }
    
    s.repo.CreateShipment(shipmentRecord)
    
    // 7. Update prescription
    s.repo.UpdatePrescription(rx.ID, bson.M{
        "$set": bson.M{
            "shipment_id": shipmentRecord.ID,
            "status":      "shipped",
            "shipped_at":  time.Now(),
        },
    })
    
    // 8. Create tracking job
    s.createTrackingJob(shipmentRecord)
    
    // 9. Publish event
    s.kafka.Publish("shipment.label.created", Event{
        PrescriptionID: rx.ID.Hex(),
        ShipmentID:     shipmentRecord.ID.Hex(),
        TrackingNumber: shipmentRecord.TrackingNumber,
        Carrier:        shipmentRecord.Carrier,
        Timestamp:      time.Now(),
    })
    
    return shipmentRecord, nil
}

func (s *ShippingService) prepareShippoRequest(rx *models.Prescription, pharmacy *models.Pharmacy) ShippoShipmentRequest {
    return ShippoShipmentRequest{
        AddressFrom: ShippoAddress{
            Name:    pharmacy.Name,
            Street1: pharmacy.Address.Line1,
            City:    pharmacy.Address.City,
            State:   pharmacy.Address.State,
            Zip:     pharmacy.Address.ZIP,
            Country: "US",
            Phone:   pharmacy.Phone,
            Email:   pharmacy.Email,
        },
        AddressTo: ShippoAddress{
            Name:    fmt.Sprintf("%s %s", rx.Patient.FirstName, rx.Patient.LastName),
            Street1: rx.Patient.Address.Line1,
            Street2: rx.Patient.Address.Line2,
            City:    rx.Patient.Address.City,
            State:   rx.Patient.Address.State,
            Zip:     rx.Patient.Address.ZIP,
            Country: "US",
            Phone:   rx.Patient.Phone,
            Email:   rx.Patient.Email,
        },
        Parcels: []ShippoParcel{
            {
                Length:       "10",
                Width:        "8",
                Height:       "4",
                DistanceUnit: "in",
                Weight:       "1",
                MassUnit:     "lb",
            },
        },
        Extra: ShippoExtra{
            SignatureConfirmation: s.getSignatureLevel(rx.Medication.IsControlled),
            Insurance: ShippoInsurance{
                Amount:   fmt.Sprintf("%.2f", rx.Adjudication.CostBreakdown.TotalDrugCost),
                Currency: "USD",
            },
        },
    }
}

func (s *ShippingService) selectBestRate(rates []ShippoRate, isControlled bool) ShippoRate {
    // Filter for acceptable service levels
    acceptable := []ShippoRate{}
    
    for _, rate := range rates {
        // For controlled substances, require faster delivery
        if isControlled && rate.EstimatedDays > 2 {
            continue
        }
        
        // Skip rates that are too slow
        if rate.EstimatedDays > 5 {
            continue
        }
        
        acceptable = append(acceptable, rate)
    }
    
    // Sort by cost
    sort.Slice(acceptable, func(i, j int) bool {
        costI, _ := strconv.ParseFloat(acceptable[i].Amount, 64)
        costJ, _ := strconv.ParseFloat(acceptable[j].Amount, 64)
        return costI < costJ
    })
    
    // Return cheapest acceptable option
    return acceptable[0]
}
```

**Shippo API Client** (`internal/clients/shippo.go`):

```go
type ShippoClient struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
}

func NewShippoClient(apiKey string) *ShippoClient {
    return &ShippoClient{
        apiKey:  apiKey,
        baseURL: "https://api.goshippo.com",
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (c *ShippoClient) CreateShipment(req ShippoShipmentRequest) (*ShippoShipment, error) {
    body, _ := json.Marshal(req)
    
    httpReq, _ := http.NewRequest("POST", c.baseURL+"/shipments", bytes.NewBuffer(body))
    httpReq.Header.Set("Authorization", fmt.Sprintf("ShippoToken %s", c.apiKey))
    httpReq.Header.Set("Content-Type", "application/json")
    
    resp, err := c.httpClient.Do(httpReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var shipment ShippoShipment
    json.NewDecoder(resp.Body).Decode(&shipment)
    
    return &shipment, nil
}

func (c *ShippoClient) PurchaseLabel(rateID string) (*ShippoTransaction, error) {
    body, _ := json.Marshal(map[string]string{
        "rate": rateID,
        "async": "false",
    })
    
    httpReq, _ := http.NewRequest("POST", c.baseURL+"/transactions", bytes.NewBuffer(body))
    httpReq.Header.Set("Authorization", fmt.Sprintf("ShippoToken %s", c.apiKey))
    httpReq.Header.Set("Content-Type", "application/json")
    
    resp, err := c.httpClient.Do(httpReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var transaction ShippoTransaction
    json.NewDecoder(resp.Body).Decode(&transaction)
    
    return &transaction, nil
}

func (c *ShippoClient) GetTracking(carrier, trackingNumber string) (*ShippoTracking, error) {
    url := fmt.Sprintf("%s/tracks/%s/%s", c.baseURL, carrier, trackingNumber)
    
    httpReq, _ := http.NewRequest("GET", url, nil)
    httpReq.Header.Set("Authorization", fmt.Sprintf("ShippoToken %s", c.apiKey))
    
    resp, err := c.httpClient.Do(httpReq)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var tracking ShippoTracking
    json.NewDecoder(resp.Body).Decode(&tracking)
    
    return &tracking, nil
}
```

---

### **TASK-5.2: Delivery Tracking Worker**
**Estimate**: 10 hours

**Tracking Worker** (`internal/workers/tracking.go`):

```go
type TrackingWorker struct {
    repo      Repository
    shippoAPI ShippoClient
    kafka     KafkaProducer
    notifySvc NotificationService
}

func (w *TrackingWorker) Start(ctx context.Context) {
    ticker := time.NewTicker(60 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            w.updateAllShipments(ctx)
        case <-ctx.Done():
            return
        }
    }
}

func (w *TrackingWorker) updateAllShipments(ctx context.Context) {
    // Fetch active shipments
    jobs := w.repo.FindTrackingJobs(bson.M{
        "status":       "pending",
        "next_check_at": bson.M{"$lte": time.Now()},
    })
    
    for _, job := range jobs {
        w.processShipment(ctx, job)
    }
}

func (w *TrackingWorker) processShipment(ctx context.Context, job TrackingJob) {
    // Lock the job
    w.repo.LockJob(job.ID, "tracking-worker-1")
    defer w.repo.UnlockJob(job.ID)
    
    // Fetch shipment
    shipment := w.repo.GetShipment(job.ShipmentID)
    
    // Get tracking info from Shippo
    tracking, err := w.shippoAPI.GetTracking(shipment.Carrier, shipment.TrackingNumber)
    if err != nil {
        log.Error("Failed to get tracking", err)
        w.repo.MarkJobFailed(job.ID, err.Error())
        return
    }
    
    // Process tracking events
    newEvents := w.extractNewEvents(shipment.TrackingHistory, tracking.TrackingHistory)
    
    for _, event := range newEvents {
        shipmentEvent := models.ShipmentEvent{
            Status:    event.Status,
            Location:  event.Location.City + ", " + event.Location.State,
            Message:   event.StatusDetails,
            Timestamp: event.OccurredAt,
        }
        
        // Add to history
        shipment.TrackingHistory = append(shipment.TrackingHistory, shipmentEvent)
        
        // Update shipment status
        shipment.Status = w.mapShippoStatus(event.Status)
        
        // Handle status changes
        w.handleStatusChange(shipment, event.Status)
    }
    
    // Update in database
    w.repo.UpdateShipment(shipment)
    
    // Schedule next check
    nextCheck := w.calculateNextCheck(shipment.Status)
    w.repo.UpdateTrackingJob(job.ID, bson.M{
        "next_check_at": nextCheck,
        "last_checked_at": time.Now(),
    })
}

func (w *TrackingWorker) handleStatusChange(shipment *models.Shipment, status string) {
    rx := w.repo.GetPrescription(shipment.PrescriptionID)
    
    switch status {
    case "TRANSIT":
        // Update prescription status
        w.repo.UpdatePrescriptionStatus(rx.ID, "in_transit")
        
        // Send notification
        w.notifySvc.SendShipmentUpdate(rx, shipment, "Your prescription is on its way!")
        
    case "OUT_FOR_DELIVERY":
        // Update status
        w.repo.UpdatePrescriptionStatus(rx.ID, "out_for_delivery")
        
        // Send notification
        w.notifySvc.SendShipmentUpdate(rx, shipment, "Your prescription is out for delivery today!")
        
    case "DELIVERED":
        // Update shipment
        deliveredAt := time.Now()
        shipment.ActualDelivery = &deliveredAt
        
        // Update prescription to completed
        w.repo.UpdatePrescription(rx.ID, bson.M{
            "$set": bson.M{
                "status":       "completed",
                "delivered_at": deliveredAt,
            },
        })
        
        // Final audit log
        w.auditLog(AuditEntry{
            EventType:      "prescription_completed",
            PrescriptionID: rx.ID.Hex(),
            UserID:         "system",
            Action:         "delivery_confirmed",
            Details: map[string]interface{}{
                "tracking_number": shipment.TrackingNumber,
                "delivered_at":    deliveredAt,
                "total_days":      time.Since(rx.CreatedAt).Hours() / 24,
            },
        })
        
        // Publish event
        w.kafka.Publish("shipment.delivered", Event{
            PrescriptionID: rx.ID.Hex(),
            ShipmentID:     shipment.ID.Hex(),
            DeliveredAt:    deliveredAt,
            Timestamp:      time.Now(),
        })
        
        // Send delivery confirmation
        w.notifySvc.SendDeliveryConfirmation(rx, shipment)
        
        // Mark tracking job complete
        w.repo.CompleteTrackingJob(shipment.ID)
        
    case "FAILURE", "RETURNED":
        // Handle delivery exceptions
        w.handleDeliveryException(shipment, status)
    }
}

func (w *TrackingWorker) calculateNextCheck(status string) time.Time {
    switch status {
    case "label_created", "picked_up":
        return time.Now().Add(2 * time.Hour) // Check every 2 hours initially
    case "in_transit":
        return time.Now().Add(4 * time.Hour) // Check every 4 hours
    case "out_for_delivery":
        return time.Now().Add(30 * time.Minute) // Check every 30 min
    default:
        return time.Now().Add(1 * time.Hour)
    }
}

func (w *TrackingWorker) handleDeliveryException(shipment *models.Shipment, status string) {
    // Update status
    shipment.Status = "exception"
    w.repo.UpdateShipment(shipment)
    
    // Alert ops team
    w.notifySvc.SendOpsAlert(Alert{
        Type:     "delivery_exception",
        Message:  fmt.Sprintf("Delivery exception for tracking %s: %s", shipment.TrackingNumber, status),
        Priority: "high",
        Data: map[string]interface{}{
            "prescription_id": shipment.PrescriptionID.Hex(),
            "tracking_number": shipment.TrackingNumber,
            "status":          status,
        },
    })
    
    // Notify patient
    rx := w.repo.GetPrescription(shipment.PrescriptionID)
    w.notifySvc.SendDeliveryException(rx, shipment, status)
}
```

---

### **TASK-5.3: Complete Notification System**
**Estimate**: 8 hours

**Notification Templates**:

```go
// Shipping Notification
func (s *NotificationService) SendShippingNotification(rx *models.Prescription, shipment *models.Shipment) error {
    emailData := map[string]interface{}{
        "patient_name":     rx.Patient.FirstName,
        "drug_name":        rx.Medication.DrugName,
        "tracking_number":  shipment.TrackingNumber,
        "tracking_url":     shipment.TrackingURL,
        "carrier":          shipment.Carrier,
        "estimated_delivery": shipment.EstimatedDelivery.Format("January 2, 2006"),
    }
    
    s.sendgrid.Send(Email{
        To:           rx.Patient.Email,
        From:         "noreply@pharmonico.com",
        Subject:      "Your Prescription Has Shipped! 📦",
        TemplateName: "shipping_notification",
        Data:         emailData,
    })
    
    // SMS
    smsMessage := fmt.Sprintf(
        "Pharmonico: Your %s has shipped via %s. Track: %s",
        rx.Medication.DrugName,
        shipment.Carrier,
        shortenURL(shipment.TrackingURL),
    )
    
    s.twilio.Send(SMS{
        To:   rx.Patient.Phone,
        From: "+1234567890",
        Body: smsMessage,
    })
    
    return nil
}

// Delivery Confirmation
func (s *NotificationService) SendDeliveryConfirmation(rx *models.Prescription, shipment *models.Shipment) error {
    emailData := map[string]interface{}{
        "patient_name":     rx.Patient.FirstName,
        "drug_name":        rx.Medication.DrugName,
        "delivered_at":     shipment.ActualDelivery.Format("January 2, 2006 at 3:04 PM"),
        "tracking_number":  shipment.TrackingNumber,
        "prescription_id":  rx.PrescriptionNumber,
    }
    
    s.sendgrid.Send(Email{
        To:           rx.Patient.Email,
        From:         "noreply@pharmonico.com",
        Subject:      "Your Prescription Has Been Delivered! ✅",
        TemplateName: "delivery_confirmation",
        Data:         emailData,
    })
    
    // SMS
    smsMessage := fmt.Sprintf(
        "Pharmonico: Your %s prescription has been delivered! Thank you for using our service.",
        rx.Medication.DrugName,
    )
    
    s.twilio.Send(SMS{
        To:   rx.Patient.Phone,
        From: "+1234567890",
        Body: smsMessage,
    })
    
    return nil
}
```

---

### **TASK-5.4: Comprehensive Audit Logging**
**Estimate**: 6 hours

**Audit Log Service** (`internal/services/audit.go`):

```go
type AuditLogService struct {
    pgDB *sql.DB
}

func (s *AuditLogService) Log(entry AuditEntry) error {
    query := `
        INSERT INTO audit_logs (
            event_type,
            prescription_id,
            user_id,
            action,
            details,
            ip_address,
            user_agent,
            created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
    `
    
    detailsJSON, _ := json.Marshal(entry.Details)
    
    _, err := s.pgDB.Exec(
        query,
        entry.EventType,
        entry.PrescriptionID,
        entry.UserID,
        entry.Action,
        detailsJSON,
        entry.IPAddress,
        entry.UserAgent,
    )
    
    return err
}

func (s *AuditLogService) GetPrescriptionAuditTrail(prescriptionID string) ([]AuditEntry, error) {
    query := `
        SELECT 
            event_type,
            user_id,
            action,
            details,
            ip_address,
            created_at
        FROM audit_logs
        WHERE prescription_id = $1
        ORDER BY created_at ASC
    `
    
    rows, err := s.pgDB.Query(query, prescriptionID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    entries := []AuditEntry{}
    for rows.Next() {
        var entry AuditEntry
        var detailsJSON []byte
        
        rows.Scan(
            &entry.EventType,
            &entry.UserID,
            &entry.Action,
            &detailsJSON,
            &entry.IPAddress,
            &entry.CreatedAt,
        )
        
        json.Unmarshal(detailsJSON, &entry.Details)
        entries = append(entries, entry)
    }
    
    return entries, nil
}
```

**Audit Log Viewer in Dashboard**:

```typescript
// src/pages/dashboard/AuditLogTab.tsx
export const AuditLogTab = () => {
  const [prescriptionId, setPrescriptionId] = useState('');
  const { auditLogs, loading } = useAuditLogs(prescriptionId);
  
  return (
    <Container>
      <TextField
        label="Prescription ID"
        value={prescriptionId}
        onChange={(e) => setPrescriptionId(e.target.value)}
        placeholder="Enter prescription ID to view audit trail"
      />
      
      {loading && <CircularProgress />}
      
      {auditLogs && (
        <Timeline>
          {auditLogs.map((log, index) => (
            <TimelineItem key={index}>
              <TimelineSeparator>
                <TimelineDot color={getEventColor(log.event_type)} />
                {index < auditLogs.length - 1 && <TimelineConnector />}
              </TimelineSeparator>
              
              <TimelineContent>
                <Paper elevation={2} sx={{ p: 2 }}>
                  <Typography variant="h6">
                    {formatEventType(log.event_type)}
                  </Typography>
                  
                  <Typography variant="body2" color="textSecondary">
                    {log.action}
                  </Typography>
                  
                  <Typography variant="caption" display="block">
                    By: {log.user_id} | {formatDate(log.created_at)}
                  </Typography>
                  
                  {log.details && (
                    <Accordion sx={{ mt: 1 }}>
                      <AccordionSummary>
                        <Typography variant="caption">View Details</Typography>
                      </AccordionSummary>
                      <AccordionDetails>
                        <pre>{JSON.stringify(log.details, null, 2)}</pre>
                      </AccordionDetails>
                    </Accordion>
                  )}
                </Paper>
              </TimelineContent>
            </TimelineItem>
          ))}
        </Timeline>
      )}
    </Container>
  );
};
```

---

### **TASK-5.5: Ops Dashboard Finalization**
**Estimate**: 12 hours

**Complete Tab System**:

```typescript
// src/pages/dashboard/Dashboard.tsx
export const Dashboard = () => {
  const [activeTab, setActiveTab] = useState('intake');
  
  const tabs = [
    { id: 'intake', label: 'Intake', count: useTabCount('received') },
    { id: 'validation', label: 'Validation', count: useTabCount('validated') },
    { id: 'enrollment', label: 'Enrollment', count: useTabCount('awaiting_enrollment') },
    { id: 'routing', label: 'Routing', count: useTabCount('awaiting_pharmacy_selection') },
    { id: 'prior-auth', label: 'Prior Auth', count: usePACount('pending') },
    { id: 'adjudication', label: 'Adjudication', count: useTabCount('adjudicated') },
    { id: 'payment', label: 'Payment', count: useTabCount('awaiting_payment') },
    { id: 'fulfillment', label: 'Fulfillment', count: useTabCount(['shipped', 'in_transit']) },
    { id: 'completed', label: 'Completed', count: useTabCount('completed') },
    { id: 'issues', label: 'Issues', count: useIssuesCount() },
    { id: 'audit', label: 'Audit Logs', count: null },
  ];
  
  return (
    <DashboardLayout>
      <Header>
        <Typography variant="h4">Pharmonico Operations Dashboard</Typography>
        <UserMenu />
      </Header>
      
      <Tabs value={activeTab} onChange={(e, v) => setActiveTab(v)}>
        {tabs.map(tab => (
          <Tab
            key={tab.id}
            value={tab.id}
            label={
              <Box display="flex" alignItems="center" gap={1}>
                {tab.label}
                {tab.count > 0 && (
                  <Chip size="small" label={tab.count} color="primary" />
                )}
              </Box>
            }
          />
        ))}
      </Tabs>
      
      <TabPanel value={activeTab} index="intake">
        <IntakeTab />
      </TabPanel>
      
      <TabPanel value={activeTab} index="validation">
        <ValidationTab />
      </TabPanel>
      
      <TabPanel value={activeTab} index="enrollment">
        <EnrollmentTab />
      </TabPanel>
      
      <TabPanel value={activeTab} index="routing">
        <PharmacyRoutingTab />
      </TabPanel>
      
      <TabPanel value={activeTab} index="prior-auth">
        <PriorAuthTab />
      </TabPanel>
      
      <TabPanel value={activeTab} index="adjudication">
        <AdjudicationTab />
      </TabPanel>
      
      <TabPanel value={activeTab} index="payment">
        <PaymentTab />
      </TabPanel>
      
      <TabPanel value={activeTab} index="fulfillment">
        <FulfillmentTab />
      </TabPanel>
      
      <TabPanel value={activeTab} index="completed">
        <CompletedTab />
      </TabPanel>
      
      <TabPanel value={activeTab} index="issues">
        <IssuesTab />
      </TabPanel>
      
      <TabPanel value={activeTab} index="audit">
        <AuditLogTab />
      </TabPanel>
    </DashboardLayout>
  );
};
```

**Global Search**:

```typescript
// src/components/GlobalSearch.tsx
export const GlobalSearch = () => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  
  const handleSearch = async () => {
    const response = await api.get('/api/v1/search', {
      params: { q: query }
    });
    setResults(response.data);
  };
  
  return (
    <Autocomplete
      freeSolo
      options={results}
      getOptionLabel={(option) => 
        `${option.prescription_number} - ${option.patient.name}`
      }
      renderInput={(params) => (
        <TextField
          {...params}
          placeholder="Search by Rx ID, Patient Name, or Tracking #"
          onChange={(e) => setQuery(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
        />
      )}
      renderOption={(props, option) => (
        <li {...props}>
          <Box>
            <Typography variant="body1">
              {option.prescription_number}
            </Typography>
            <Typography variant="caption" color="textSecondary">
              {option.patient.name} - {option.medication.drug_name}
            </Typography>
            <Chip
              size="small"
              label={option.status}
              color={getStatusColor(option.status)}
            />
          </Box>
        </li>
      )}
      onChange={(e, value) => {
        if (value?.id) {
          navigate(`/prescription/${value.id}`);
        }
      }}
    />
  );
};
```

---

## **Sprint 5 Deliverables:**

✅ Shippo integration for label generation  
✅ Delivery tracking worker with polling  
✅ Real-time shipment status updates  
✅ Complete notification system (all templates)  
✅ Comprehensive audit logging  
✅ Audit log viewer in dashboard  
✅ All dashboard tabs completed  
✅ Global search functionality  
✅ Delivery confirmation flow  
✅ Exception handling for failed deliveries

---

# 🚀 **SPRINT 6 — Testing, Polish & Documentation (Week 16)**

## **Goal**: End-to-end testing, bug fixes, performance optimization, comprehensive documentation

---

### **TASK-6.1: End-to-End Integration Tests**
**Estimate**: 16 hours

**E2E Test Suite** (`backend/tests/e2e/prescription_flow_test.go`):

```go
func TestCompletePresc riptionFlow(t *testing.T) {
    // Setup
    cleanDB(t)
    seedTestData(t)
    
    // Step 1: Intake
    t.Run("Intake Prescription", func(t *testing.T) {
        ncpdpPayload := loadTestNCPDP("humira_prescription.xml")
        
        resp := makeRequest(t, "POST", "/api/v1/prescriptions/intake", ncpdpPayload)
        assert.Equal(t, 200, resp.StatusCode)
        
        var result map[string]interface{}
        json.NewDecoder(resp.Body).Decode(&result)
        
        prescriptionID = result["prescription_id"].(string)
        assert.NotEmpty(t, prescriptionID)
        
        // Verify Kafka event published
        event := waitForKafkaEvent(t, "prescription.intake.received", 5*time.Second)
        assert.NotNil(t, event)
    })
    
    // Step 2: Validation
    t.Run("Validation Process", func(t *testing.T) {
        // Wait for validation worker to process
        time.Sleep(15 * time.Second)
        
        // Check prescription status
        rx := getPrescription(t, prescriptionID)
        assert.Equal(t, "validated", rx.Status)
        assert.True(t, rx.ValidationChecks.NPIValid)
        assert.True(t, rx.ValidationChecks.NDCValid)
        
        // Verify Kafka event
        event := waitForKafkaEvent(t, "prescription.validation.completed", 5*time.Second)
        assert.NotNil(t, event)
    })
    
    // Step 3: Enrollment
    t.Run("Enrollment Flow", func(t *testing.T) {
        // Initiate enrollment
        resp := makeRequest(t, "POST", "/api/v1/enrollment/initiate", map[string]string{
            "prescription_id": prescriptionID,
        })
        assert.Equal(t, 200, resp.StatusCode)
        
        var enrollmentResp map[string]interface{}
        json.NewDecoder(resp.Body).Decode(&enrollmentResp)
        
        token := extractTokenFromMagicLink(enrollmentResp["magic_link"].(string))
        
        // Validate token
        resp = makeRequest(t, "GET", "/api/v1/enrollment/validate/"+token, nil)
        assert.Equal(t, 200, resp.StatusCode)
        
        // Submit enrollment
        enrollmentData := map[string]interface{}{
            "token": token,
            "insurance": map[string]string{
                "payer_name":  "Blue Cross Blue Shield",
                "member_id":   "TEST123",
                "bin":         "610014",
            },
            "hipaa_consent": map[string]interface{}{
                "authorization_text": "I consent...",
                "signature":          "data:image/png;base64,...",
                "signature_name":     "John Doe",
            },
        }
        
        resp = makeRequest(t, "POST", "/api/v1/enrollment/submit", enrollmentData)
        assert.Equal(t, 200, resp.StatusCode)
        
        // Verify status change
        time.Sleep(2 * time.Second)
        rx := getPrescription(t, prescriptionID)
        assert.Equal(t, "enrolled", rx.Status)
    })
    
    // Step 4: Pharmacy Routing
    t.Run("Pharmacy Routing and Selection", func(t *testing.T) {
        // Wait for routing worker
        time.Sleep(25 * time.Second)
        
        rx := getPrescription(t, prescriptionID)
        assert.Equal(t, "awaiting_pharmacy_selection", rx.Status)
        assert.NotEmpty(t, rx.PharmacyRecommendations)
        assert.GreaterOrEqual(t, len(rx.PharmacyRecommendations), 1)
        
        // Select top pharmacy
        pharmacyID := rx.PharmacyRecommendations[0].PharmacyID
        
        resp := makeRequest(t, "POST", 
            fmt.Sprintf("/api/v1/prescriptions/%s/select-pharmacy", prescriptionID),
            map[string]string{"pharmacy_id": pharmacyID},
        )
        assert.Equal(t, 200, resp.StatusCode)
        
        // Verify selection
        rx = getPrescription(t, prescriptionID)
        assert.Equal(t, "pharmacy_selected", rx.Status)
        assert.Equal(t, pharmacyID, rx.SelectedPharmacyID)
    })
    
    // Step 5: Adjudication
    t.Run("Insurance Adjudication", func(t *testing.T) {
        // Wait for adjudication worker
        time.Sleep(20 * time.Second)
        
        rx := getPrescription(t, prescriptionID)
        assert.Equal(t, "adjudicated", rx.Status)
        assert.NotNil(t, rx.AdjudicationID)
        
        // Fetch adjudication details
        adjudication := getAdjudication(t, rx.AdjudicationID)
        assert.Equal(t, "approved", adjudication.PrimaryInsurance.Status)
        assert.Greater(t, adjudication.CostBreakdown.FinalPatientCopay, 0.0)
    })
    
    // Step 6: Payment
    t.Run("Payment Process", func(t *testing.T) {
        // Get payment link
        rx := getPrescription(t, prescriptionID)
        payment := getPayment(t, rx.PaymentID)
        
        assert.NotEmpty(t, payment.StripePaymentLink)
        assert.Equal(t, "pending", payment.Status)
        
        // Simulate successful payment via Stripe webhook
        webhookPayload := createStripeWebhookPayload(payment.StripeSessionID)
        
        resp := makeRequest(t, "POST", "/api/v1/webhooks/stripe", webhookPayload)
        assert.Equal(t, 200, resp.StatusCode)
        
        // Verify payment marked as paid
        time.Sleep(2 * time.Second)
        payment = getPayment(t, payment.ID)
        assert.Equal(t, "paid", payment.Status)
        assert.NotNil(t, payment.PaidAt)
        
        // Verify prescription status
        rx = getPrescription(t, prescriptionID)
        assert.Equal(t, "paid", rx.Status)
    })
    
    // Step 7: Shipping
    t.Run("Shipping and Fulfillment", func(t *testing.T) {
        // Wait for shipping worker
        time.Sleep(35 * time.Second)
        
        rx := getPrescription(t, prescriptionID)
        assert.Equal(t, "shipped", rx.Status)
        assert.NotNil(t, rx.ShipmentID)
        
        // Fetch shipment details
        shipment := getShipment(t, rx.ShipmentID)
        assert.NotEmpty(t, shipment.TrackingNumber)
        assert.NotEmpty(t, shipment.LabelURL)
        assert.Equal(t, "label_created", shipment.Status)
    })
    
    // Step 8: Delivery
    t.Run("Delivery Tracking and Completion", func(t *testing.T) {
        // Simulate delivery status update
        rx := getPrescription(t, prescriptionID)
        shipment := getShipment(t, rx.ShipmentID)
        
        // Mock Shippo webhook for delivery
        mockShippoDelivery(t, shipment.TrackingNumber)
        
        // Wait for tracking worker to process
        time.Sleep(65 * time.Second)
        
        // Verify final status
        rx = getPrescription(t, prescriptionID)
        assert.Equal(t, "completed", rx.Status)
        assert.NotNil(t, rx.DeliveredAt)
        
        // Verify shipment status
        shipment = getShipment(t, rx.ShipmentID)
        assert.Equal(t, "delivered", shipment.Status)
        assert.NotNil(t, shipment.ActualDelivery)
    })
    
    // Step 9: Audit Trail
    t.Run("Complete Audit Trail", func(t *testing.T) {
        auditLogs := getAuditLogs(t, prescriptionID)
        
        // Verify key events are logged
        events := extractEventTypes(auditLogs)
        expectedEvents := []string{
            "intake",
            "validation",
            "enrollment",
            "pharmacy_selection",
            "adjudication",
            "payment",
            "shipping",
            "delivery",
        }
        
        for _, expected := range expectedEvents {
            assert.Contains(t, events, expected)
        }
        
        // Verify HIPAA compliance fields
        for _, log := range auditLogs {
            assert.NotEmpty(t, log.Timestamp)
            assert.NotEmpty(t, log.UserID)
            assert.NotEmpty(t, log.Action)
        }
    })
}
```

---

### **TASK-6.2: Unit Tests for Critical Components**
**Estimate**: 8 hours

```go
// Validation Rules Tests
func TestNPIValidation(t *testing.T) {
    validator := NewNPIValidator()
    
    tests := []struct {
        npi      string
        expected bool
    }{
        {"1234567890", true},
        {"123456789", false},   // Too short
        {"12345678901", false}, // Too long
        {"abcdefghij", false},  // Non-numeric
    }
    
    for _, tt := range tests {
        result := validator.Validate(&models.Prescription{
            Prescriber: models.Prescriber{NPI: tt.npi},
        })
        
        hasError := len(result) > 0
        assert.Equal(t, !tt.expected, hasError)
    }
}

// Pharmacy Scoring Tests
func TestPharmacyScoring(t *testing.T) {
    service := NewRoutingService()
    
    rx := &models.Prescription{
        Patient: models.Patient{
            Address: models.Address{
                Coordinates: models.GeoLocation{Lat: 42.3601, Lng: -71.0589},
            },
        },
        Insurance: models.InsuranceProfile{
            PayerName: "Blue Cross Blue Shield",
        },
    }
    
    pharmacies := []models.Pharmacy{
        createTestPharmacy("nearby", 5.0, "preferred"),
        createTestPharmacy("far", 50.0, "preferred"),
        createTestPharmacy("nearby_standard", 5.0, "standard"),
    }
    
    scored := service.scorePharmacies(pharmacies, rx)
    
    // Nearby + preferred should score highest
    assert.Equal(t, "nearby", scored[0].PharmacyID)
    assert.Greater(t, scored[0].Score.TotalScore, scored[1].Score.TotalScore)
}

// Redis Cache Tests
func TestRedisCaching(t *testing.T) {
    redis := setupTestRedis(t)
    defer redis.FlushDB(context.Background())
    
    // Test magic link storage
    token := uuid.New().String()
    data := map[string]interface{}{
        "prescription_id": "rx_123",
        "expires_at":      time.Now().Add(48 * time.Hour),
    }
    
    redis.Set(context.Background(), "magic_link:"+token, data, 48*time.Hour)
    
    // Retrieve
    result := redis.Get(context.Background(), "magic_link:"+token).Val()
    assert.NotEmpty(t, result)
    
    // Verify TTL
    ttl := redis.TTL(context.Background(), "magic_link:"+token).Val()
    assert.Greater(t, ttl.Hours(), 47.0)
}
```

---

### **TASK-6.3: Performance Testing & Optimization**
**Estimate**: 8 hours

**Load Test Script** (`scripts/load_test.js`):

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 100 },  // Ramp up to 100 users
    { duration: '5m', target: 100 },  // Stay at 100 users
    { duration: '2m', target: 200 },  // Ramp to 200 users
    { duration: '5m', target: 200 },  // Stay at 200
    { duration: '2m', target: 0 },    // Ramp down
  ],
};

export default function() {
  // Test intake endpoint
  let intakeResp = http.post(
    'http://localhost:8080/api/v1/prescriptions/intake',
    JSON.stringify(generateMockPrescription()),
    { headers: { 'Content-Type': 'application/json' } }
  );
  
  check(intakeResp, {
    'intake status is 200': (r) => r.status === 200,
    'intake response time < 500ms': (r) => r.timings.duration < 500,
  });
  
  sleep(1);
  
  // Test search endpoint
  let searchResp = http.get('http://localhost:8080/api/v1/prescriptions?status=received');
  
  check(searchResp, {
    'search status is 200': (r) => r.status === 200,
    'search response time < 200ms': (r) => r.timings.duration < 200,
  });
  
  sleep(1);
}
```

**Performance Optimizations**:

1. **Database Indexing**:
```sql
-- Add indexes for frequently queried fields
CREATE INDEX idx_prescriptions_status ON prescriptions(status);
CREATE INDEX idx_prescriptions_patient ON prescriptions(patient_id);
CREATE INDEX idx_prescriptions_created ON prescriptions(created_at);
CREATE INDEX idx_jobs_status_locked ON validation_jobs(status, locked_at);
```

2. **Redis Connection Pooling**:
```go
func NewRedisClient() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:         "redis:6379",
        PoolSize:     20,
        MinIdleConns: 5,
        MaxRetries:   3,
    })
}
```

3. **MongoDB Query Optimization**:
```go
// Use projection to fetch only needed fields
func (r *Repository) GetPrescriptionBasic(id string) (*models.Prescription, error) {
    projection := bson.M{
        "status": 1,
        "patient.first_name": 1,
        "patient.last_name": 1,
        "medication.drug_name": 1,
    }
    
    var rx models.Prescription
    err := r.db.Collection("prescriptions").
        FindOne(context.Background(), bson.M{"_id": id}, options.FindOne().SetProjection(projection)).
        Decode(&rx)
    
    return &rx, err
}
```

---

### **TASK-6.4: Comprehensive Documentation**
**Estimate**: 12 hours

**README.md**:
````markdown
# Pharmonico - Prescription Fulfillment System

A complete prescription fulfillment platform that handles the entire workflow from intake to delivery.

## Features

- ✅ NCPDP SCRIPT intake
- ✅ Automated prescription validation
- ✅ Patient enrollment with magic links
- ✅ Intelligent pharmacy routing
- ✅ Insurance adjudication coordination
- ✅ Manufacturer program integration
- ✅ Stripe payment processing
- ✅ Shippo shipping integration
- ✅ Real-time delivery tracking
- ✅ HIPAA-compliant audit logging

## Tech Stack

- **Backend**: Go 1.21
- **Frontend**: React 18 + TypeScript
- **Databases**: MongoDB, PostgreSQL, Redis
- **Message Queue**: Kafka
- **Storage**: MinIO (S3-compatible)
- **Payments**: Stripe
- **Shipping**: Shippo
- **Email**: SendGrid
- **SMS**: Twilio

## Quick Start

```bash
# Clone repository
git clone https://github.com/your-org/pharmonico.git
cd pharmonico

# Start development environment
make dev

# Seed databases
make seed

# Run tests
make test
```

## Services

- **API**: http://localhost:8080
- **Frontend**: http://localhost:3000
- **Maildev**: http://localhost:1080
- **MinIO Console**: http://localhost:9001

## Documentation

- [Architecture](docs/architecture/)
- [API Reference](docs/api/)
- [Workflows](docs/workflows/)
- [Deployment](docs/deployment/)

## License

MIT
````

**API Documentation** (`docs/api/README.md`):
````markdown
# API Documentation

Base URL: `http://localhost:8080/api/v1`

## Authentication

All API endpoints require JWT authentication except:
- Enrollment endpoints (magic link based)
- Webhooks (signature verified)

```
Authorization: Bearer <token>
```

## Endpoints

### Prescriptions

#### Create Prescription
```
POST /prescriptions/intake
Content-Type: application/xml or application/json

Response:
{
  "prescription_id": "rx_abc123",
  "status": "received"
}
```

#### Get Prescription
```
GET /prescriptions/:id

Response:
{
  "id": "rx_abc123",
  "status": "adjudicated",
  "patient": {...},
  "medication": {...}
}
```

### Enrollment

#### Initiate Enrollment
```
POST /enrollment/initiate
{
  "prescription_id": "rx_abc123"
}

Response:
{
  "enrollment_id": "enr_xyz",
  "magic_link": "https://enroll.pharmonico.com/enroll/token"
}
```

[See full API documentation...](./full-api-spec.md)
````

---

### **TASK-6.5: Deployment Documentation**
**Estimate**: 6 hours

**Docker Deployment**:

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  api:
    image: pharmonico/api:latest
    restart: always
    environment:
      - ENVIRONMENT=production
      - MONGODB_URI=${MONGODB_URI}
      - POSTGRES_URI=${POSTGRES_URI}
      - REDIS_URI=${REDIS_URI}
      - KAFKA_BROKERS=${KAFKA_BROKERS}
      - STRIPE_SECRET_KEY=${STRIPE_SECRET_KEY}
      - SHIPPO_API_KEY=${SHIPPO_API_KEY}
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '2'
          memory: 2G
```

**Environment Variables Template**:

```bash
# .env.example
# Database
MONGODB_URI=mongodb://localhost:27017/pharmonico
POSTGRES_URI=postgres://user:pass@localhost:5432/pharmonico
REDIS_URI=redis://localhost:6379

# Kafka
KAFKA_BROKERS=localhost:9092

# Storage
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin

# External APIs
STRIPE_SECRET_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...
SHIPPO_API_KEY=shippo_test_...
SENDGRID_API_KEY=SG...
TWILIO_ACCOUNT_SID=AC...
TWILIO_AUTH_TOKEN=...
TWILIO_PHONE_NUMBER=+1...

# App
JWT_SECRET=your-secret-key-here
ENVIRONMENT=development
```

---

## **Sprint 6 Deliverables:**

✅ Complete E2E test suite  
✅ Unit tests for all critical components  
✅ Performance testing with K6  
✅ Database query optimization  
✅ Redis connection pooling  
✅ Comprehensive README  
✅ Complete API documentation  
✅ Deployment guide  
✅ Environment configuration templates  
✅ Bug fixes and polish

---

# 📊 **Project Summary**

## **Total Timeline: 16 Weeks**

| Phase | Duration | Deliverables |
|-------|----------|-------------|
| **Sprint 0** | 2 weeks | Infrastructure, Docker, CI/CD |
| **Sprint 1** | 2 weeks | Intake, Validation, Basic Dashboard |
| **Sprint 2** | 2 weeks | Enrollment, Magic Links, HIPAA |
| **Sprint 3** | 3 weeks | Routing, Program Management |
| **Sprint 4** | 3 weeks | Adjudication, Payment, Prior Auth |
| **Sprint 5** | 3 weeks | Shipping, Tracking, Notifications |
| **Sprint 6** | 1 week | Testing, Polish, Documentation |

## **Key Technologies Mastered**

### Backend
✅ Go microservices architecture  
✅ PostgreSQL job queue with row-level locking  
✅ MongoDB document modeling  
✅ Redis caching strategies  
✅ Kafka event-driven architecture  
✅ RESTful API design  
✅ JWT authentication  
✅ HIPAA compliance

### Frontend
✅ React with TypeScript  
✅ Material-UI component library  
✅ Real-time updates via WebSocket  
✅ Form validation  
✅ File upload handling  
✅ Electronic signature capture

### Integrations
✅ Stripe payment processing  
✅ Shippo shipping labels  
✅ SendGrid email templates  
✅ Twilio SMS  
✅ MinIO object storage

### DevOps
✅ Docker containerization  
✅ Docker Compose orchestration  
✅ GitHub Actions CI/CD  
✅ Load testing with K6

## **Production-Ready Features**

- ✅ Complete prescription fulfillment workflow
- ✅ Automated validation with retry logic
- ✅ Real-time pharmacy capacity tracking
- ✅ Intelligent routing algorithm
- ✅ Two-step insurance adjudication
- ✅ Manufacturer program integration
- ✅ Prior authorization handling
- ✅ Secure payment processing
- ✅ Automated shipping and tracking
- ✅ Multi-channel notifications
- ✅ Comprehensive audit logging
- ✅ Operations dashboard
- ✅ Patient enrollment portal
- ✅ HIPAA compliance

## **Next Steps (Post-MVP)**

### Phase 2 Enhancements
- Multiple refills management
- Prescription transfer handling
- Patient portal with order history
- Pharmacy performance analytics
- Cost savings reports
- Insurance denial appeals
- Automated prior authorization submission
- Prescription synchronization
- Medication adherence tracking

### Technical Improvements
- GraphQL API
- Mobile app (React Native)
- Advanced caching with Redis Cluster
- Kafka Streams for complex event processing
- Elasticsearch for full-text search
- Prometheus + Grafana monitoring
- Distributed tracing with Jaeger
- Blue-green deployments
- Auto-scaling with Kubernetes

---

# 🎓 **Learning Objectives Achieved**

By completing this project, you will have:

✅ Built a **production-ready microservices architecture**  
✅ Implemented **event-driven design** with Kafka  
✅ Mastered **Go** for backend services  
✅ Created a complete **React** dashboard  
✅ Integrated multiple **third-party APIs**  
✅ Implemented **HIPAA-compliant** systems  
✅ Built **background workers** with job queues  
✅ Designed **complex database schemas**  
✅ Implemented **caching strategies** with Redis  
✅ Created **comprehensive test suites**  
✅ Deployed with **Docker** containers  
✅ Understood **US healthcare workflows**  

---

**🎉 Congratulations on completing Pharmonico!**

This project demonstrates enterprise-level software engineering skills and deep understanding of:
- Healthcare system workflows
- Distributed systems architecture
- Event-driven design patterns
- Third-party API integration
- HIPAA compliance requirements
- Production-ready deployment practices

You now have a portfolio project that showcases:
- **Technical Breadth**: Full-stack + DevOps + integrations
- **Domain Knowledge**: Healthcare/pharmacy industry
- **Best Practices**: Testing, documentation, security
- **Scale**: Microservices, caching, job queues
- **Real-World Value**: Solves actual business problems