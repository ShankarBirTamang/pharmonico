

# ✅ **SPRINT 1 — Intake & Validation **

**Goal:** Implement prescription intake, deduplication, parsing, validation pipeline, and basic Ops dashboard with real-time updates.

---

# **TASK 1.1 — NCPDP Intake API (12 hours)**

### **Backend API Development**

* **Subtask 1.1.1** — Create route `POST /api/v1/prescriptions/intake`
* **Subtask 1.1.2** — Parse request body (XML or JSON)
* **Subtask 1.1.3** — Implement NCPDP parser (`pkg/ncpdp/parser.go`)
* **Subtask 1.1.4** — Extract patient, prescriber, medication, insurance
* **Subtask 1.1.5** — Validate XML structure (well-formedness check)

### **Deduplication & Caching**

* **Subtask 1.1.6** — Generate dedup hash from core fields (patient + drug + date)
* **Subtask 1.1.7** — Check Redis for duplicates (TTL 5 mins)
* **Subtask 1.1.8** — Store dedup key in Redis if new

### **Persistence**

* **Subtask 1.1.9** — Insert prescription into MongoDB
* **Subtask 1.1.10** — Set status: `"received"`

### **Event Publishing**

* **Subtask 1.1.11** — Publish Kafka event: `prescription.intake.received`
* **Subtask 1.1.12** — Structure event payload with prescription metadata

### **API Response**

* **Subtask 1.1.13** — Return `{ prescription_id }`

### **Testing**

* **Subtask 1.1.14** — Test valid NCPDP XML intake
* **Subtask 1.1.15** — Test duplicate detection logic
* **Subtask 1.1.16** — Test malformed XML
* **Subtask 1.1.17** — Test missing required fields
* **Subtask 1.1.18** — Test Kafka publish success/failure paths

---

# **TASK 1.2 — Validation Worker (16 hours)**

### **Worker Setup**

* **Subtask 1.2.1** — Implement worker ticker loop (10-second interval)
* **Subtask 1.2.2** — Fetch pending jobs from PostgreSQL using row-level locking
* **Subtask 1.2.3** — Batch size = 10 jobs

### **Concurrent Job Processing**

* **Subtask 1.2.4** — Spawn goroutines for each job
* **Subtask 1.2.5** — Add wait group for concurrency control
* **Subtask 1.2.6** — Handle job context cancellation

### **Validation Rules Implementation**

* **Subtask 1.2.7** — Create `Validator` interface
* **Subtask 1.2.8** — Implement **NPI validation**

  * format (10-digit)
  * registry lookup
* **Subtask 1.2.9** — Implement **DEA validation**

  * format
  * checksum
  * control-substance compatibility
* **Subtask 1.2.10** — Implement **NDC validation**

  * format
  * FDA drug lookup
  * active/discontinued status
* **Subtask 1.2.11** — Add base rule: required fields (patient, prescriber, medication)

### **Outcome Handling**

* **Subtask 1.2.12** — Mark job `completed` or `failed`
* **Subtask 1.2.13** — Update MongoDB prescription:

  * status = `"validated"` OR `"validation_failed"`
  * store validation errors
* **Subtask 1.2.14** — Publish Kafka event: `prescription.validation.completed` (only when valid)
* **Subtask 1.2.15** — Trigger Ops notification when failed
* **Subtask 1.2.16** — Add audit log record

### **Testing**

* **Subtask 1.2.17** — Unit test each validator (NPI/DEA/NDC)
* **Subtask 1.2.18** — Test job retry logic
* **Subtask 1.2.19** — Test concurrency behavior
* **Subtask 1.2.20** — Integration test with PostgreSQL + MongoDB

---

# **TASK 1.3 — Basic Operations Dashboard (12 hours)**

### **Dashboard Layout**

* **Subtask 1.3.1** — Create Dashboard page structure with header + tabs
* **Subtask 1.3.2** — Add "Intake" tab
* **Subtask 1.3.3** — Add "Validation" tab

### **Intake Tab**

* **Subtask 1.3.4** — Hook: `usePrescriptions({ status: 'received' })`
* **Subtask 1.3.5** — Render `<PrescriptionCard />`
* **Subtask 1.3.6** — Add “Validate” button to trigger manual validation
* **Subtask 1.3.7** — Add loading and empty states

### **Validation Tab**

* **Subtask 1.3.8** — Hook: Fetch prescriptions with status `['validated', 'validation_failed']`
* **Subtask 1.3.9** — Show validation errors on failed items
* **Subtask 1.3.10** — Display success/failed labels

### **Real-Time WebSocket Updates**

* **Subtask 1.3.11** — Connect Dashboard to WebSocket server
* **Subtask 1.3.12** — Listen to Kafka → WebSocket → React updates
* **Subtask 1.3.13** — On event: `invalidateQueries(['prescriptions'])`
* **Subtask 1.3.14** — Test real-time UI refresh scenario

### **UI Testing**

* **Subtask 1.3.15** — Test card rendering
* **Subtask 1.3.16** — Test WebSocket reconnection handling
* **Subtask 1.3.17** — Test mobile responsiveness

---

# **SPRINT 1 Deliverables Checklist**

* ✅ NCPDP Intake API endpoint
* ✅ Redis duplicate detection
* ✅ MongoDB insert + Kafka event publishing
* ✅ Validation Worker with NPI/DEA/NDC rules
* ✅ PostgreSQL job queue integration
* ✅ Ops dashboard (Intake + Validation tabs)
* ✅ Real-time updates via WebSocket
* ✅ Unit & integration tests

