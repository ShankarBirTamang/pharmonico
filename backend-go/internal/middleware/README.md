# Middleware Package

This package contains HTTP middleware for the API server.

## Structure

Middleware functions should be organized by concern:
- `logging.go` - Request logging middleware
- `recovery.go` - Panic recovery middleware
- `cors.go` - CORS configuration middleware
- `correlation.go` - Correlation ID middleware
- `auth.go` - Authentication/authorization middleware

Each middleware should:
- Accept and return `http.Handler` or router-compatible middleware
- Be composable with other middleware
- Handle errors gracefully
- Log appropriately

