# 🔄 **PhilMyMeds — Improved End-to-End Architecture**

## **Architecture Overview**

PhilMyMeds follows a **microservices architecture** where:
- **Pharmacy handles manufacturer program adjudication** (not our enrollment system)
- **Kafka** provides event-driven communication
- **Redis** handles caching and session management
- **PostgreSQL** manages job queues and audit logs
- **MongoDB** stores business data

---

## **1. Prescription Intake (NCPDP Entry Point)**

### **1.1 Provider Submits Prescription**
- **Real world**: eRx via NCPDP SCRIPT standard
- **PhilMyMeds**: Mock data or Gemini-generated NCPDP payload

### **1.2 API Receives Prescription**
```
POST /api/v1/prescriptions/intake
```

**Process:**
1. Parse NCPDP SCRIPT XML format
2. Extract patient, prescriber, medication, insurance data
3. Store in **MongoDB** `prescriptions` collection
4. Initial status: `"received"`
5. Cache recent prescription in **Redis** (5-min TTL) for deduplication
6. Create **Validation Job** in PostgreSQL `validation_jobs` table

**Kafka Event:**
```json
Topic: "prescription.intake.received"
{
  "event_id": "evt_xyz",
  "prescription_id": "rx_abc123",
  "patient_id": "pat_def456",
  "drug_ndc": "00002-7510-02",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## **2. Validation Worker Processes Intake**

### **2.1 Validation Worker (Go Background Service)**

**Trigger**: Polls PostgreSQL `validation_jobs` table every 10 seconds

**Process:**
1. Fetch jobs with `status = 'pending'` using row-level locking
2. Parse and validate prescription:
   - **NPI Validation**: Check prescriber NPI against registry
   - **DEA Validation**: Verify DEA number for controlled substances
   - **NDC Validation**: Confirm drug code exists in FDA database
   - **SIG Validation**: Check directions format (route, frequency, quantity)
   - **Required Fields**: Ensure all mandatory NCPDP fields present
   - **Patient Demographics**: Verify name/DOB consistency

**Outcomes:**

**If Valid:**
1. Update MongoDB: `status = "validated"`
2. Add validation checks to prescription document
3. Mark job as `completed` in PostgreSQL
4. Log to PostgreSQL `audit_logs`
5. Publish Kafka event

```json
Topic: "prescription.validation.completed"
{
  "event_id": "evt_abc",
  "prescription_id": "rx_abc123",
  "validation_result": "passed",
  "timestamp": "2024-01-15T10:31:00Z"
}
```

**If Invalid:**
1. Update MongoDB: `status = "validation_failed"`
2. Store error details in `validation_errors` array
3. Mark job as `failed` in PostgreSQL
4. **No Kafka event** - Ops team handles manually via dashboard
5. Send alert notification to ops team

---

## **3. Patient Enrollment Flow**

### **3.1 Ops Team Triggers Enrollment**

**Trigger**: Ops reviews validated prescription in dashboard and clicks "Start Enrollment"

**API Endpoint:**
```
POST /api/v1/enrollment/initiate
Body: { "prescription_id": "rx_abc123" }
```

**Process:**
1. Generate unique magic link token (UUID)
2. Store token in **Redis** with 48-hour TTL:
   ```
   Key: "magic_link:{token}"
   Value: {
     "prescription_id": "rx_abc123",
     "patient_id": "pat_def456",
     "expires_at": "2024-01-17T10:35:00Z",
     "used": false
   }
   TTL: 48 hours
   ```
3. Create enrollment record in MongoDB with `status = "pending"`
4. Send email via **SendGrid** and SMS via **Twilio**:
   - Magic link URL: `https://enroll.phil-my-meds.com/enroll/{token}`
   - Link expires in 48 hours

### **3.2 Patient Opens Magic Link**

**Frontend Route**: `/enroll/:token`

**Validation Flow:**
```
GET /api/v1/enrollment/validate/:token
```

1. Check Redis for token existence and expiration
2. If valid, return prescription and patient data
3. If invalid/expired, show error message

### **3.3 Patient Completes Enrollment**

**React Enrollment Portal Collects:**

1. **Insurance Information** (verify/update):
   - Insurance carrier name
   - Member ID
   - Group number
   - BIN (6 digits)
   - PCN
   - Plan type detection (commercial vs government)

