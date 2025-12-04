# 📘 Product Requirements — User Stories & Acceptance Criteria

## **Epic: Intake & Validation**

---

### **US-INT-01 — Prescription Intake (NCPDP SCRIPT)**

**As an** HCP system (simulator)
**I want to** submit prescriptions in NCPDP SCRIPT format
**So that** the system can ingest them.

#### ✅ Acceptance Criteria

* `POST /api/intake` accepts a valid NCPDP SCRIPT payload.
* System creates a **Prescription** record in **MongoDB** with status `received`.
* Returns **400 Bad Request** on malformed SCRIPT payloads.

#### 📌 Requirements

* Parse NCPDP SCRIPT XML.
* API authentication (if needed later).
* MongoDB schema for Prescription.
* Error handling & validation for malformed XML.

---

### **US-INT-02 — Automated Validation Worker**

**As the system**, I want an automated validation worker to validate required fields so that bad prescriptions are flagged.

#### ✅ Acceptance Criteria

* Worker picks validation job.
* Runs rules:

  * NPI present
  * NDC present
  * Quantity present
  * Directions present
* Updates prescription status to:

  * `validated`
  * `validation_issue` (with list of errors stored in MongoDB)
* Retry policy: **3 retries**.

#### 📌 Requirements

* Background worker queue (e.g., BullMQ / Go worker).
* Validation rules engine.
* Error list persistence in MongoDB.
* Retry + backoff configuration.

---

## **Epic: Patient Enrollment**

---

### **US-ENR-01 — Magic Enrollment Link**

**As the system**, create a magic enrollment link for new patients so they can complete insurance & consent.

#### ✅ Acceptance Criteria

* Magic link generated with expiry (e.g., 72 hours).
* Link sent via **SendGrid** (email) or **Twilio** (SMS).
* Visiting the link opens an **Enrollment SPA route**.
* Route displays pre-filled patient data.

#### 📌 Requirements

* Token generator with expiry.
* Magic link table/schema.
* SendGrid/Twilio integration.
* SPA route capable of loading initial patient data.

---

### **US-ENR-02 — Insurance Card Upload**

**As a patient**, I want to upload an insurance card image so ops can verify coverage.

#### ✅ Acceptance Criteria

* Uploaded image stored in local S3-compatible storage (e.g., MinIO).
* File metadata linked to enrollment record.
* Validates file type + size.
* Optional: OCR placeholder for future processing.

#### 📌 Requirements

* File upload endpoint.
* MinIO (or AWS S3) connection for development.
* MIME & size validation.
* Update enrollment record in DB.

---

## **Epic: Pharmacy Routing & Ops**

---

### **US-PHR-01 — Generate Pharmacy Recommendations**

**As the system**, generate ranked pharmacy recommendations based on location, capacity, and insurance network.

#### ✅ Acceptance Criteria

* API returns **top N pharmacies**.
* Each recommendation includes scoring details.
* Ops dashboard displays ranked list.

#### 📌 Requirements

* Scoring engine (location, network, capacity).
* Pharmacy dataset.
* Recommendation API.

---

### **US-PHR-02 — Manual Pharmacy Selection**

**As an operator**, I want to manually select a pharmacy so the workflow can proceed.

#### ✅ Acceptance Criteria

* Selecting a pharmacy updates prescription status to `pharmacy_selected`.
* Triggers adjudication job automatically.

#### 📌 Requirements

* Ops dashboard UI for selecting pharmacy.
* DB update logic.
* Message/queue trigger for adjudication.

---

## **Epic: Insurance Adjudication**

---

### **US-ADJ-01 — Adjudication Simulator**

**As a pharmacy simulator**, run an adjudication call and return cost breakdown.

#### ✅ Acceptance Criteria

* Adjudication API returns:

  ```json
  {
    "reimbursement": number,
    "coupon_amount": number,
    "patient_copay": number
  }
  ```
* System stores results.
* System displays breakdown to patient.

#### 📌 Requirements

* Fake adjudication microservice or in-process simulator.
* DB persistence of adjudication results.

---

## **Epic: Payments**

---

### **US-PAY-01 — Stripe Payment Link**

**As the system**, generate Stripe payment link and send to patient.

#### ✅ Acceptance Criteria

* Stripe Test Payment link generated.
* Link sent via email/SMS.
* Webhook updates payment record to `paid`.
* Triggers fulfillment job after payment.

#### 📌 Requirements

* Stripe API integration.
* Payment status tracking (DB).
* Webhook endpoint.
* Fulfillment queue trigger.

---

## **Epic: Shipping**

---

### **US-SHP-01 — Shippo Label Creation**

**As the pharmacy**, create a Shippo label and store tracking number; notify patient.

#### ✅ Acceptance Criteria

* Shippo label created in dev mode.
* Tracking number saved.
* Status updated to `shipped`.
* Tracking updates reflected in system.
* When carrier marks delivered → status becomes `delivered`.

#### 📌 Requirements

* Shippo dev mode integration.
* Shipping DB schema.
* Webhook or polling for tracking updates.

---

## **Epic: Audit & Logging**

---

### **US-AUD-01 — Audit Log of State Changes**

**As an auditor**, I want every state change recorded with timestamp and context.

#### ✅ Acceptance Criteria

* Every status change writes an entry to `audit_logs` (PostgreSQL):

  * timestamp
  * actor
  * previous status
  * new status
  * payload
* Logs searchable in Ops dashboard.

#### 📌 Requirements

* Postgres table for audit logs.
* Middleware or event hook to log changes.
* Dashboard query UI.

