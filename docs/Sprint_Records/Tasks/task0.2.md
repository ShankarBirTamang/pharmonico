# 🚀 **TASK-0.2 — Create Docker Compose Dev Environment (Enhanced & Expanded Version)**

## **Overview**

This task creates the **entire local development infrastructure** for Pharmonico using Docker Compose.
It includes:

* Databases (MongoDB, PostgreSQL)
* Cache (Redis)
* Event streaming backbone (Kafka + Zookeeper)
* Storage (MinIO)
* Email testing (Maildev)
* Application services (API, Worker, Frontend)
* Optional reverse proxy (Nginx)

The environment must be:

* **Fully isolated**
* **Ready for hot-reload development**
* **Consistent across machines**
* **Easy to tear down and rebuild**
* **Checked using healthchecks to avoid race conditions**

---

# 📦 **Sub-Task Breakdown (Expanded)**

| Sub-Task  | Description                                            | Dependencies  |
| --------- | ------------------------------------------------------ | ------------- |
| **0.2.1** | Database Services: MongoDB + PostgreSQL                | None          |
| **0.2.2** | Storage Service: MinIO                                 | None          |
| **0.2.3** | Developer Tools: Maildev                               | None          |
| **0.2.4** | Go Backend Dockerfiles (API + Worker)                  | 0.2.1         |
| **0.2.5** | React Frontend Dockerfile                              | None          |
| **0.2.6** | Redis + Kafka Infrastructure (NEW, added)              | None          |
| **0.2.7** | Docker Compose Orchestration (Tie Everything Together) | 0.2.1–0.2.6   |
| **0.2.8** | Nginx Reverse Proxy (Optional)                         | 0.2.4 + 0.2.5 |

---

# 📌 **TASK-0.2.1 — Database Services (MongoDB + PostgreSQL)**

### ✔️ Purpose

Set up both databases with persistent volumes, authentication, and health checks.

---

### 🔧 MongoDB Configuration (Expanded)

```yaml
mongodb:
  image: mongo:7.0
  container_name: mongodb
  ports:
    - "27017:27017"
  volumes:
    - mongodb-data:/data/db
  healthcheck:
    test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
    interval: 10s
    timeout: 5s
    retries: 5
    start_period: 10s
  environment:
    - MONGO_INITDB_DATABASE=pharmonico
  networks:
    - pharmonico-network
  restart: unless-stopped
```

**Key Features Added:**

* Persistent named volume
* Startup healthcheck (avoids API connecting too early)
* Pre-created database `pharmonico`

---

### 🔧 PostgreSQL Configuration (Expanded)

```yaml
postgres:
  image: postgres:16-alpine
  container_name: postgres
  environment:
    - POSTGRES_USER=postgres
    - POSTGRES_PASSWORD=postgres
    - POSTGRES_DB=pharmonico
  ports:
    - "5432:5432"
  volumes:
    - postgres-data:/var/lib/postgresql/data
  healthcheck:
    test: ["CMD-SHELL", "pg_isready -U postgres"]
    interval: 10s
    timeout: 5s
    retries: 5
    start_period: 10s
  networks:
    - pharmonico-network
  restart: unless-stopped
```

**Improvements:**

* Proper healthchecks
* Persistent volume
* Environment variables externalized for production later

---

# 📌 **TASK-0.2.2 — Object Storage (MinIO)**

### ✔️ Purpose

Provide S3-compatible storage for uploads (insurance cards) & generated files (shipping labels).

---

### 🔧 MinIO + Bucket Initialization

```yaml
minio:
  image: minio/minio:latest
  container_name: minio
  environment:
    - MINIO_ROOT_USER=minioadmin
    - MINIO_ROOT_PASSWORD=minioadmin
  command: server /data --console-address ":9001"
  ports:
    - "9000:9000"
    - "9001:9001"
  volumes:
    - minio-data:/data
  networks:
    - pharmonico-network
  restart: unless-stopped

minio-init:
  image: minio/mc:latest
  depends_on:
    - minio
  entrypoint: >
    /bin/sh -c "
    sleep 5;
    mc alias set myminio http://minio:9000 minioadmin minioadmin;
    mc mb myminio/pharmonico --ignore-existing;
    mc anonymous set download myminio/pharmonico;
    "
  networks:
    - pharmonico-network
```

**Improvements Added:**

* Auto-create MinIO bucket on startup
* Allow anonymous download for dev mode
* Clear separation of object storage vs bucket-init logic

---

# 📌 **TASK-0.2.3 — Developer Tools (Maildev)**