2. **Insurance Card Images**:
   - Upload front and back images
   - Store in **MinIO** with encryption
   - Return URLs to frontend

3. **HIPAA Consent**:
   - Display HIPAA authorization text
   - Collect electronic signature (canvas drawing)
   - Capture signature name
   - Record timestamp and IP address

**NO Manufacturer Program Selection** - This is handled by pharmacy during adjudication

**Submission:**
```
POST /api/v1/enrollment/submit
Body: {
  "prescription_id": "rx_abc123",
  "token": "magic_link_token",
  "insurance": {
    "payer_name": "Blue Cross Blue Shield",
    "member_id": "ABC123456789",
    "group_number": "12345",
    "bin": "610014",
    "pcn": "MEDDCO",
    "card_front_url": "minio://...",
    "card_back_url": "minio://..."
  },
  "hipaa_consent": {
    "authorization_text": "I authorize...",
    "signature": "data:image/png;base64,...",
    "signature_name": "John Doe",
    "signature_date": "2024-01-15T14:00:00Z",
    "ip_address": "192.168.1.1"
  }
}
```

**Process:**
1. Validate token from Redis
2. Mark token as `used` in Redis
3. Store/update insurance profile in MongoDB
4. Store enrollment data in MongoDB
5. Update prescription: `status = "enrolled"`
6. Log to PostgreSQL audit_logs (HIPAA compliance)
7. Publish Kafka event

```json
Topic: "patient.enrollment.completed"
{
  "event_id": "evt_ghi",
  "prescription_id": "rx_abc123",
  "enrollment_id": "enr_xyz789",
  "insurance_verified": true,
  "timestamp": "2024-01-15T14:05:00Z"
}
```

---

## **4. Pharmacy Routing & Selection**

### **4.1 Routing Worker**

**Trigger**: Consumes `patient.enrollment.completed` Kafka event

**Process:**
1. Create routing job in PostgreSQL
2. Fetch prescription with patient and insurance data
3. Query MongoDB `pharmacies` collection
4. Apply filtering criteria:
   - **Geographic**: Calculate distance from patient address
   - **Insurance Network**: Pharmacy accepts patient's insurance
   - **Specialty Capability**: Pharmacy can handle the drug (NDC)
   - **Capacity**: Check current workload via Redis cache
   - **Cold Chain**: If drug requires refrigeration
   - **Controlled Substances**: If pharmacy has DEA license

5. Check pharmacy capacity from Redis:
   ```
   Key: "pharmacy_capacity:{pharmacy_id}"
   Value: {
     "current_daily_rx": 45,
     "max_daily_rx": 100,
     "utilization": 0.45
   }
   TTL: 5 minutes
   ```

6. Score pharmacies using weighted algorithm:
   - Distance (30%)
   - Insurance network tier (25%)
   - Available capacity (20%)
   - Performance metrics (15%)
   - Fulfillment speed (10%)

7. Store top 5 recommendations in MongoDB
8. Update prescription: `status = "awaiting_pharmacy_selection"`

### **4.2 Ops Team Selects Pharmacy**

**Dashboard Display:**
- Shows ranked list with scores breakdown
- Displays pharmacy details, distance, network tier, capacity

**API Endpoint:**
```
POST /api/v1/prescriptions/{id}/select-pharmacy
Body: {
  "pharmacy_id": "pharm_abc123",
  "selected_by": "ops_user_john",
  "selection_reason": "Best overall match"
}
```

**Process:**
1. Update MongoDB: `selected_pharmacy_id`, `status = "pharmacy_selected"`
2. Increment pharmacy capacity counter in Redis
3. Update MongoDB pharmacy capacity
4. Log selection to audit_logs
5. Publish Kafka event

```json
Topic: "pharmacy.selected"
{
  "event_id": "evt_jkl",
  "prescription_id": "rx_abc123",
  "pharmacy_id": "pharm_abc123",
  "pharmacy_name": "SpecialtyRx Pharmacy",
  "timestamp": "2024-01-15T14:30:00Z"
}
```

---

## **5. Insurance Adjudication (Pharmacy Handles This)**

### **5.1 Adjudication Worker**

