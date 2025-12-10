// Package handlers provides HTTP request handlers
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pharmonico/backend-gogit/internal/models"
	"github.com/pharmonico/backend-gogit/pkg/ncpdp"
)

// PrescriptionHandler handles prescription-related requests
type PrescriptionHandler struct {
	deps *Dependencies
}

// NewPrescriptionHandler creates a new prescription handler
func NewPrescriptionHandler(deps *Dependencies) *PrescriptionHandler {
	return &PrescriptionHandler{
		deps: deps,
	}
}

// Intake handles POST /api/v1/prescriptions/intake
// Subtask 1.1.1: Create route POST /api/v1/prescriptions/intake
// Subtask 1.1.2: Parse request body (XML or JSON)
func (h *PrescriptionHandler) Intake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req models.IntakeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Payload == "" {
		http.Error(w, "Payload is required", http.StatusBadRequest)
		return
	}

	// Determine format (default to XML if not specified)
	format := strings.ToLower(req.Format)
	if format == "" {
		format = "xml"
	}

	// Parse based on format
	var prescription *models.Prescription
	var parseErr error

	switch format {
	case "xml":
		// Subtask 1.1.3: Implement NCPDP parser
		// Subtask 1.1.4: Extract patient, prescriber, medication, insurance
		// Subtask 1.1.5: Validate XML structure (well-formedness check)
		prescription, parseErr = ncpdp.ParseXML(req.Payload)
		if parseErr != nil {
			log.Printf("Error parsing XML: %v", parseErr)
			http.Error(w, fmt.Sprintf("Failed to parse XML: %v", parseErr), http.StatusBadRequest)
			return
		}
	case "json":
		prescription, parseErr = ncpdp.ParseJSON(req.Payload)
		if parseErr != nil {
			log.Printf("Error parsing JSON: %v", parseErr)
			http.Error(w, fmt.Sprintf("Failed to parse JSON: %v", parseErr), http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Invalid format. Supported formats: xml, json", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if err := validateRequiredFields(prescription); err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, fmt.Sprintf("Missing required fields: %v", err), http.StatusBadRequest)
		return
	}

	// For now, just return success with a mock prescription_id
	// We'll add persistence, deduplication, and Kafka in the next steps
	response := models.IntakeResponse{
		PrescriptionID: "rx_" + prescription.Patient.ID + "_" + prescription.Medication.NDC,
		Message:        "Prescription received successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// validateRequiredFields checks that all required fields are present
func validateRequiredFields(p *models.Prescription) error {
	var errors []string

	// Patient required fields
	if p.Patient.FirstName == "" {
		errors = append(errors, "patient.first_name is required")
	}
	if p.Patient.LastName == "" {
		errors = append(errors, "patient.last_name is required")
	}
	if p.Patient.DateOfBirth == "" {
		errors = append(errors, "patient.date_of_birth is required")
	}

	// Prescriber required fields
	if p.Prescriber.NPI == "" {
		errors = append(errors, "prescriber.npi is required")
	}
	if p.Prescriber.FirstName == "" {
		errors = append(errors, "prescriber.first_name is required")
	}
	if p.Prescriber.LastName == "" {
		errors = append(errors, "prescriber.last_name is required")
	}

	// Medication required fields
	if p.Medication.NDC == "" {
		errors = append(errors, "medication.ndc is required")
	}
	if p.Medication.Name == "" {
		errors = append(errors, "medication.name is required")
	}
	if p.Medication.Quantity == 0 {
		errors = append(errors, "medication.quantity is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}

	return nil
}
