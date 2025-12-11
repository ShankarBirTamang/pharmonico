// Package handlers provides HTTP request handlers
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pharmonico/backend-gogit/internal/models"
	"github.com/pharmonico/backend-gogit/pkg/ncpdp"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// Subtask 1.1.6: Generate dedup hash from core fields (patient + drug + date)
	// Generate deduplication hash
	patientID := prescription.Patient.ID
	if patientID == "" {
		// Use a composite key if patient ID is not available
		patientID = fmt.Sprintf("%s_%s", prescription.Patient.FirstName, prescription.Patient.LastName)
	}
	dateWritten := prescription.DateWritten
	if dateWritten == "" {
		// Use current date if date written is not available
		dateWritten = time.Now().Format("2006-01-02")
	}
	dedupHash := ncpdp.GenerateDedupHash(patientID, prescription.Medication.NDC, dateWritten)

	// Subtask 1.1.7: Check Redis for duplicates (TTL 5 mins)
	ctx := r.Context()
	exists, err := h.deps.Redis.Exists(ctx, dedupHash)
	if err != nil {
		log.Printf("Error checking Redis for duplicate: %v", err)
		// Continue processing if Redis check fails (don't block intake)
	} else if exists {
		// Subtask 1.1.7: Duplicate detected
		log.Printf("Duplicate prescription detected: %s", dedupHash)
		response := models.IntakeResponse{
			PrescriptionID: "",
			Message:        "Duplicate prescription detected. This prescription was recently submitted.",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		json.NewEncoder(w).Encode(response)
		return
	}

	// Subtask 1.1.8: Store dedup key in Redis if new (TTL 5 mins)
	// Use SetNX for atomic operation - only set if key doesn't exist
	ttl := 5 * time.Minute
	wasSet, err := h.deps.Redis.SetNX(ctx, dedupHash, "1", ttl)
	if err != nil {
		log.Printf("Error storing dedup key in Redis: %v", err)
		// Continue processing even if Redis storage fails
	} else if !wasSet {
		// Key was already set (race condition - another request got there first)
		log.Printf("Race condition: duplicate prescription detected during SetNX: %s", dedupHash)
		response := models.IntakeResponse{
			PrescriptionID: "",
			Message:        "Duplicate prescription detected. This prescription was recently submitted.",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		json.NewEncoder(w).Encode(response)
		return
	}

	// Subtask 1.1.9: Insert prescription into MongoDB
	// Subtask 1.1.10: Set status: "received"
	now := time.Now()
	prescription.Status = models.StatusReceived
	prescription.CreatedAt = now
	prescription.UpdatedAt = now
	prescription.OriginalPayload = req.Payload

	// Generate prescription_id if not already set
	if prescription.PrescriptionID == "" {
		// Generate a unique prescription ID using timestamp and patient/medication info
		prescription.PrescriptionID = fmt.Sprintf("rx_%d_%s_%s", now.Unix(), prescription.Patient.ID, prescription.Medication.NDC)
	}

	// Insert into MongoDB
	collection := h.deps.MongoClient.GetCollection("prescriptions")
	result, err := collection.InsertOne(ctx, prescription)
	if err != nil {
		log.Printf("Error inserting prescription into MongoDB: %v", err)
		http.Error(w, "Failed to save prescription", http.StatusInternalServerError)
		return
	}

	// Get the inserted ID
	insertedID := result.InsertedID
	var prescriptionID string
	if oid, ok := insertedID.(primitive.ObjectID); ok {
		prescriptionID = oid.Hex()
	} else {
		prescriptionID = prescription.PrescriptionID
	}

	log.Printf("Prescription inserted successfully with ID: %s", prescriptionID)

	// Subtask 1.1.13: Return { prescription_id }
	response := models.IntakeResponse{
		PrescriptionID: prescriptionID,
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