**Trigger**: Consumes `pharmacy.selected` Kafka event

**Process:**
1. Create adjudication job in PostgreSQL
2. Send prescription details to pharmacy via API
3. **Pharmacy performs TWO-STEP adjudication**:

### **Step 1: Primary Insurance Claim**

Pharmacy submits claim to patient's insurance:
```
Pharmacy → Insurance Payer (Blue Cross)
NCPDP Transaction with patient insurance BIN/PCN

Response:
- Claim approved
- Drug cost: $6,500
- Insurance pays: $5,000
- Patient copay: $1,500
```

### **Step 2: Secondary Claim to Manufacturer Program**

**Our System's Role:**
- Provide pharmacy with manufacturer program credentials
- Query MongoDB `manufacturer_programs` by drug NDC
- Check Redis cache first (1-hour TTL)
- Return program BIN/PCN to pharmacy

```
GET /api/v1/programs/lookup
Query: { "ndc": "00002-7510-02", "insurance_type": "commercial" }

Response:
{
  "program_id": "prog_abbvie_humira_2024",
  "program_name": "Humira Complete Savings Card",
  "credentials": {
    "bin": "004682",
    "pcn": "CNRX",
    "group_id": "HUMIRA"
  },
  "eligibility": {
    "commercial_only": true,
    "target_copay": 5.00
  }
}
```

**Pharmacy submits secondary claim:**
```
Pharmacy → Manufacturer Copay Hub (AbbVie)
NCPDP Transaction with program BIN/PCN

Response:
- Program approved
- Primary copay: $1,500
- Discount applied: $1,495
- Final patient copay: $5
```

### **5.2 Pharmacy Returns Results**

Pharmacy sends complete adjudication results back to our system:

```
POST /api/v1/adjudication/results
Body: {
  "prescription_id": "rx_abc123",
  "primary_insurance": {
    "claim_id": "CLM123",
    "status": "approved",
    "drug_cost": 6500.00,
    "insurance_paid": 5000.00,
    "patient_copay": 1500.00
  },
  "manufacturer_programs": [
    {
      "program_id": "prog_abbvie_humira_2024",
      "program_name": "Humira Complete Savings Card",
      "status": "approved",
      "discount_amount": 1495.00,
      "reduced_copay": 5.00
    }
  ],
  "cost_breakdown": {
    "total_drug_cost": 6500.00,
    "insurance_covered": 5000.00,
    "initial_copay": 1500.00,
    "manufacturer_discount": 1495.00,
    "final_patient_copay": 5.00
  }
}
```

**Process:**
1. Store adjudication results in MongoDB `adjudications` collection
2. Update prescription with cost breakdown
3. Update prescription: `status = "adjudicated"`
4. Mark adjudication job as complete
5. Cache adjudication result in Redis (30-min TTL)
6. Publish Kafka event

```json
Topic: "insurance.adjudication.completed"
{
  "event_id": "evt_mno",
  "prescription_id": "rx_abc123",
  "final_copay": 5.00,
  "programs_applied": ["Humira Complete Savings Card"],
  "patient_savings": 1495.00,
  "timestamp": "2024-01-15T15:00:00Z"
}
```

### **5.3 Prior Authorization Handling**

**If Insurance Requires PA:**

1. Pharmacy system detects PA requirement in initial claim response
2. Pharmacy notifies our system via API
3. Our system creates PA record in MongoDB
4. Ops team contacts prescriber to initiate PA request
5. Track PA status: pending → approved/denied
6. Once approved, pharmacy resubmits claim
7. Continue with manufacturer program adjudication

---

## **6. Payment Collection**

### **6.1 Payment Link Generation**

**Trigger**: Consumes `insurance.adjudication.completed` Kafka event

**API Process:**
```
POST /api/v1/payments/create-link
Body: {
  "prescription_id": "rx_abc123",
  "amount": 5.00,
  "currency": "USD"
}
```

**Process:**
1. Create Stripe Checkout Session
2. Store in MongoDB `payments` collection:
   - `stripe_session_id`
   - `stripe_payment_link`
   - `amount: 5.00`
   - `status: "pending"`
   - `link_expires_at`: 24 hours from now
3. Publish Kafka event

