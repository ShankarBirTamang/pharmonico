# ✅ **Pharmonico — Updated Sprint Plan (Final Version)**

This sprint plan includes **Redis caching**, **Kafka event workflow**, **MinIO**, **Stripe shipping**, **NCPDP intake**, and **worker system improvements**.

Copy and paste directly.

---

# 🚀 **SPRINT 0 — Monorepo + Infra Bootstrap**

### **TASK-0.1: Monorepo Initialization**

* Create project structure:

  ```
  backend-go/
  frontend-react/
  infra/
  docs/
  scripts/
  ```
* Add `.editorconfig`, `.gitignore`, and base folder scaffolding.

---

### **TASK-0.2: Docker Compose Developer Environment**

Provision core infra:

**Services**

* `golang-api` (Go API server)
* `worker` (background processing; Kafka consumer)
* `react-app`
* `mongodb`
* `postgres`
* `redis`
* `kafka + zookeeper`
* `minio`
* `maildev`
* `nginx` (optional reverse proxy)

Provide:

* Health checks
* Volume mounts
* Hot reload

---

### **TASK-0.3: Bootstrap Scripts + README**

* Add development scripts:

  * `make dev`
  * `make dev-down`
  * `make seed`
  * `make test`
* Add updated **architecture diagram**, including Kafka + Redis
* Add onboarding instructions in `README.md`

---

### **TASK-0.4: Seed Scripts**

* MongoDB seeds:

  * pharmacies
  * patients
  * sample prescriptions
* PostgreSQL seeds:

  * job queue schema (`jobs` table)
  * audit logs table
* Insert test data to validate flows.

---

### **TASK-0.5: CI Pipeline Skeleton**

* GitHub Actions:

  * Go Build + Test
  * React Build + Test
  * Linting (Go + JS)
  * Docker build check
* Cache dependencies

---

<br>

# 🚀 **SPRINT 1 — Intake + Validation Worker (+ Kafka Routing)**

### **TASK-1.1: POST `/api/intake`**

* Accept NCPDP-like payload
* Store in MongoDB
* Mark as `received`
* Link patient or auto-create patient record
* Publish `intake_received` Kafka event

**Requirements:**
*NCPDP SCRIPT intake structure + validation rules*

---

### **TASK-1.2: Kafka Event Topics + Minimal Job Enqueuer**

Topics:

* `intake_received`
* `validate_prescription`

Steps:

* When intake happens → publish Kafka event
* Create fallback enqueue to PostgreSQL job queue (resilience)

---

### **TASK-1.3: Validation Worker (Go)**

Worker responsibilities:

* Consume `validate_prescription` from Kafka
* Validate:

  * NPI
  * NDC
  * dosage form
  * duplicate check
  * patient demographics
* Update Mongo:

  * `validated`
  * OR `validation_issue`

**Requirements:**
*NCPDP validation rules, prescription structure*

---

### **TASK-1.4: Basic Ops UI**

React dashboard:

* **Intake Tab:**

  * Show all received intake items
  * Status: `received`, `validated`, `validation_issue`
* **Validation Tab:**

  * Show errors, warnings, metadata
  * Click item → detail drawer

---

<br>

# 🚀 **SPRINT 2 — Enrollment Flow + Magic Links + Redis**

### **TASK-2.1: Magic Link API (Token via Redis)**

Endpoint:

* `/api/enrollment/start`
  Process:
* Create Redis token (TTL: 30 mins)
* Attach patient + prescription ID
* Generate magic link URL
* Store audit entry

---

### **TASK-2.2: Email/SMS Integration**

* Use Maildev for local email
* Twilio test mode for SMS (or mock)
* Templates:

  * “Start Enrollment”
  * “Verify Insurance”
  * “Upload Insurance Card”

---

### **TASK-2.3: Enrollment SPA**

Routes:

* `/enroll/:token`
  Features:
* Token validation (Redis lookup)
* Pages:

  * Insurance details
  * Consent form
  * File upload → MinIO storage
* On success:

  * Update Mongo → `enrollment_completed`
  * Publish Kafka event `enrollment_completed`

**Requirements:**
*Enrollment workflow, authorization rules*

---

### **TASK-2.4: Enrollment Worker**

