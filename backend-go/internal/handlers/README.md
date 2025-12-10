# Handlers Package

This package contains HTTP request handlers for the API endpoints.

## Structure

Handlers should be organized by domain/resource:
- `health.go` - Health check endpoints
- `prescriptions.go` - Prescription-related endpoints
- `patients.go` - Patient-related endpoints
- `pharmacies.go` - Pharmacy-related endpoints
- etc.

Each handler should:
- Accept dependencies via constructor or struct fields
- Return HTTP handlers compatible with the router
- Handle request validation
- Call appropriate services
- Return proper HTTP responses

