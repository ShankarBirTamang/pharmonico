# Worker Service - Kafka-Based Event Processing

## Overview

The worker service processes Kafka events in an event-driven architecture. All processing is done through Kafka consumers, **NOT** PostgreSQL job queues.

## Architecture

### Task 8.4: Removed PostgreSQL Job Queue Dependencies

- ✅ **PostgreSQL polling loops removed** - All workers now consume from Kafka topics
- ✅ **Job queue tables removed** - See migration `000_drop_job_queue_tables.sql`
- ✅ **Kafka-only workflow** - All event processing is event-driven via Kafka

### Event Flow

1. **Consume** - Worker polls Kafka for events (8.3.1)
2. **Process** - Handler processes business logic (8.3.2)
3. **Emit** - Handler publishes next event to Kafka (8.3.3)

## Worker Handlers

All handlers implement the `Handler` interface and process events from their respective Kafka topics:

- **ValidationWorker** - `prescription.intake.received` → `prescription.validation.completed`
- **EnrollmentWorker** - `prescription.validation.completed` → `patient.enrollment.completed`
- **RoutingWorker** - `patient.enrollment.completed` → `pharmacy.selected`
- **AdjudicationWorker** - `pharmacy.selected` → `insurance.adjudication.completed`
- **PaymentWorker** - `insurance.adjudication.completed` → `payment.link.created` / `payment.completed`
- **ShippingWorker** - `payment.completed` → `shipment.label.created`
- **DeliveryWorker** - `shipment.label.created` → (tracks delivery)

## PostgreSQL Usage

PostgreSQL is **only** used for:
- ✅ Audit logs (HIPAA compliance)
- ❌ **NOT** for job queues (removed in Task 8.4)

## Migration

The migration `000_drop_job_queue_tables.sql` ensures all job queue tables are removed:
- `validation_jobs`
- `enrollment_jobs`
- `routing_jobs`
- `adjudication_jobs`
- `payment_jobs`
- `shipping_jobs`
- `tracking_jobs`
- `job_queue`
- `dead_letter_queue` (PostgreSQL version - Kafka DLQ is used instead)

## Error Handling

- Failed messages are sent to Kafka dead letter queue (`dead_letter_queue` topic)
- Unhandled topics are sent to DLQ
- Correlation IDs are propagated through all events for traceability

