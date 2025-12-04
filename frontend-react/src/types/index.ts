// Prescription Types
export interface Prescription {
  id: string
  patientId: string
  prescriberId: string
  status: PrescriptionStatus
  medications: Medication[]
  createdAt: string
  updatedAt: string
}

export type PrescriptionStatus =
  | 'received'
  | 'validated'
  | 'validation_issue'
  | 'awaiting_enrollment'
  | 'awaiting_routing'
  | 'routed'
  | 'fulfilled'

export interface Medication {
  name: string
  dosage: string
  quantity: number
  refills: number
}

// Pharmacy Types
export interface Pharmacy {
  id: string
  name: string
  address: Address
  phone: string
  email: string
  specialties: string[]
  acceptedInsurers: string[]
  capacity: number
  loadFactor: number
  coordinates: GeoCoordinates
}

export interface Address {
  street: string
  city: string
  state: string
  zipCode: string
}

export interface GeoCoordinates {
  lat: number
  lng: number
}

// Patient Types
export interface Patient {
  id: string
  firstName: string
  lastName: string
  dateOfBirth: string
  insurance: Insurance
  address: Address
}

export interface Insurance {
  provider: string
  planId: string
  memberId: string
}

// API Response Types
export interface ApiResponse<T> {
  data: T
  message?: string
  success: boolean
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

