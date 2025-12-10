import { useState } from 'react'
import { Card, CardHeader, CardTitle, CardContent } from '../components/ui'
import { PrescriptionForm, PrescriptionFormData } from '../components/forms'
import { api } from '../services/api'
import { getPrescriptionExample } from '../services/examples'
import { convertFormDataToXML, convertExampleToFormData } from '../utils/xmlConverter'

const initialFormData: PrescriptionFormData = {
  // Patient
  patientId: '',
  patientFirstName: '',
  patientLastName: '',
  patientDateOfBirth: '',
  patientStreet: '',
  patientCity: '',
  patientState: '',
  patientZipCode: '',
  patientPhone: '',

  // Prescriber
  prescriberId: '',
  prescriberNPI: '',
  prescriberDEA: '',
  prescriberFirstName: '',
  prescriberLastName: '',
  prescriberStreet: '',
  prescriberCity: '',
  prescriberState: '',
  prescriberZipCode: '',
  prescriberPhone: '',

  // Medication
  medicationNDC: '',
  medicationName: '',
  medicationQuantity: '',
  medicationRefills: '',
  medicationDosage: '',
  medicationDirections: '',
  dateWritten: new Date().toISOString().split('T')[0],

  // Insurance
  insuranceBIN: '',
  insurancePCN: '',
  insuranceGroupID: '',
  insuranceMemberID: '',
  insurancePlanName: '',
}

export function Intake() {
  const [formData, setFormData] = useState<PrescriptionFormData>(initialFormData)
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<{ prescription_id?: string; error?: string } | null>(null)
  const [errors, setErrors] = useState<Partial<Record<keyof PrescriptionFormData, string>>>({})

  const validateForm = (): boolean => {
    const newErrors: Partial<Record<keyof PrescriptionFormData, string>> = {}

    // Patient required fields
    if (!formData.patientFirstName.trim()) {
      newErrors.patientFirstName = 'First name is required'
    }
    if (!formData.patientLastName.trim()) {
      newErrors.patientLastName = 'Last name is required'
    }
    if (!formData.patientDateOfBirth) {
      newErrors.patientDateOfBirth = 'Date of birth is required'
    }

    // Prescriber required fields
    if (!formData.prescriberNPI.trim()) {
      newErrors.prescriberNPI = 'NPI is required'
    } else if (formData.prescriberNPI.length !== 10 || !/^\d+$/.test(formData.prescriberNPI)) {
      newErrors.prescriberNPI = 'NPI must be 10 digits'
    }
    if (!formData.prescriberFirstName.trim()) {
      newErrors.prescriberFirstName = 'First name is required'
    }
    if (!formData.prescriberLastName.trim()) {
      newErrors.prescriberLastName = 'Last name is required'
    }

    // Medication required fields
    if (!formData.medicationNDC.trim()) {
      newErrors.medicationNDC = 'NDC is required'
    }
    if (!formData.medicationName.trim()) {
      newErrors.medicationName = 'Medication name is required'
    }
    if (!formData.medicationQuantity || parseInt(formData.medicationQuantity) <= 0) {
      newErrors.medicationQuantity = 'Quantity must be greater than 0'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setResult(null)

    if (!validateForm()) {
      return
    }

    setLoading(true)

    try {
      // Convert form data to XML
      const xmlPayload = convertFormDataToXML(formData)
      
      // Submit to API
      const response = await api.submitPrescription(xmlPayload, 'xml')
      setResult({ prescription_id: response.prescription_id })
      
      // Clear form on success
      setFormData(initialFormData)
      setErrors({})
    } catch (error) {
      setResult({ error: error instanceof Error ? error.message : 'Failed to submit prescription' })
    } finally {
      setLoading(false)
    }
  }

  const handleLoadExample = () => {
    setResult(null)
    setErrors({})

    // Get a random static example
    const example = getPrescriptionExample()
    
    // Convert example to form data
    const exampleFormData = convertExampleToFormData(example)
    setFormData(exampleFormData)
  }

  return (
    <div className="max-w-6xl mx-auto">
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Prescription Intake Form</CardTitle>
              <p className="text-sm text-slate-600 dark:text-slate-400 mt-2">
                Fill out the prescription form below as a Healthcare Provider. All fields marked with <span className="text-red-500">*</span> are required.
              </p>
            </div>
            <button
              type="button"
              onClick={handleLoadExample}
              className="btn-secondary flex items-center gap-2"
            >
              <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Load Example
            </button>
          </div>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <PrescriptionForm
              formData={formData}
              onChange={setFormData}
              errors={errors}
            />

            <div className="flex gap-4 pt-4 border-t border-slate-200 dark:border-slate-700">
              <button
                type="submit"
                disabled={loading}
                className="btn-primary flex-1 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? (
                  <span className="flex items-center justify-center gap-2">
                    <span className="inline-block h-4 w-4 animate-spin rounded-full border-2 border-solid border-white border-r-transparent"></span>
                    Submitting Prescription...
                  </span>
                ) : (
                  'Submit Prescription'
                )}
              </button>
              <button
                type="button"
                onClick={() => {
                  setFormData(initialFormData)
                  setErrors({})
                  setResult(null)
                }}
                className="btn-secondary"
                disabled={loading}
              >
                Clear Form
              </button>
            </div>
          </form>

          {result && (
            <div
              className={`mt-6 p-4 rounded-lg ${
                result.prescription_id
                  ? 'bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800'
                  : 'bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800'
              }`}
            >
              {result.prescription_id ? (
                <div>
                  <p className="text-emerald-800 dark:text-emerald-400 font-medium">✓ Success!</p>
                  <p className="text-sm text-emerald-700 dark:text-emerald-300 mt-1">
                    Prescription ID: <code className="font-mono bg-emerald-100 dark:bg-emerald-900/50 px-2 py-1 rounded">{result.prescription_id}</code>
                  </p>
                  <p className="text-xs text-emerald-600 dark:text-emerald-400 mt-2">
                    The prescription has been received and is being processed.
                  </p>
                </div>
              ) : (
                <div>
                  <p className="text-red-800 dark:text-red-400 font-medium">✗ Error</p>
                  <p className="text-sm text-red-700 dark:text-red-300 mt-1">{result.error}</p>
                </div>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
