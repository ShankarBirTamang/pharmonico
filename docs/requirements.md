# Pharmonico - Prescription Fulfillment System Requirements Document

## 1. Project Overview

### 1.1 Purpose
A learning project to understand the US healthcare system through building a prescription fulfillment service that handles the complete workflow from prescription intake to delivery. This is an educational project, not intended for real-world production use.

### 1.2 Business Problem
The system simulates a specialty pharmacy fulfillment service that acts as an intermediary between healthcare providers (HCPs), patients, pharmacies, and insurance companies. It streamlines the complex process of prescription fulfillment by managing enrollment, insurance verification, payment collection, and delivery coordination.

### 1.3 Key Learning Objectives
- Understanding US healthcare workflow, terminology, and NCPDP standards
- Asynchronous job processing patterns using database polling
- Full-stack development with modern technologies
- Third-party API integrations in healthcare context
- HIPAA compliance concepts and patient data handling

---

## 2. Business Workflow

### 2.1 High-Level Process Flow

**Stage 1: Prescription Intake**
- Healthcare providers submit prescriptions to the system
- For this learning project, Gemini API will generate mock prescriptions in NCPDP SCRIPT format
- System receives prescription data including patient, prescriber, medication, and insurance details

**Stage 2: Prescription Validation**
- Automated validation worker verifies prescription format and required fields
- Checks for data completeness and basic compliance rules
- Validates prescriber credentials (NPI, DEA numbers)
- Ensures medication information is complete (NDC code, quantity, directions)

**Stage 3: Patient Enrollment**
- System checks if patient is already enrolled in the system
- If not enrolled, generates a secure magic link sent to patient
- Patient accesses enrollment portal through the link where they:
  - Enter or verify insurance information (carrier, member ID, group number, BIN, PCN)
  - Optionally upload insurance card images
  - Accept and enroll in manufacturer coupon programs
  - Provide HIPAA consent and electronic signature
- Enrollment worker monitors completion status