### ✔️ Purpose

Capture emails in dev without sending real emails.

```yaml
maildev:
  image: maildev/maildev:2.1.0
  container_name: maildev
  ports:
    - "1080:1080"
    - "1025:1025"
  networks:
    - pharmonico-network
  restart: unless-stopped
```

**UI:** [http://localhost:1080](http://localhost:1080)
**SMTP:** localhost:1025

---

# 📌 **TASK-0.2.4 — Go Backend Dockerfiles (API + Worker)**

### ✔️ Purpose

Provide isolated containerized environments for:

* API Server
* Background Worker (for jobs + Kafka consumers)

Both are built using the same Dockerfile with multi-stage builds.

---

### 🛠️ backend-go/Dockerfile (Expanded + Improved)

Includes:

* dependency caching
* timezone support
* separation between API & Worker builds

*(You already posted the file; here is the expanded, final version)*

```dockerfile
# ---------- BUILD STAGE ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build API
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api

# Build Worker
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/worker ./cmd/worker


# ---------- API RUNTIME ----------
FROM alpine:3.19 AS api
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /bin/api /app/api
EXPOSE 8080
CMD ["/app/api"]


# ---------- WORKER RUNTIME ----------
FROM alpine:3.19 AS worker
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /bin/worker /app/worker
CMD ["/app/worker"]
```

---

# 📌 **TASK-0.2.5 — React Frontend Dockerfile (Improved)**

Includes hot reload mode in development.

*(Your existing file is solid; adding final prod details for completeness)*

```dockerfile
# Development mode
FROM node:20-alpine AS development
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 5173
CMD ["npm", "run", "dev"]


# Production build
FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build


# Production runtime (Nginx)
FROM nginx:alpine AS production
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

---

# 📌 **TASK-0.2.6 — Redis + Kafka Infrastructure (NEW Details Added)**

### ✔️ Purpose

Enable:

* Magic link tokens (Redis)
* Event-based pipelines (Kafka)
* Reliable background processing

---

### 🔧 Redis

```yaml
redis:
  image: redis:7-alpine
  container_name: redis
  ports:
    - "6379:6379"
  volumes:
    - redis-data:/data
  command: ["redis-server", "--save", "60", "1"]
  networks:
    - pharmonico-network
  restart: unless-stopped
```

---

### 🔧 Kafka + Zookeeper (Production-grade Local Setup)

```yaml
zookeeper:
  image: confluentinc/cp-zookeeper:7.5.0
  environment:
    ZOOKEEPER_CLIENT_PORT: 2181
    ZOOKEEPER_TICK_TIME: 2000
  ports:
    - "2181:2181"
  networks:
    - pharmonico-network

kafka:
  image: confluentinc/cp-kafka:7.5.0
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
```

### Kafka auto-topic creation can be added later.

---

# 📌 **TASK-0.2.7 — Docker Compose Orchestration (Master File)**

This ties **all services together**:

```yaml
version: "3.8"

services:
  # (all services from Tasks 0.2.1–0.2.6 go here)

volumes:
  mongodb-data:
  postgres-data:
  minio-data:
  redis-data:

networks:
  pharmonico-network:
    driver: bridge
```

---

# 📌 **TASK-0.2.8 — Optional Nginx Reverse Proxy**

Used only if you want a **single entry point** (`http://localhost`).

```yaml
nginx:
  image: nginx:alpine
  ports:
    - "80:80"
  volumes:
    - ./infra/nginx/nginx.conf:/etc/nginx/nginx.conf:ro
  depends_on:
    - api
    - frontend
  networks:
    - pharmonico-network
```

---

# 🧪 **Verification Checklist (Expanded)**

| Sub-Task   | Verification Command                                                          |
| ---------- | ----------------------------------------------------------------------------- |
| MongoDB    | `docker compose exec mongodb mongosh --eval "db.adminCommand('ping')"`        |
| PostgreSQL | `docker compose exec postgres psql -U postgres -c "SELECT 1"`                 |
| Redis      | `docker compose exec redis redis-cli ping`                                    |
| Kafka      | `docker compose exec kafka kafka-topics --bootstrap-server kafka:9092 --list` |
| MinIO      | Visit [http://localhost:9001](http://localhost:9001)                          |
| Maildev    | Visit [http://localhost:1080](http://localhost:1080)                          |
| API        | Visit [http://localhost:8080/health](http://localhost:8080/health)            |
| Frontend   | Visit [http://localhost:5173](http://localhost:5173)                          |
| Nginx      | Visit [http://localhost](http://localhost) (if enabled)                       |

