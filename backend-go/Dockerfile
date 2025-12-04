# =============================================================================
# Pharmonico Backend - Multi-stage Dockerfile
# Builds: API, Worker, and Scheduler services
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Builder - Compiles all Go binaries
# -----------------------------------------------------------------------------
FROM golang:1.21-alpine AS builder

# Install git and certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go.mod and go.sum first (for better layer caching)
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build API binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /bin/api ./cmd/api

# Build Worker binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /bin/worker ./cmd/worker

# Build Scheduler binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /bin/scheduler ./cmd/scheduler

# -----------------------------------------------------------------------------
# Stage 2a: API Runtime - Minimal image for API server
# -----------------------------------------------------------------------------
FROM alpine:3.19 AS api

# Add certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /bin/api .

# Use non-root user
USER appuser

# API listens on port 8080
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./api"]

# -----------------------------------------------------------------------------
# Stage 2b: Worker Runtime - Minimal image for Kafka consumer
# -----------------------------------------------------------------------------
FROM alpine:3.19 AS worker

RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /bin/worker .

USER appuser

CMD ["./worker"]

# -----------------------------------------------------------------------------
# Stage 2c: Scheduler Runtime - Minimal image for cron jobs
# -----------------------------------------------------------------------------
FROM alpine:3.19 AS scheduler

RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /bin/scheduler .

USER appuser

CMD ["./scheduler"]


