import { useState } from 'react'

export interface PrescriptionFormData {
  // Patient Information
  patientId: string
  patientFirstName: string
  patientLastName: string
  patientDateOfBirth: string
  patientStreet: string
  patientCity: string
  patientState: string
  patientZipCode: string
  patientPhone: string

  // Prescriber Information
  prescriberId: string
  prescriberNPI: string
  prescriberDEA: string
  prescriberFirstName: string
  prescriberLastName: string
  prescriberStreet: string
  prescriberCity: string
  prescriberState: string
  prescriberZipCode: string
  prescriberPhone: string

  // Medication Information
  medicationNDC: string
  medicationName: string
  medicationQuantity: string
  medicationRefills: string
  medicationDosage: string
  medicationDirections: string
  dateWritten: string

  // Insurance Information
  insuranceBIN: string
  insurancePCN: string
  insuranceGroupID: string
  insuranceMemberID: string
  insurancePlanName: string
}

interface PrescriptionFormProps {
  formData: PrescriptionFormData
  onChange: (data: PrescriptionFormData) => void
  errors?: Partial<Record<keyof PrescriptionFormData, string>>
}

export function PrescriptionForm({ formData, onChange, errors = {} }: PrescriptionFormProps) {
  const handleChange = (field: keyof PrescriptionFormData, value: string) => {
    onChange({ ...formData, [field]: value })
  }

  return (
    <div className="space-y-8">
      {/* Patient Information Section */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-slate-900 dark:text-white border-b border-slate-200 dark:border-slate-700 pb-2">
          Patient Information
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              Patient ID <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              value={formData.patientId}
              onChange={(e) => handleChange('patientId', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="PAT001"
            />
            {errors.patientId && <p className="text-xs text-red-600 mt-1">{errors.patientId}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              First Name <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              value={formData.patientFirstName}
              onChange={(e) => handleChange('patientFirstName', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="John"
            />
            {errors.patientFirstName && <p className="text-xs text-red-600 mt-1">{errors.patientFirstName}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              Last Name <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              value={formData.patientLastName}
              onChange={(e) => handleChange('patientLastName', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="Doe"
            />
            {errors.patientLastName && <p className="text-xs text-red-600 mt-1">{errors.patientLastName}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              Date of Birth <span className="text-red-500">*</span>
            </label>
            <input
              type="date"
              value={formData.patientDateOfBirth}
              onChange={(e) => handleChange('patientDateOfBirth', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
            />
            {errors.patientDateOfBirth && <p className="text-xs text-red-600 mt-1">{errors.patientDateOfBirth}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Street Address</label>
            <input
              type="text"
              value={formData.patientStreet}
              onChange={(e) => handleChange('patientStreet', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="123 Main St"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">City</label>
            <input
              type="text"
              value={formData.patientCity}
              onChange={(e) => handleChange('patientCity', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="New York"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">State</label>
            <input
              type="text"
              value={formData.patientState}
              onChange={(e) => handleChange('patientState', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="NY"
              maxLength={2}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Zip Code</label>
            <input
              type="text"
              value={formData.patientZipCode}
              onChange={(e) => handleChange('patientZipCode', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="10001"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Phone</label>
            <input
              type="tel"
              value={formData.patientPhone}
              onChange={(e) => handleChange('patientPhone', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="555-1234"
            />
          </div>
        </div>
      </div>

      {/* Prescriber Information Section */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-slate-900 dark:text-white border-b border-slate-200 dark:border-slate-700 pb-2">
          Prescriber Information
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              Prescriber ID
            </label>
            <input
              type="text"
              value={formData.prescriberId}
              onChange={(e) => handleChange('prescriberId', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="PRES001"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              NPI <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              value={formData.prescriberNPI}
              onChange={(e) => handleChange('prescriberNPI', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="1234567890"
              maxLength={10}
            />
            {errors.prescriberNPI && <p className="text-xs text-red-600 mt-1">{errors.prescriberNPI}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">DEA Number</label>
            <input
              type="text"
              value={formData.prescriberDEA}
              onChange={(e) => handleChange('prescriberDEA', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="AB1234567"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              First Name <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              value={formData.prescriberFirstName}
              onChange={(e) => handleChange('prescriberFirstName', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="Jane"
            />
            {errors.prescriberFirstName && <p className="text-xs text-red-600 mt-1">{errors.prescriberFirstName}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              Last Name <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              value={formData.prescriberLastName}
              onChange={(e) => handleChange('prescriberLastName', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="Smith"
            />
            {errors.prescriberLastName && <p className="text-xs text-red-600 mt-1">{errors.prescriberLastName}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Street Address</label>
            <input
              type="text"
              value={formData.prescriberStreet}
              onChange={(e) => handleChange('prescriberStreet', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="456 Medical Blvd"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">City</label>
            <input
              type="text"
              value={formData.prescriberCity}
              onChange={(e) => handleChange('prescriberCity', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="New York"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">State</label>
            <input
              type="text"
              value={formData.prescriberState}
              onChange={(e) => handleChange('prescriberState', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="NY"
              maxLength={2}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Zip Code</label>
            <input
              type="text"
              value={formData.prescriberZipCode}
              onChange={(e) => handleChange('prescriberZipCode', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="10002"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Phone</label>
            <input
              type="tel"
              value={formData.prescriberPhone}
              onChange={(e) => handleChange('prescriberPhone', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="555-5678"
            />
          </div>
        </div>
      </div>

      {/* Medication Information Section */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-slate-900 dark:text-white border-b border-slate-200 dark:border-slate-700 pb-2">
          Medication Information
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              NDC <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              value={formData.medicationNDC}
              onChange={(e) => handleChange('medicationNDC', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="00002-7510-02"
            />
            {errors.medicationNDC && <p className="text-xs text-red-600 mt-1">{errors.medicationNDC}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              Medication Name <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              value={formData.medicationName}
              onChange={(e) => handleChange('medicationName', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="Humira"
            />
            {errors.medicationName && <p className="text-xs text-red-600 mt-1">{errors.medicationName}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              Quantity <span className="text-red-500">*</span>
            </label>
            <input
              type="number"
              value={formData.medicationQuantity}
              onChange={(e) => handleChange('medicationQuantity', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="2"
              min="1"
            />
            {errors.medicationQuantity && <p className="text-xs text-red-600 mt-1">{errors.medicationQuantity}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Refills</label>
            <input
              type="number"
              value={formData.medicationRefills}
              onChange={(e) => handleChange('medicationRefills', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="3"
              min="0"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Dosage</label>
            <input
              type="text"
              value={formData.medicationDosage}
              onChange={(e) => handleChange('medicationDosage', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="40mg"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Date Written</label>
            <input
              type="date"
              value={formData.dateWritten}
              onChange={(e) => handleChange('dateWritten', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
            />
          </div>
          <div className="md:col-span-2">
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Directions</label>
            <textarea
              value={formData.medicationDirections}
              onChange={(e) => handleChange('medicationDirections', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              rows={3}
              placeholder="Take 2 injections every 2 weeks"
            />
          </div>
        </div>
      </div>

      {/* Insurance Information Section */}
      <div className="space-y-4">
        <h3 className="text-lg font-semibold text-slate-900 dark:text-white border-b border-slate-200 dark:border-slate-700 pb-2">
          Insurance Information (Optional)
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">BIN</label>
            <input
              type="text"
              value={formData.insuranceBIN}
              onChange={(e) => handleChange('insuranceBIN', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="004682"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">PCN</label>
            <input
              type="text"
              value={formData.insurancePCN}
              onChange={(e) => handleChange('insurancePCN', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="CNRX"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Group ID</label>
            <input
              type="text"
              value={formData.insuranceGroupID}
              onChange={(e) => handleChange('insuranceGroupID', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="GROUP123"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Member ID</label>
            <input
              type="text"
              value={formData.insuranceMemberID}
              onChange={(e) => handleChange('insuranceMemberID', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="MEM123456"
            />
          </div>
          <div className="md:col-span-2">
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Plan Name</label>
            <input
              type="text"
              value={formData.insurancePlanName}
              onChange={(e) => handleChange('insurancePlanName', e.target.value)}
              className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white"
              placeholder="Blue Cross Blue Shield"
            />
          </div>
        </div>
      </div>
    </div>
  )
}

