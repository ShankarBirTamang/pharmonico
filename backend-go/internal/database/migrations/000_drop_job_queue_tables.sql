-- Migration: Drop job queue tables (if they exist)
-- This ensures we're using Kafka-only workflow, not PostgreSQL job queues

-- Drop all job queue related tables
DROP TABLE IF EXISTS tracking_jobs CASCADE;
DROP TABLE IF EXISTS shipping_jobs CASCADE;
DROP TABLE IF EXISTS payment_jobs CASCADE;
DROP TABLE IF EXISTS adjudication_jobs CASCADE;
DROP TABLE IF EXISTS routing_jobs CASCADE;
DROP TABLE IF EXISTS enrollment_jobs CASCADE;
DROP TABLE IF EXISTS validation_jobs CASCADE;

-- Drop generic job queue tables
DROP TABLE IF EXISTS job_queue CASCADE;
DROP TABLE IF EXISTS dead_letter_queue CASCADE;

