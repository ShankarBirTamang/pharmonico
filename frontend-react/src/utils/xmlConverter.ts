// Utility to convert prescription form data to NCPDP XML format

import { PrescriptionFormData } from '../components/forms/PrescriptionForm'

/**
 * Convert prescription form data to NCPDP XML format
 */
export function convertFormDataToXML(formData: PrescriptionFormData): string {
  const messageId = `MSG${Date.now()}`
  const timestamp = new Date().toISOString()
  const dateWritten = formData.dateWritten || new Date().toISOString().split('T')[0]

  let xml = `<?xml version="1.0" encoding="UTF-8"?>\n`
  xml += `<Message>\n`
  xml += `  <Header>\n`
  xml += `    <MessageID>${messageId}</MessageID>\n`
  xml += `    <Timestamp>${timestamp}</Timestamp>\n`
  xml += `  </Header>\n`
  xml += `  <Body>\n`
  xml += `    <Prescription DateWritten="${dateWritten}">\n`

  // Patient Information
  xml += `      <Patient ID="${formData.patientId || 'PAT001'}">\n`
  xml += `        <FirstName>${escapeXML(formData.patientFirstName)}</FirstName>\n`
  xml += `        <LastName>${escapeXML(formData.patientLastName)}</LastName>\n`
  xml += `        <DateOfBirth>${formData.patientDateOfBirth}</DateOfBirth>\n`
  if (formData.patientStreet || formData.patientCity || formData.patientState || formData.patientZipCode) {
    xml += `        <Address>\n`
    if (formData.patientStreet) xml += `          <Street>${escapeXML(formData.patientStreet)}</Street>\n`
    if (formData.patientCity) xml += `          <City>${escapeXML(formData.patientCity)}</City>\n`
    if (formData.patientState) xml += `          <State>${escapeXML(formData.patientState)}</State>\n`
    if (formData.patientZipCode) xml += `          <ZipCode>${escapeXML(formData.patientZipCode)}</ZipCode>\n`
    xml += `        </Address>\n`
  }
  if (formData.patientPhone) {
    xml += `        <Phone>${escapeXML(formData.patientPhone)}</Phone>\n`
  }
  xml += `      </Patient>\n`

  // Prescriber Information
  xml += `      <Prescriber ID="${formData.prescriberId || 'PRES001'}">\n`
  xml += `        <NPI>${escapeXML(formData.prescriberNPI)}</NPI>\n`
  if (formData.prescriberDEA) {
    xml += `        <DEA>${escapeXML(formData.prescriberDEA)}</DEA>\n`
  }
  xml += `        <FirstName>${escapeXML(formData.prescriberFirstName)}</FirstName>\n`
  xml += `        <LastName>${escapeXML(formData.prescriberLastName)}</LastName>\n`
  if (formData.prescriberStreet || formData.prescriberCity || formData.prescriberState || formData.prescriberZipCode) {
    xml += `        <Address>\n`
    if (formData.prescriberStreet) xml += `          <Street>${escapeXML(formData.prescriberStreet)}</Street>\n`
    if (formData.prescriberCity) xml += `          <City>${escapeXML(formData.prescriberCity)}</City>\n`
    if (formData.prescriberState) xml += `          <State>${escapeXML(formData.prescriberState)}</State>\n`
    if (formData.prescriberZipCode) xml += `          <ZipCode>${escapeXML(formData.prescriberZipCode)}</ZipCode>\n`
    xml += `        </Address>\n`
  }
  if (formData.prescriberPhone) {
    xml += `        <Phone>${escapeXML(formData.prescriberPhone)}</Phone>\n`
  }
  xml += `      </Prescriber>\n`

  // Medication Information
  xml += `      <Medication>\n`
  xml += `        <NDC>${escapeXML(formData.medicationNDC)}</NDC>\n`
  xml += `        <Name>${escapeXML(formData.medicationName)}</Name>\n`
  xml += `        <Quantity>${formData.medicationQuantity || 1}</Quantity>\n`
  if (formData.medicationRefills) {
    xml += `        <Refills>${formData.medicationRefills}</Refills>\n`
  }
  if (formData.medicationDosage) {
    xml += `        <Dosage>${escapeXML(formData.medicationDosage)}</Dosage>\n`
  }
  if (formData.medicationDirections) {
    xml += `        <Directions>${escapeXML(formData.medicationDirections)}</Directions>\n`
  }
  xml += `      </Medication>\n`

  // Insurance Information (optional)
  if (
    formData.insuranceBIN ||
    formData.insurancePCN ||
    formData.insuranceGroupID ||
    formData.insuranceMemberID ||
    formData.insurancePlanName
  ) {
    xml += `      <Insurance>\n`
    if (formData.insuranceBIN) xml += `        <BIN>${escapeXML(formData.insuranceBIN)}</BIN>\n`
    if (formData.insurancePCN) xml += `        <PCN>${escapeXML(formData.insurancePCN)}</PCN>\n`
    if (formData.insuranceGroupID) xml += `        <GroupID>${escapeXML(formData.insuranceGroupID)}</GroupID>\n`
    if (formData.insuranceMemberID) xml += `        <MemberID>${escapeXML(formData.insuranceMemberID)}</MemberID>\n`
    if (formData.insurancePlanName) xml += `        <PlanName>${escapeXML(formData.insurancePlanName)}</PlanName>\n`
    xml += `      </Insurance>\n`
  }

  xml += `    </Prescription>\n`
  xml += `  </Body>\n`
  xml += `</Message>`

  return xml
}

/**
 * Escape XML special characters
 */
function escapeXML(str: string): string {
  if (!str) return ''
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&apos;')
}

/**
 * Convert  example to form data
 */
export function convertExampleToFormData(example: any): PrescriptionFormData {
  return {
    // Patient
    patientId: example.patient?.id || '',
    patientFirstName: example.patient?.firstName || '',
    patientLastName: example.patient?.lastName || '',
    patientDateOfBirth: example.patient?.dateOfBirth || '',
    patientStreet: example.patient?.street || '',
    patientCity: example.patient?.city || '',
    patientState: example.patient?.state || '',
    patientZipCode: example.patient?.zipCode || '',
    patientPhone: example.patient?.phone || '',

    // Prescriber
    prescriberId: example.prescriber?.id || '',
    prescriberNPI: example.prescriber?.npi || '',
    prescriberDEA: example.prescriber?.dea || '',
    prescriberFirstName: example.prescriber?.firstName || '',
    prescriberLastName: example.prescriber?.lastName || '',
    prescriberStreet: example.prescriber?.street || '',
    prescriberCity: example.prescriber?.city || '',
    prescriberState: example.prescriber?.state || '',
    prescriberZipCode: example.prescriber?.zipCode || '',
    prescriberPhone: example.prescriber?.phone || '',

    // Medication
    medicationNDC: example.medication?.ndc || '',
    medicationName: example.medication?.name || '',
    medicationQuantity: example.medication?.quantity?.toString() || '',
    medicationRefills: example.medication?.refills?.toString() || '',
    medicationDosage: example.medication?.dosage || '',
    medicationDirections: example.medication?.directions || '',
    dateWritten: example.dateWritten || new Date().toISOString().split('T')[0],

    // Insurance
    insuranceBIN: example.insurance?.bin || '',
    insurancePCN: example.insurance?.pcn || '',
    insuranceGroupID: example.insurance?.groupID || '',
    insuranceMemberID: example.insurance?.memberID || '',
    insurancePlanName: example.insurance?.planName || '',
  }
}