```json
Topic: "payment.link.created"
{
  "event_id": "evt_pqr",
  "prescription_id": "rx_abc123",
  "payment_id": "pay_xyz",
  "amount": 5.00,
  "timestamp": "2024-01-15T15:05:00Z"
}
```

### **6.2 Send Payment Notification**

**Notification Service** consumes `payment.link.created`:

**Email Template:**
```
Subject: Complete Your Payment - $5.00 Copay

Hi John,

Your prescription for Humira is ready! Here's your cost breakdown:

Original Drug Cost:        $6,500.00
Insurance Coverage:        $5,000.00
Initial Copay:             $1,500.00
Manufacturer Discount:     $1,495.00
                          -----------
Your Final Copay:              $5.00
You Saved:                 $1,495.00

[Pay $5.00 Now] → {stripe_payment_link}

This link expires in 24 hours.
```

**SMS:**
```
PhilMyMeds: Your Humira Rx is ready. Final copay: $5.00 (saved $1,495!). 
Pay now: {short_link}
```

### **6.3 Patient Completes Payment**

1. Patient clicks link → Stripe Checkout
2. Enters credit card info
3. Completes payment

### **6.4 Stripe Webhook Handler**

```
POST /api/v1/webhooks/stripe
```

**Process:**
1. Verify Stripe webhook signature
2. Handle `checkout.session.completed` event
3. Update MongoDB payment: `status = "paid"`, `paid_at = NOW()`
4. Update prescription: `status = "paid"`
5. Cache payment receipt in Redis (24-hour TTL)
6. Log payment to audit_logs
7. Publish Kafka event

```json
Topic: "payment.completed"
{
  "event_id": "evt_stu",
  "prescription_id": "rx_abc123",
  "payment_id": "pay_xyz",
  "amount_paid": 5.00,
  "timestamp": "2024-01-15T15:30:00Z"
}
```

### **6.5 Payment Monitoring Worker**

**Background Job**: Monitors payment timeouts

1. Polls PostgreSQL `payment_jobs` every 30 seconds
2. Checks for expired payment links (24 hours)
3. Sends reminder notification 2 hours before expiration
4. Marks prescription as `"payment_timeout"` if expired
5. Notifies ops team for follow-up

---

## **7. Shipping & Fulfillment**

### **7.1 Shipping Worker**

**Trigger**: Consumes `payment.completed` Kafka event

**Process:**
1. Create shipping job in PostgreSQL
2. Fetch prescription and pharmacy data
3. Call **Shippo API** to create shipment:

```javascript
// Shippo Integration
{
  "address_from": {
    "name": "SpecialtyRx Pharmacy",
    "street1": "123 Main St",
    "city": "Boston",
    "state": "MA",
    "zip": "02101"
  },
  "address_to": {
    "name": "John Doe",
    "street1": "456 Oak Ave",
    "city": "Cambridge",
    "state": "MA",
    "zip": "02138"
  },
  "parcels": [{
    "length": "10",
    "width": "8",
    "height": "4",
    "distance_unit": "in",
    "weight": "1",
    "mass_unit": "lb"
  }],
  "extra": {
    "signature_confirmation": "STANDARD" // For controlled substances
  }
}
```

**Shippo Response:**
```json
{
  "shipment_id": "shp_abc",
  "rates": [...],
  "label_url": "https://shippo.com/label.pdf",
  "tracking_number": "1Z999AA10123456784",
  "tracking_url": "https://www.ups.com/track?tracknum=1Z999..."
}
```

4. Store shipment in MongoDB `shipments` collection
5. Update prescription: `status = "shipped"`
6. Publish Kafka event

```json
Topic: "shipment.label.created"
{
  "event_id": "evt_vwx",
  "prescription_id": "rx_abc123",
  "shipment_id": "ship_xyz",
  "tracking_number": "1Z999AA10123456784",
  "carrier": "UPS",
  "timestamp": "2024-01-15T16:00:00Z"
}
```

### **7.2 Shipping Notification**

**Notification Service** consumes event and sends:

**Email:**
```
Subject: Your Prescription Has Shipped!

Hi John,

Your Humira prescription has been shipped via UPS.

Tracking Number: 1Z999AA10123456784
Estimated Delivery: January 18, 2024

Track Your Package: {tracking_url}
```