**Stage 4: Pharmacy Routing & Selection**
- System automatically filters and ranks partner pharmacies based on:
  - Patient geographic location and shipping address
  - Pharmacy capacity and current workload
  - Insurance network contracts (whether pharmacy accepts patient's insurance)
  - Specialty medication capabilities if needed
- Operations team reviews filtered pharmacy options presented in dashboard
- Ops team manually selects the most appropriate pharmacy from recommendations

**Stage 5: Insurance Adjudication**
- Selected pharmacy receives the prescription details
- Pharmacy runs insurance claim through their system (simulated via API in this project)
- Insurance adjudication determines:
  - Insurance reimbursement amount
  - Manufacturer coupon discount application
  - Final patient copay amount
- Results are recorded and communicated back to the system

**Stage 6: Payment Collection**
- System generates Stripe payment link for the patient copay amount
- Patient receives communication (email/SMS) with:
  - Breakdown of costs (drug cost, insurance coverage, coupon discount, copay)
  - Payment link and instructions
  - Estimated delivery timeline
- Payment worker monitors payment completion status
- Upon successful payment, prescription moves to fulfillment

**Stage 7: Fulfillment & Shipping**
- Pharmacy prepares and fills the prescription
- System integrates with Shippo to generate shipping label
- Tracking number is captured and stored
- Patient receives shipping notification with tracking information
- Shipping worker monitors delivery status updates
- Upon delivery confirmation, prescription is marked complete

### 2.2 Operations Dashboard

The operations team uses a centralized dashboard with tab-based views representing each workflow stage:

- **Intake Tab**: Newly received prescriptions awaiting validation
- **Validation Tab**: Prescriptions being validated or with validation issues
- **Enrollment Tab**: Prescriptions awaiting patient enrollment
- **Pharmacy Routing Tab**: Ready for pharmacy selection with filtered recommendations
- **Insurance Processing Tab**: Awaiting or processing insurance adjudication
- **Payment Pending Tab**: Awaiting patient payment
- **Fulfillment Tab**: Paid prescriptions being filled and shipped
- **Completed Tab**: Delivered prescriptions

Each prescription card displays on the appropriate tab based on current status. Ops team can view details, take actions (like selecting pharmacy), and track progress.

---

## 3. Technology Stack

### 3.1 Backend
- **Language**: Go (Golang)
- **Framework**: Standard net/http library
- **Databases**: 
  - MongoDB for business data (prescriptions, patients, pharmacies, enrollments, payments, shipments)
  - PostgreSQL for job queue and audit logs
- **Authentication**: JWT tokens for API authentication

### 3.2 Frontend
- **Framework**: React
- **Architecture**: Single-page application serving both operations dashboard and patient enrollment portal
- **State Management**: Context API or Zustand
- **UI Library**: Material-UI or Tailwind CSS
- **Routing**: Different routes for ops dashboard vs patient portal

### 3.3 Job Processing System
- **Pattern**: Database polling (PostgreSQL-based job queue)
- **Workers**: Continuous goroutines polling for jobs
  - **Validation Worker**: Processes prescription validation
  - **Enrollment Worker**: Checks enrollment completion status
  - **Payment Worker**: Monitors payment status
  - **Shipping Worker**: Tracks shipment updates
- **Job Management**: Retry logic, error handling, status tracking

### 3.4 Third-Party Integrations

**Gemini API**
- Purpose: Generate mock NCPDP-formatted prescriptions
- Usage: Simulate HCP prescription submissions

**Insurance Adjudication API (Custom Built)**
- Purpose: Simulate pharmacy insurance claim processing
- Returns: Mock adjudication results with cost breakdown

**Stripe (Test Mode)**
- Purpose: Payment processing
- Features: Payment link generation, payment intent tracking, webhook handling

**Shippo**
- Purpose: Shipping label generation and tracking
- Features: Label creation, carrier selection, tracking updates

**SendGrid / Twilio**
- Purpose: Patient communications
- Features: Email and SMS notifications for enrollment links, payment requests, shipping updates

### 3.5 Development Environment
- **Containerization**: Docker and Docker Compose
- **Local Development**: All services running in containers
- **Database Setup**: MongoDB and PostgreSQL containers with seed data

---

## 4. User Roles & Access

### 4.1 Operations Team
- Login with JWT authentication
- Access to full operations dashboard
- Capabilities:
  - View all prescriptions across all stages
  - Select pharmacies from filtered recommendations
  - Trigger insurance adjudication manually
  - View patient and prescription details
  - Search and filter prescriptions
  - Access audit logs

### 4.2 Patients
- Access via magic link (no traditional login for enrollment)
- Optional JWT-based login for order tracking (future enhancement)
- Capabilities:
  - Complete enrollment form
  - Upload insurance information
  - Accept manufacturer coupon terms
  - Provide HIPAA consent
  - Make payment via Stripe
  - View order status and tracking

### 4.3 System/Automated Processes
- Background workers processing jobs
- Automated notifications
- Status updates based on external webhooks (Stripe, Shippo)

---

## 5. Key Features

### 5.1 Prescription Management
- NCPDP SCRIPT format support
- Validation engine with configurable rules
- Status tracking through complete lifecycle
- Audit trail for all changes

### 5.2 Patient Enrollment
- Secure magic link generation with expiration
- Insurance information capture
- Manufacturer coupon enrollment
- HIPAA consent with electronic signature
- Optional insurance card image upload

### 5.3 Intelligent Pharmacy Routing
- Multi-criteria filtering (location, capacity, insurance contracts)
- Weighted scoring and ranking algorithm
- Manual selection from top recommendations
- Pharmacy capacity tracking

### 5.4 Payment Processing
- Stripe integration for secure payments
- Payment link generation
- Cost breakdown transparency
- Payment status tracking
- Webhook handling for payment confirmation

### 5.5 Shipping & Tracking
- Shippo integration for label generation
- Multi-carrier support
- Real-time tracking updates
- Delivery confirmation
- Customer notifications at each milestone

### 5.6 Communication System
- Multi-channel notifications (Email & SMS)
- Event-driven messaging:
  - Enrollment invitation
  - Payment request
  - Shipping confirmation
  - Delivery notification
- Template-based messaging

---

## 6. Success Criteria

- Complete end-to-end prescription flow from intake to delivery
- All worker processes functioning with proper error handling
- Operations dashboard displaying prescriptions in correct stage tabs
- Patient enrollment flow with insurance and consent capture
- Successful Stripe payment integration
- Successful Shippo shipping integration
- Proper audit logging of all actions
- Understanding of US healthcare terminology and workflows