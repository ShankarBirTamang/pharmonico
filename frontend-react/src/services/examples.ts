// Prescription example service - provides static prescription examples

export interface PrescriptionExample {
  patient: {
    id: string
    firstName: string
    lastName: string
    dateOfBirth: string
    street?: string
    city?: string
    state?: string
    zipCode?: string
    phone?: string
  }
  prescriber: {
    id: string
    npi: string
    dea?: string
    firstName: string
    lastName: string
    street?: string
    city?: string
    state?: string
    zipCode?: string
    phone?: string
  }
  medication: {
    ndc: string
    name: string
    quantity: number
    refills?: number
    dosage?: string
    directions?: string
  }
  dateWritten: string
  insurance?: {
    bin?: string
    pcn?: string
    groupID?: string
    memberID?: string
    planName?: string
  }
}

/**
 * Get a random prescription example from static examples
 */
export function getPrescriptionExample(): PrescriptionExample {
  const examples: PrescriptionExample[] = [
    {
      patient: {
        id: 'PAT001',
        firstName: 'John',
        lastName: 'Doe',
        dateOfBirth: '1980-05-15',
        street: '123 Main St',
        city: 'New York',
        state: 'NY',
        zipCode: '10001',
        phone: '555-1234',
      },
      prescriber: {
        id: 'PRES001',
        npi: '1234567890',
        dea: 'AB1234567',
        firstName: 'Jane',
        lastName: 'Smith',
        street: '456 Medical Blvd',
        city: 'New York',
        state: 'NY',
        zipCode: '10002',
        phone: '555-5678',
      },
      medication: {
        ndc: '00002-7510-02',
        name: 'Humira',
        quantity: 2,
        refills: 3,
        dosage: '40mg',
        directions: 'Take 2 injections every 2 weeks',
      },
      dateWritten: '2024-01-15',
      insurance: {
        bin: '004682',
        pcn: 'CNRX',
        memberID: 'MEM123456',
        planName: 'Blue Cross Blue Shield',
      },
    },
    {
      patient: {
        id: 'PAT002',
        firstName: 'Sarah',
        lastName: 'Johnson',
        dateOfBirth: '1975-08-22',
        street: '789 Oak Avenue',
        city: 'Los Angeles',
        state: 'CA',
        zipCode: '90001',
        phone: '555-9876',
      },
      prescriber: {
        id: 'PRES002',
        npi: '9876543210',
        dea: 'CD9876543',
        firstName: 'Michael',
        lastName: 'Brown',
        street: '321 Health Center Dr',
        city: 'Los Angeles',
        state: 'CA',
        zipCode: '90002',
        phone: '555-5432',
      },
      medication: {
        ndc: '68180-500-01',
        name: 'Enbrel',
        quantity: 4,
        refills: 2,
        dosage: '50mg',
        directions: 'Inject once weekly',
      },
      dateWritten: '2024-01-20',
      insurance: {
        bin: '003585',
        pcn: 'ADV',
        memberID: 'MEM789012',
        planName: 'Aetna',
      },
    },
    {
      patient: {
        id: 'PAT003',
        firstName: 'Robert',
        lastName: 'Williams',
        dateOfBirth: '1990-12-05',
        street: '456 Elm Street',
        city: 'Chicago',
        state: 'IL',
        zipCode: '60601',
        phone: '555-2468',
      },
      prescriber: {
        id: 'PRES003',
        npi: '1122334455',
        dea: 'EF1122334',
        firstName: 'Emily',
        lastName: 'Davis',
        street: '789 Clinic Way',
        city: 'Chicago',
        state: 'IL',
        zipCode: '60602',
        phone: '555-1357',
      },
      medication: {
        ndc: '50458-300-01',
        name: 'Stelara',
        quantity: 1,
        refills: 5,
        dosage: '45mg/0.5mL',
        directions: 'Inject subcutaneously every 12 weeks',
      },
      dateWritten: '2024-01-25',
      insurance: {
        bin: '004682',
        pcn: 'CNRX',
        memberID: 'MEM345678',
        planName: 'UnitedHealthcare',
      },
    },
    {
      patient: {
        id: 'PAT004',
        firstName: 'Maria',
        lastName: 'Garcia',
        dateOfBirth: '1985-03-18',
        street: '321 Pine Street',
        city: 'Houston',
        state: 'TX',
        zipCode: '77001',
        phone: '555-3691',
      },
      prescriber: {
        id: 'PRES004',
        npi: '5566778899',
        dea: 'GH5566778',
        firstName: 'David',
        lastName: 'Wilson',
        street: '654 Medical Plaza',
        city: 'Houston',
        state: 'TX',
        zipCode: '77002',
        phone: '555-7410',
      },
      medication: {
        ndc: '50242-040-01',
        name: 'Remicade',
        quantity: 1,
        refills: 2,
        dosage: '100mg',
        directions: 'Infuse intravenously every 8 weeks',
      },
      dateWritten: '2024-01-28',
      insurance: {
        bin: '003585',
        pcn: 'ADV',
        memberID: 'MEM456789',
        planName: 'Cigna',
      },
    },
    {
      patient: {
        id: 'PAT005',
        firstName: 'James',
        lastName: 'Anderson',
        dateOfBirth: '1972-11-30',
        street: '987 Maple Drive',
        city: 'Phoenix',
        state: 'AZ',
        zipCode: '85001',
        phone: '555-8520',
      },
      prescriber: {
        id: 'PRES005',
        npi: '2233445566',
        dea: 'IJ2233445',
        firstName: 'Lisa',
        lastName: 'Martinez',
        street: '147 Healthcare Blvd',
        city: 'Phoenix',
        state: 'AZ',
        zipCode: '85002',
        phone: '555-9630',
      },
      medication: {
        ndc: '00069-1010-01',
        name: 'Cosentyx',
        quantity: 2,
        refills: 4,
        dosage: '150mg',
        directions: 'Inject subcutaneously once monthly',
      },
      dateWritten: '2024-02-01',
      insurance: {
        bin: '004682',
        pcn: 'CNRX',
        memberID: 'MEM567890',
        planName: 'Anthem',
      },
    },
  ]

  // Return a random example
  return examples[Math.floor(Math.random() * examples.length)]
}
