// Package ncpdp provides NCPDP SCRIPT format parsing functionality
package ncpdp

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/pharmonico/backend-gogit/internal/models"
)

// NCPDPScript represents the root element of an NCPDP SCRIPT message
type NCPDPScript struct {
	XMLName xml.Name `xml:"Message"`
	Header  Header   `xml:"Header"`
	Body    Body     `xml:"Body"`
}

// Header contains message header information
type Header struct {
	MessageID string `xml:"MessageID"`
	RelatesTo string `xml:"RelatesTo,omitempty"`
	Timestamp string `xml:"Timestamp"`
}

// Body contains the prescription data
type Body struct {
	Prescription PrescriptionXML `xml:"Prescription"`
}

// PrescriptionXML represents prescription data in NCPDP format
type PrescriptionXML struct {
	Patient     PatientXML    `xml:"Patient"`
	Prescriber  PrescriberXML `xml:"Prescriber"`
	Medication  MedicationXML `xml:"Medication"`
	Insurance   InsuranceXML  `xml:"Insurance,omitempty"`
	DateWritten string        `xml:"DateWritten"`
}

// PatientXML represents patient information in NCPDP format
type PatientXML struct {
	ID          string     `xml:"ID,attr"`
	FirstName   string     `xml:"FirstName"`
	LastName    string     `xml:"LastName"`
	DateOfBirth string     `xml:"DateOfBirth"`
	Address     AddressXML `xml:"Address,omitempty"`
	Phone       string     `xml:"Phone,omitempty"`
}

// PrescriberXML represents prescriber information in NCPDP format
type PrescriberXML struct {
	ID        string     `xml:"ID,attr"`
	NPI       string     `xml:"NPI"`
	DEA       string     `xml:"DEA,omitempty"`
	FirstName string     `xml:"FirstName"`
	LastName  string     `xml:"LastName"`
	Address   AddressXML `xml:"Address,omitempty"`
	Phone     string     `xml:"Phone,omitempty"`
}

// MedicationXML represents medication information in NCPDP format
type MedicationXML struct {
	NDC        string `xml:"NDC"`
	Name       string `xml:"Name"`
	Quantity   int    `xml:"Quantity"`
	Refills    int    `xml:"Refills,omitempty"`
	Dosage     string `xml:"Dosage,omitempty"`
	Directions string `xml:"Directions,omitempty"`
}

// InsuranceXML represents insurance information in NCPDP format
type InsuranceXML struct {
	BIN      string `xml:"BIN,omitempty"`
	PCN      string `xml:"PCN,omitempty"`
	GroupID  string `xml:"GroupID,omitempty"`
	MemberID string `xml:"MemberID,omitempty"`
	PlanName string `xml:"PlanName,omitempty"`
}

// AddressXML represents address information in NCPDP format
type AddressXML struct {
	Street  string `xml:"Street,omitempty"`
	City    string `xml:"City,omitempty"`
	State   string `xml:"State,omitempty"`
	ZipCode string `xml:"ZipCode,omitempty"`
}