---

## **8. Delivery Tracking**

### **8.1 Delivery Tracking Worker**

**Background Job**: Polls Shippo API every 60 seconds

**Process:**
1. Fetch active shipments from PostgreSQL `tracking_jobs`
2. Call Shippo tracking API for status updates
3. Process tracking events:
   - `in_transit` → Update shipment status
   - `out_for_delivery` → Send notification to patient
   - `delivered` → Mark prescription as `"completed"`
   - `exception` → Alert ops team

4. Store tracking history in MongoDB shipment document
5. Update prescription status based on tracking

### **8.2 Delivery Confirmation**

**When Status = "delivered":**

1. Update MongoDB: `prescription.status = "completed"`
2. Update MongoDB: `shipment.status = "delivered"`
3. Write final audit log to PostgreSQL
4. Publish Kafka event

```json
Topic: "shipment.delivered"
{
  "event_id": "evt_yz1",
  "prescription_id": "rx_abc123",
  "delivered_at": "2024-01-18T14:30:00Z",
  "timestamp": "2024-01-18T14:35:00Z"
}
```

5. Send delivery confirmation notification

---

## **9. System Components**

### **9.1 Services**

| Service | Port | Purpose |
|---------|------|---------|
| **API Gateway** | 8080 | JWT auth, rate limiting, request routing |
| **Prescription Service** | 8081 | Prescription CRUD, status management |
| **Enrollment Service** | 8082 | Magic links, enrollment flow |
| **Routing Service** | 8083 | Pharmacy scoring and selection |
| **Program Service** | 8084 | Manufacturer program lookup |
| **Adjudication Service** | 8085 | Coordinates with pharmacy for claims |
| **Payment Service** | 8086 | Stripe integration, payment tracking |
| **Shipping Service** | 8087 | Shippo integration, label generation |
| **Notification Service** | 8088 | Email/SMS via SendGrid/Twilio |

### **9.2 Background Workers**

| Worker | Poll Interval | Job Table |
|--------|---------------|-----------|
| **Validation Worker** | 10 seconds | `validation_jobs` |
| **Enrollment Monitor** | 30 seconds | `enrollment_jobs` |
| **Routing Worker** | 20 seconds | `routing_jobs` |
| **Adjudication Worker** | 15 seconds | `adjudication_jobs` |
| **Payment Monitor** | 30 seconds | `payment_jobs` |
| **Shipping Worker** | 30 seconds | `shipping_jobs` |
| **Delivery Tracker** | 60 seconds | `tracking_jobs` |

### **9.3 Data Stores**

**MongoDB Collections:**
- `prescriptions` - Main prescription documents
- `patients` - Patient demographics
- `prescribers` - Healthcare providers
- `pharmacies` - Partner pharmacies
- `insurance_profiles` - Patient insurance info
- `manufacturer_programs` - Copay program catalog
- `enrollments` - Patient enrollment records
- `adjudications` - Insurance claim results
- `prior_authorizations` - PA tracking
- `payments` - Stripe payment records
- `shipments` - Shippo shipment data
- `notifications` - Communication log
- `users` - Ops team accounts
- `file_assets` - Insurance cards, labels, etc.

**PostgreSQL Tables:**
- `validation_jobs` - Validation queue
- `enrollment_jobs` - Enrollment monitoring
- `routing_jobs` - Pharmacy routing queue
- `adjudication_jobs` - Adjudication queue
- `payment_jobs` - Payment monitoring
- `shipping_jobs` - Shipping queue
- `tracking_jobs` - Delivery tracking
- `audit_logs` - HIPAA-compliant audit trail

**Redis Keys:**
- `rate_limit:{identifier}` - API rate limiting
- `rx:recent:{id}` - Deduplication cache (5 min)
- `rx:dedup:{hash}` - Duplicate detection (5 min)
- `magic_link:{token}` - Enrollment tokens (48 hours)
- `pharmacy_capacity:{id}` - Real-time capacity (5 min)
- `programs:ndc:{ndc}` - Program cache (1 hour)
- `adjudication:{id}` - Adjudication results (30 min)
- `payment_receipt:{id}` - Payment receipts (24 hours)
- `session:{session_id}` - JWT sessions

### **9.4 Kafka Topics**

