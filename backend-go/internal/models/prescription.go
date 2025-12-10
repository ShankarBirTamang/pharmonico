// Package models provides data models for the application
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrescriptionStatus represents the status of a prescription
type PrescriptionStatus string

const (
	StatusReceived           PrescriptionStatus = "received"
	StatusValidated          PrescriptionStatus = "validated"
	StatusValidationFailed   PrescriptionStatus = "validation_failed"
	StatusAwaitingEnrollment PrescriptionStatus = "awaiting_enrollment"
	StatusAwaitingRouting    PrescriptionStatus = "awaiting_routing"
	StatusRouted             PrescriptionStatus = "routed"
	StatusFulfilled          PrescriptionStatus = "fulfilled"
)

// Prescription represents a prescription document in MongoDB
type Prescription struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PrescriptionID string             `bson:"prescription_id" json:"prescription_id"`
	Status         PrescriptionStatus `bson:"status" json:"status"`

	// Patient information
	Patient PatientInfo `bson:"patient" json:"patient"`

	// Prescriber information
	Prescriber PrescriberInfo `bson:"prescriber" json:"prescriber"`

	// Medication information
	Medication MedicationInfo `bson:"medication" json:"medication"`

	// Insurance information
	Insurance InsuranceInfo `bson:"insurance,omitempty" json:"insurance,omitempty"`

	// Validation errors (if any)
	ValidationErrors []string `bson:"validation_errors,omitempty" json:"validation_errors,omitempty"`

	// Date written (from prescription)
	DateWritten string `bson:"date_written,omitempty" json:"date_written,omitempty"`

	// Metadata
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`

	// Original NCPDP payload (for reference)
	OriginalPayload string `bson:"original_payload,omitempty" json:"original_payload,omitempty"`
}

// PatientInfo contains patient demographic information
type PatientInfo struct {
	ID          string  `bson:"id,omitempty" json:"id,omitempty"`
	FirstName   string  `bson:"first_name" json:"first_name"`
	LastName    string  `bson:"last_name" json:"last_name"`
	DateOfBirth string  `bson:"date_of_birth" json:"date_of_birth"`
	Address     Address `bson:"address,omitempty" json:"address,omitempty"`
	Phone       string  `bson:"phone,omitempty" json:"phone,omitempty"`
}

// PrescriberInfo contains prescriber information
type PrescriberInfo struct {
	ID        string  `bson:"id,omitempty" json:"id,omitempty"`
	NPI       string  `bson:"npi" json:"npi"`
	DEA       string  `bson:"dea,omitempty" json:"dea,omitempty"`
	FirstName string  `bson:"first_name" json:"first_name"`
	LastName  string  `bson:"last_name" json:"last_name"`
	Address   Address `bson:"address,omitempty" json:"address,omitempty"`
	Phone     string  `bson:"phone,omitempty" json:"phone,omitempty"`
}

// MedicationInfo contains medication details
type MedicationInfo struct {
	NDC        string `bson:"ndc" json:"ndc"`
	Name       string `bson:"name" json:"name"`
	Quantity   int    `bson:"quantity" json:"quantity"`
	Refills    int    `bson:"refills,omitempty" json:"refills,omitempty"`
	Dosage     string `bson:"dosage,omitempty" json:"dosage,omitempty"`
	Directions string `bson:"directions,omitempty" json:"directions,omitempty"`
}

// InsuranceInfo contains insurance information
type InsuranceInfo struct {
	BIN      string `bson:"bin,omitempty" json:"bin,omitempty"`
	PCN      string `bson:"pcn,omitempty" json:"pcn,omitempty"`
	GroupID  string `bson:"group_id,omitempty" json:"group_id,omitempty"`
	MemberID string `bson:"member_id,omitempty" json:"member_id,omitempty"`
	PlanName string `bson:"plan_name,omitempty" json:"plan_name,omitempty"`
}

// Address represents a physical address
type Address struct {
	Street  string `bson:"street,omitempty" json:"street,omitempty"`
	City    string `bson:"city,omitempty" json:"city,omitempty"`
	State   string `bson:"state,omitempty" json:"state,omitempty"`
	ZipCode string `bson:"zip_code,omitempty" json:"zip_code,omitempty"`
}

// IntakeRequest represents the request body for prescription intake
type IntakeRequest struct {
	Payload string `json:"payload"` // Can be XML or JSON string
	Format  string `json:"format"`  // "xml" or "json"
}

// IntakeResponse represents the response from prescription intake
type IntakeResponse struct {
	PrescriptionID string `json:"prescription_id"`
	Message        string `json:"message,omitempty"`
}