// ParseXML parses an NCPDP SCRIPT XML string and returns a Prescription model
func ParseXML(xmlData string) (*models.Prescription, error) {
	// Validate XML is well-formed
	if err := validateXMLWellFormed(xmlData); err != nil {
		return nil, fmt.Errorf("invalid XML structure: %w", err)
	}

	var script NCPDPScript
	if err := xml.Unmarshal([]byte(xmlData), &script); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	// Extract and convert to Prescription model
	prescription := &models.Prescription{
		Status:          models.StatusReceived,
		DateWritten:     script.Body.Prescription.DateWritten,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		OriginalPayload: xmlData,
	}

	// Extract patient information
	prescription.Patient = models.PatientInfo{
		ID:          script.Body.Prescription.Patient.ID,
		FirstName:   script.Body.Prescription.Patient.FirstName,
		LastName:    script.Body.Prescription.Patient.LastName,
		DateOfBirth: script.Body.Prescription.Patient.DateOfBirth,
		Phone:       script.Body.Prescription.Patient.Phone,
	}
	if script.Body.Prescription.Patient.Address.Street != "" {
		prescription.Patient.Address = models.Address{
			Street:  script.Body.Prescription.Patient.Address.Street,
			City:    script.Body.Prescription.Patient.Address.City,
			State:   script.Body.Prescription.Patient.Address.State,
			ZipCode: script.Body.Prescription.Patient.Address.ZipCode,
		}
	}

	// Extract prescriber information
	prescription.Prescriber = models.PrescriberInfo{
		ID:        script.Body.Prescription.Prescriber.ID,
		NPI:       script.Body.Prescription.Prescriber.NPI,
		DEA:       script.Body.Prescription.Prescriber.DEA,
		FirstName: script.Body.Prescription.Prescriber.FirstName,
		LastName:  script.Body.Prescription.Prescriber.LastName,
		Phone:     script.Body.Prescription.Prescriber.Phone,
	}
	if script.Body.Prescription.Prescriber.Address.Street != "" {
		prescription.Prescriber.Address = models.Address{
			Street:  script.Body.Prescription.Prescriber.Address.Street,
			City:    script.Body.Prescription.Prescriber.Address.City,
			State:   script.Body.Prescription.Prescriber.Address.State,
			ZipCode: script.Body.Prescription.Prescriber.Address.ZipCode,
		}
	}

	// Extract medication information
	prescription.Medication = models.MedicationInfo{
		NDC:        script.Body.Prescription.Medication.NDC,
		Name:       script.Body.Prescription.Medication.Name,
		Quantity:   script.Body.Prescription.Medication.Quantity,
		Refills:    script.Body.Prescription.Medication.Refills,
		Dosage:     script.Body.Prescription.Medication.Dosage,
		Directions: script.Body.Prescription.Medication.Directions,
	}

	// Extract insurance information (optional)
	if script.Body.Prescription.Insurance.BIN != "" || script.Body.Prescription.Insurance.MemberID != "" {
		prescription.Insurance = models.InsuranceInfo{
			BIN:      script.Body.Prescription.Insurance.BIN,
			PCN:      script.Body.Prescription.Insurance.PCN,
			GroupID:  script.Body.Prescription.Insurance.GroupID,
			MemberID: script.Body.Prescription.Insurance.MemberID,
			PlanName: script.Body.Prescription.Insurance.PlanName,
		}
	}

	return prescription, nil
}

// ParseJSON parses an NCPDP SCRIPT JSON string and returns a Prescription model
// Note: This is a simplified implementation. In production, you'd have a proper JSON schema
func ParseJSON(jsonData string) (*models.Prescription, error) {
	// For now, we'll support a simplified JSON format
	// In production, you'd want to use a proper JSON schema validator
	return nil, fmt.Errorf("JSON parsing not yet implemented - please use XML format")
}

// validateXMLWellFormed checks if the XML string is well-formed
// Subtask 1.1.5: Validate XML structure (well-formedness check)
func validateXMLWellFormed(xmlData string) error {
	if strings.TrimSpace(xmlData) == "" {
		return fmt.Errorf("XML data is empty")
	}

	// Use xml.NewDecoder to validate well-formedness
	decoder := xml.NewDecoder(strings.NewReader(xmlData))

	// Try to decode the entire document to check for well-formedness
	for {
		_, err := decoder.Token()
		if err != nil {
			// EOF means we successfully parsed the entire document
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("XML is not well-formed: %w", err)
		}
	}
}

// GenerateDedupHash generates a deduplication hash from core prescription fields
// Subtask 1.1.6: Generate dedup hash from core fields (patient + drug + date)
// Uses SHA256 hash of: patient_id + drug_ndc + date_written
func GenerateDedupHash(patientID, drugNDC, dateWritten string) string {
	// Create a composite key from core fields
	compositeKey := fmt.Sprintf("%s:%s:%s", patientID, drugNDC, dateWritten)
	
	// Generate SHA256 hash for consistent, fixed-length keys
	hash := sha256.Sum256([]byte(compositeKey))
	hashString := hex.EncodeToString(hash[:])
	
	// Return Redis key with prefix
	return fmt.Sprintf("rx:dedup:%s", hashString)
}