| Topic | Producer | Consumers |
|-------|----------|-----------|
| `prescription.intake.received` | Prescription Service | Validation Worker |
| `prescription.validation.completed` | Validation Worker | Dashboard, Analytics |
| `patient.enrollment.completed` | Enrollment Service | Routing Worker |
| `pharmacy.selected` | Routing Service | Adjudication Worker |
| `insurance.adjudication.completed` | Adjudication Service | Payment Service |
| `payment.link.created` | Payment Service | Notification Service |
| `payment.completed` | Payment Service | Shipping Worker |
| `shipment.label.created` | Shipping Service | Notification Service |
| `shipment.delivered` | Shipping Service | Notification Service, Analytics |

---

## **10. Key Architecture Decisions**

### **10.1 Why Pharmacy Handles Program Adjudication**

✅ **Matches Real-World Process**: This is how PHIL Inc and all specialty pharmacies actually work

✅ **Simplifies Patient Experience**: Patient doesn't need to understand complex programs

✅ **Real-Time Pricing**: Patient knows exact final copay immediately

✅ **Automatic Application**: Pharmacy tries all available programs without patient action

✅ **NCPDP Standard**: Uses industry-standard claim routing

### **10.2 Why Kafka for Events**

✅ **Decoupling**: Services don't directly depend on each other

✅ **Reliability**: Events persisted, can replay if needed

✅ **Real-Time Updates**: Dashboard subscribes to all events for live updates

✅ **Audit Trail**: Every state change captured

✅ **Scalability**: Easy to add new consumers

### **10.3 Why Redis for Caching**

✅ **Performance**: Reduces MongoDB queries by 80%+

✅ **Session Management**: Fast JWT session lookups

✅ **Rate Limiting**: Prevents API abuse

✅ **Magic Links**: Fast token validation with automatic expiration

✅ **Capacity Tracking**: Real-time pharmacy workload without DB hits

### **10.4 Why PostgreSQL for Jobs**

✅ **ACID Transactions**: Reliable job locking and state management

✅ **Row-Level Locking**: Multiple workers don't process same job

✅ **Audit Trail**: Immutable log for HIPAA compliance

✅ **Query Performance**: Fast queries for pending jobs

✅ **Reliability**: Jobs persist across service restarts

---

## **11. HIPAA Compliance Features**

### **11.1 Data Encryption**

- **At Rest**: MongoDB encryption, MinIO server-side encryption
- **In Transit**: TLS 1.3 for all API calls
- **PHI Fields**: Insurance cards, patient data encrypted

### **11.2 Audit Logging**

PostgreSQL `audit_logs` table captures:
- Every status change
- Who made the change (user ID or "system")
- Timestamp
- IP address
- Before/after state

### **11.3 Access Controls**

- JWT authentication for ops team
- Magic links for patients (time-limited, single-use)
- Role-based permissions (admin, ops_manager, ops_agent)
- API rate limiting per user

### **11.4 Data Retention**

- Active prescriptions: Indefinite
- Completed prescriptions: 7 years (HIPAA requirement)
- Audit logs: 7 years
- Temporary data (Redis): Auto-expires based on TTL

---

## **12. Monitoring & Observability**

### **12.1 Metrics**

- API response times
- Worker job processing latency
- Kafka lag (events waiting to be processed)
- Redis hit/miss ratio
- Database connection pool usage
- External API failures (Stripe, Shippo)

### **12.2 Logging**

- Structured JSON logs (Zap or Zerolog)
- Log levels: DEBUG, INFO, WARN, ERROR
- Correlation IDs for request tracing
- ELK Stack for log aggregation (optional)

### **12.3 Alerting**

- Job failures after max retries
- Worker crashes
- Stripe webhook failures
- Shippo API errors
- Payment link expirations
- Prescription timeout thresholds

---

## **End Result**

A **production-ready, HIPAA-compliant prescription fulfillment system** that:

✅ Handles complete workflow from intake to delivery

✅ Automatically finds maximum patient savings via manufacturer programs

✅ Provides real-time visibility to operations team

✅ Integrates with real-world pharmacy systems

✅ Scales horizontally with microservices architecture

✅ Maintains complete audit trail for compliance

✅ Provides excellent patient experience with transparent pricing