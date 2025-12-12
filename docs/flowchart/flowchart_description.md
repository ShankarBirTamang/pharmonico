Here is a **clear explanation of the PhilMyMeds flowchart**, based on the requirements, architecture, and workflow described inside both documents you uploaded.
(Your flowchart image wasn’t provided here, so I am explaining the *official* improved workflow from the docs you shared.)

---

# ✅ **PhilMyMeds Flowchart — Simple & Clear Explanation**

PhilMyMeds simulates a **real specialty pharmacy** prescription pipeline.
The full flow has **7 major stages**, each with workers, Kafka events, Redis caching, and MongoDB writes.

Below is the **easy-to-understand explanation** of each stage in the flowchart.

---

# **1️⃣ Prescription Intake (Stage 1)**

**Who sends it?**
— Doctor’s EHR (simulated by Gemini API)

**What happens?**

1. System receives NCPDP SCRIPT XML.
2. Parse patient, prescriber, medication, insurance.
3. Create prescription record in MongoDB.
4. Add to Redis for deduplication.
5. Push “prescription.intake.received” event to Kafka.
6. Create validation job in PostgreSQL.

Outcome → Status: **received**

---

# **2️⃣ Prescription Validation (Stage 2)**

**Worker picks up the validation job** and checks:

* Required fields
* NPI validation (doctor ID)
* DEA validation (if controlled drugs)
* NDC validation (drug code)
* SIG (directions)
* Patient demographics & insurance completeness

If valid →

* Update status to **validated**
* Emit “prescription.validation.completed” Kafka event

If invalid →

* Status = **validation_failed**
* Ops fixes manually

---

# **3️⃣ Patient Enrollment (Stage 3)**

Triggered by Ops user.

### Actions:

1. System generates **magic link** → stored in Redis with 48h TTL.
2. Patient opens the enrollment portal.
3. Patient submits:

   * Updated insurance
   * Insurance card uploads
   * HIPAA consent
   * Signature

### Important:

🚨 **No manufacturer program handling here**
(As per real-world process—pharmacy handles this during adjudication.)

Outcome → Status: **enrolled**
Kafka: **patient.enrollment.completed**

---

# **4️⃣ Pharmacy Routing & Selection (Stage 4)**

Triggered by routing worker after enrollment.

### Steps:

1. Fetch pharmacies matching:

   * Insurance network
   * Geographic proximity
   * Capacity (via Redis)
   * NDC compatibility (specialty or refrigerator storage etc.)

2. Create a ranked list using weighted scoring.

3. Ops team manually selects pharmacy.

Outcome → Status: **pharmacy_selected**
Kafka: **pharmacy.selected**

---

# **5️⃣ Insurance Adjudication (Stage 5)**

🔥 **This is where manufacturer program matching happens**
— **BY THE PHARMACY**, not by our system.

Pharmacy performs:

### Step 1 — Primary Insurance Claim

Outputs:

* Insurance paid
* Patient initial copay
* Drug cost

### Step 2 — Manufacturer Program Claim

(Using manufacturer program BIN/PCN)

Outputs:

* Discount applied
* Final adjusted copay

Pharmacy sends complete breakdown back to PhilMyMeds.

Outcome → Status: **adjudicated**
Kafka: **insurance.adjudication.completed**

---

# **6️⃣ Payment Collection (Stage 6)**

Triggered automatically after adjudication.

### Steps:

1. Create Stripe payment link
2. Save payment record
3. Send email/SMS to patient
4. Stripe webhook marks payment complete

Outcome → Status: **paid**
Kafka: **payment.completed**

---

# **7️⃣ Fulfillment & Shipping (Stage 7)**

Triggered by shipping worker.

### Steps:

1. Pharmacy fills prescription
2. Shippo used to create shipping label
3. Tracking number generated
4. Delivery events tracked
5. Final confirmation updates status

Outcome → Status: **completed**
Kafka: **shipment.delivered**

---

# ⭐ **Flow Summary Diagram (Text Version)**

```
Stage 1 → Intake  
Stage 2 → Validation  
Stage 3 → Enrollment  
Stage 4 → Pharmacy Routing  
Stage 5 → Insurance + Manufacturer Program Adjudication (Pharmacy)  
Stage 6 → Payment  
Stage 7 → Shipping + Delivery  
```


