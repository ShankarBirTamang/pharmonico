// Package middleware provides HTTP middleware functions
package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// CorrelationIDKey is the context key for correlation ID
type CorrelationIDKey struct{}

// CorrelationIDMiddleware extracts or generates a correlation ID for each request
func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get correlation ID from header
		correlationID := r.Header.Get("X-Correlation-ID")

		// Generate new correlation ID if not present
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Add correlation ID to response header
		w.Header().Set("X-Correlation-ID", correlationID)

		// Add correlation ID to request context
		ctx := context.WithValue(r.Context(), CorrelationIDKey{}, correlationID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetCorrelationID extracts the correlation ID from the request context
func GetCorrelationID(r *http.Request) string {
	if id, ok := r.Context().Value(CorrelationIDKey{}).(string); ok {
		return id
	}
	return ""
}
