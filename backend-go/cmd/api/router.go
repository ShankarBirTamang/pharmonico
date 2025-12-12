// Package main provides router setup
package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/phil-my-meds/backend-gogit/internal/handlers"
	appMiddleware "github.com/phil-my-meds/backend-gogit/internal/middleware"
)

// setupRouter configures and returns the HTTP router
func (s *Server) setupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Add middleware in order (first added is outermost)
	// 1. CORS - must be first to handle preflight requests
	r.Use(appMiddleware.CORSMiddleware)

	// 2. Correlation ID - extract or generate correlation ID
	r.Use(appMiddleware.CorrelationIDMiddleware)

	// 3. Real IP - get client's real IP address
	r.Use(middleware.RealIP)

	// 4. Logging - log all requests
	r.Use(appMiddleware.LoggingMiddleware)

	// 5. Panic recovery - recover from panics
	r.Use(appMiddleware.RecoveryMiddleware)

	// Create handler dependencies with all injected services
	deps := handlers.NewDependencies(
		s.MongoClient,
		s.Postgres,
		s.Redis,
		s.KafkaProducer,
	)

	// Health check endpoint (outside /api/v1)
	healthHandler := handlers.NewHealthHandler(deps)
	r.Get("/health", healthHandler.GetHealth)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Prescription routes
		prescriptionHandler := handlers.NewPrescriptionHandler(deps)
		r.Post("/prescriptions/intake", prescriptionHandler.Intake)
	})

	return r
}
