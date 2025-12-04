# ✅ **SPRINT 0 — Repo, Infrastructure & CI Bootstrap (UPDATED)**

### *(Fully aligned with Kafka, Redis, MinIO, Workers, Multi-Service Architecture)*

Sprint 0 lays the full technical foundation of Pharmonico.
It prepares the monorepo, Dockerized dev environment, infrastructure services, seeds, and CI.

---

# 📦 **TASK-0.1 — Initialize Monorepo (Updated)**

### 🎯 Purpose

Define a clean, scalable monorepo that supports:

* Multiple backend services (API, Workers)
* React frontend
* Event-driven architecture (Kafka)
* Infrastructure-as-code layout

### 📁 Final Directory Structure

(Updated based on Kafka + Redis + observability)

```
pharmonico/
├── backend-go/
│   ├── cmd/
│   │   ├── api/                 # API server entry point
│   │   ├── worker/              # Core worker consumer
│   │   └── scheduler/           # (Optional) cron, polling jobs
│   ├── internal/
│   │   ├── config/              # env, config loading
│   │   ├── database/            # Mongo, Postgres, Redis
│   │   ├── kafka/               # Producer/consumer helpers
│   │   ├── handlers/            # HTTP routing
│   │   ├── middleware/          # cors, auth, logging
│   │   ├── models/              # Mongo + Postgres models
│   │   ├── services/            # Business logic
│   │   └── workers/             # Kafka consumers
│   ├── pkg/
│   │   └── ncpdp/               # NCPDP parser
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── frontend-react/
│   ├── src/
│   ├── public/
│   ├── package.json
│   ├── Dockerfile
│   ├── vite.config.ts
│   └── tailwind.config.js
│
├── infra/
│   ├── docker/
│   │   ├── zookeeper/
│   │   ├── kafka/
│   │   ├── redis/
│   │   └── postgres/
│   └── nginx/
│       └── nginx.conf
│
├── scripts/
│   └── seeds/
│       ├── mongo/
│       └── postgres/
│
├── docs/
│   ├── sprint_plan.md
│   ├── requirements.md
│   └── user_stories.md
│
├── .github/workflows/ci.yml
├── docker-compose.yml
├── Makefile
└── README.md
```

### 🆕 **Notable Changes**

* Added **kafka/** module in backend
* Added **redis** as shared caching + token storage
* Added **scheduler** stub for cron/polling jobs
* Extended **infra/docker/** for per-service overrides

---

# 🚢 **TASK-0.2 — Docker Compose Dev Environment (Updated)**

### 🎯 Purpose

Run the full Pharmonico stack with:

* Event-driven flows (Kafka)
* Caching + token storage (Redis)
* File storage (MinIO)
* Two databases (Mongo + Postgres)
* Development-friendly mounts
* Health checks on all services

---

### 🏗️ Services Included (Updated Table)

| Service       | Purpose                          | Ports      |
| ------------- | -------------------------------- | ---------- |
| **api**       | Go API Server                    | 8080       |
| **worker**    | Kafka consumer worker            | —          |
| **scheduler** | Polling jobs                     | —          |
| **frontend**  | React Vite dev server            | 5173       |
| **mongodb**   | Primary business DB              | 27017      |
| **postgres**  | Jobs + audit DB                  | 5432       |
| **redis**     | Caching, magic links, rate-limit | 6379       |
| **zookeeper** | Kafka dependency                 | 2181       |
| **kafka**     | Event streaming                  | 9092       |
| **minio**     | S3-compatible storage            | 9000, 9001 |
| **maildev**   | Email testing                    | 1080, 1025 |
| **nginx**     | (Optional) reverse proxy         | 80         |

---

### 🔧 Additional Improvements Added

* Kafka configured with internal network & advertised listeners
* Redis persistence volume added
* MinIO automated bucket creation using `minio-init`
* Mongo + Postgres health check before API starts
* Worker depends_on Kafka (with retries)
* Shared network `pharmonico-network`

---

# 📑 **TASK-0.3 — README + Makefile (Updated)**

### 🎯 Purpose

Create developer-friendly onboarding with:

* Streamlined Makefile
* Fully updated README (includes Kafka + Redis)
* Architecture diagrams (updated)

### 🔨 Makefile (New Commands Added)

| Command                   | Description                |
| ------------------------- | -------------------------- |
| `make kafka-topics`       | List Kafka topics          |
| `make kafka-create-topic` | Create topic from template |
| `make redis-cli`          | Redis shell                |
| `make logs-api`           | Only API logs              |
| `make logs-worker`        | Only Worker logs           |
| `make ps`                 | Show all dev containers    |

---

### 📘 README Updated To Include:

* Redis in magic-link workflow
* Kafka in event-driven workflow
* Updated architecture ASCII diagram
* Updated Quick Start (now includes Kafka)
* Updated troubleshooting section:

  * Kafka topic not found
  * Redis key debugging
  * MinIO permissions fix

---

# 🌱 **TASK-0.4 — Seed Scripts (Updated)**

### 🎯 Purpose

Populate initial data to support Sprint 1–4 workflows.

### 🗄️ Updated MongoDB Seeds

* Pharmacies (with scoring metadata):

  * accepted insurers
  * specialties
  * handling capacity
  * load factor
  * geo coordinates
* Patients (sample data)
* Sample prescriptions with updated statuses:

  * `received`
  * `validated`
  * `validation_issue`
  * `awaiting_enrollment`
  * `awaiting_routing`

### 🆕 Added Seed:

* Kafka topic initializer script (optional)
* PostgreSQL `job_queue` upgraded schema for:

  * exponential backoff
  * dead-letter queue (DLQ)
* Audit logs schema fully updated
* Redis dummy tokens for testing magic links

---

# ⚙️ **TASK-0.5 — CI/CD Skeleton (Updated)**

### 🎯 Purpose

Ensure backend + frontend + infrastructure builds correctly.

### 🔁 Pipeline Improvements

* Added Kafka container build check
* Added Redis integration test (Ping)
* Run Go tests with race detector:

  ```
  go test -race ./...
  ```
* Build multi-stage Docker images for:

  * API
  * Worker
  * Frontend
* Added checks:

  * `docker-compose config` syntax validation
  * Lint YAML files
  * Verify seeds syntax (JS + SQL)

### 📦 Artifacts Produced

* Backend test coverage
* Frontend test coverage
* Docker image build logs
* Linting reports

---

# 🚀 Getting Started (UPDATED)

### Step 1 — Start All Services

```
make dev
```

### Step 2 — Verify Kafka

```
make kafka-topics
```

Expect:

* intake_received
* validate_prescription
* enrollment_completed
* pharmacy_recommendation_requested

### Step 3 — Verify Redis

```
make redis-cli
> keys *
```

### Step 4 — Verify MinIO

Open browser → [http://localhost:9001](http://localhost:9001)
Default credentials:

```
admin / minioadmin
```

---

# 📘 Sprint 0 Summary (Updated)

| Component                    | Status | Notes                                  |
| ---------------------------- | ------ | -------------------------------------- |
| Monorepo structure           | ✅      | Includes kafka/, redis/, scheduler/    |
| Dockerized environment       | ✅      | Full infra w/ health checks            |
| Seeds for Mongo + PG + Redis | ✅      | Includes pharmacy scoring + tokens     |
| CI/CD pipeline               | ✅      | Includes Docker, Go, Node, Kafka tests |
| Updated README               | ✅      | Full architecture + quick start        |