Responsibilities:

* Consume `enrollment_completed` topic
* Verify insurance submission
* Validate insurance data format
* Update Mongo → `ready_for_routing`
* Publish `pharmacy_recommendation_requested`

---

<br>

# 🚀 **SPRINT 3 — Pharmacy Routing + Adjudication Engine**

### **TASK-3.1: Pharmacy Scoring Engine**

Inputs:

* Distance (zip/county)
* Specialty match (e.g., Oncology)
* Current load/capacity
* Insurance acceptance
* Working hours (optional)

Outputs:

* Sorted pharmacy list
* Embed scores + metadata

---

### **TASK-3.2: Ops UI — Pharmacy Selection**

Add to React UI:

* Ranked list with scores
* Manual override button
* Save selected pharmacy
* Publish Kafka event `pharmacy_selected`

---

### **TASK-3.3: Mock Adjudication API**

Simulate:

* claim submission
* copay calculation
* coupon discount
* reimbursement

Store adjudication data in Mongo:

* copay
* plan paid
* patient responsibility

---

### **TASK-3.4: Adjudication Worker**

When event `run_adjudication` fires:

* Trigger mock adjudication
* Update Mongo: `adjudicated`
* Publish `payment_link_required`

---

<br>

# 🚀 **SPRINT 4 — Payments + Shipping + Webhooks + Notifications + Audit + Observability**

### **TASK-4.1: Stripe Payment Link**

Endpoint:

* `/api/payments/create-link`

Flow:

* Create Stripe checkout session
* Save session ID → Mongo
* Return payment URL
* Publish `awaiting_payment`

---

### **TASK-4.2: Stripe Webhook Handler**

Endpoint:

* `/webhook/stripe`
  Includes:
* Signature verification
* On payment success:

  * Mark as `paid`
  * Publish `start_shipping`

---

### **TASK-4.3: Shippo Shipping Integration**

Workflow:

* Create shipment
* Generate label
* Save tracking number → Mongo
* Publish `shipping_created`

---

### **TASK-4.4: Delivery Tracking Worker**

* Poll Shippo API (test mode)
* Status updates:

  * `in_transit`
  * `out_for_delivery`
  * `delivered`
* Notification logic
* Audit logs for updates

---

### **TASK-4.5: Audit Log System**

PostgreSQL tables:

* `audit_logs`
  Log every state transition:
* user
* timestamp
* event source
* metadata

UI addition:

* Audit Log viewer in Ops Dashboard

---

### **TASK-4.6: Notification System**

Email:

* enrollment started
* payment link
* shipping
* delivery
  SMS (optional):
* short alerts
  Use Maildev + Twilio mock

---

### **TASK-4.7: Ops Dashboard Finalization**

Add:

* Search (patient, Rx ID)
* Filters (status, pharmacy, insurance)
* Workflow timeline component
* Job failure view

---

### **TASK-4.8: Observability**

* Structured logs (Zap or Logrus)
* `/metrics` endpoint (Prometheus format)
* Worker failure counters
* Kafka lag monitor (optional)

---

### **TASK-4.9: E2E Testing**

Test full pipeline:

* Intake → Validation
* Enrollment → Routing
* Routing → Adjudication
* Payment → Shipping
* Delivery + Notifications
* Worker retries + crash recovery

---

### **TASK-4.10: Deployment Prep**

Optional but recommended:

* Multi-stage Dockerfiles
* `docker-compose.prod.yml`
* Environment templates:

  * `.env.example`
  * `infra/env/*`

---

# 📦 **SPRINT SUMMARY TABLE**

| Sprint       | Focus Area                                                                  |
| ------------ | --------------------------------------------------------------------------- |
| **Sprint 0** | Infra, Monorepo, Docker, Seeds, CI                                          |
| **Sprint 1** | Intake, Validation Worker, Kafka Events, Ops Intake UI                      |
| **Sprint 2** | Enrollment, Magic Links (Redis), Insurance, Enrollment Worker               |
| **Sprint 3** | Pharmacy Routing, Scoring Engine, Adjudication Pipeline                     |
| **Sprint 4** | Payments, Shipping, Notifications, Webhooks, Audits, Observability, Testing |

---

