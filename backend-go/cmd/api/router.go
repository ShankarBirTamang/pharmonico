// Package main provides router setup
package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pharmonico/backend-gogit/internal/handlers"
)

// setupRouter configures and returns the HTTP router
func (s *Server) setupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Add basic middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check endpoint (outside /api/v1)
	healthHandler := handlers.NewHealthHandler()
	r.Get("/health", healthHandler.GetHealth)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Add v1 routes here
		// Example: r.Get("/prescriptions", prescriptionHandler.List)
	})

	return r
}
