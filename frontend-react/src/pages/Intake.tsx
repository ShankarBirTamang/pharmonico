import { useState } from 'react'
import { Card, CardHeader, CardTitle, CardContent } from '../components/ui'
import { api } from '../services/api'

export function Intake() {
  const [xmlPayload, setXmlPayload] = useState('')
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<{ prescription_id?: string; error?: string } | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setResult(null)

    try {
      const response = await api.submitPrescription(xmlPayload, 'xml')
      setResult({ prescription_id: response.prescription_id })
      setXmlPayload('') // Clear form on success
    } catch (error) {
      setResult({ error: error instanceof Error ? error.message : 'Failed to submit prescription' })
    } finally {
      setLoading(false)
    }
  }

  const exampleXML = `<Message>
  <Header>
    <MessageID>MSG001</MessageID>
    <Timestamp>2024-01-15T10:30:00Z</Timestamp>
  </Header>
  <Body>
    <Prescription DateWritten="2024-01-15">
      <Patient ID="PAT001">
        <FirstName>John</FirstName>
        <LastName>Doe</LastName>
        <DateOfBirth>1980-05-15</DateOfBirth>
        <Address>
          <Street>123 Main St</Street>
          <City>New York</City>
          <State>NY</State>
          <ZipCode>10001</ZipCode>
        </Address>
        <Phone>555-1234</Phone>
      </Patient>
      <Prescriber ID="PRES001">
        <NPI>1234567890</NPI>
        <DEA>AB1234567</DEA>
        <FirstName>Jane</FirstName>
        <LastName>Smith</LastName>
        <Address>
          <Street>456 Medical Blvd</Street>
          <City>New York</City>
          <State>NY</State>
          <ZipCode>10002</ZipCode>
        </Address>
        <Phone>555-5678</Phone>
      </Prescriber>
      <Medication>
        <NDC>00002-7510-02</NDC>
        <Name>Humira</Name>
        <Quantity>2</Quantity>
        <Refills>3</Refills>
        <Dosage>40mg</Dosage>
        <Directions>Take 2 injections every 2 weeks</Directions>
      </Medication>
      <Insurance>
        <BIN>004682</BIN>
        <PCN>CNRX</PCN>
        <MemberID>MEM123456</MemberID>
        <PlanName>Blue Cross Blue Shield</PlanName>
      </Insurance>
    </Prescription>
  </Body>
</Message>`

  const loadExample = () => {
    setXmlPayload(exampleXML)
    setResult(null)
  }

  return (
    <div className="max-w-4xl mx-auto">
      <Card>
        <CardHeader>
          <CardTitle>Prescription Intake</CardTitle>
          <p className="text-sm text-slate-600 dark:text-slate-400 mt-2">
            Submit NCPDP SCRIPT prescription in XML format as a Healthcare Provider
          </p>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <div className="flex items-center justify-between mb-2">
                <label htmlFor="xml-payload" className="block text-sm font-medium text-slate-700 dark:text-slate-300">
                  NCPDP XML Payload
                </label>
                <button
                  type="button"
                  onClick={loadExample}
                  className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
                >
                  Load Example
                </button>
              </div>
              <textarea
                id="xml-payload"
                value={xmlPayload}
                onChange={(e) => setXmlPayload(e.target.value)}
                rows={20}
                className="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-slate-800 dark:border-slate-600 dark:text-white font-mono text-sm"
                placeholder="Paste or enter NCPDP XML here..."
                required
              />
            </div>

            <button
              type="submit"
              disabled={loading || !xmlPayload.trim()}
              className="btn-primary w-full disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? (
                <span className="flex items-center justify-center gap-2">
                  <span className="inline-block h-4 w-4 animate-spin rounded-full border-2 border-solid border-white border-r-transparent"></span>
                  Submitting...
                </span>
              ) : (
                'Submit Prescription'
              )}
            </button>
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

          {/* Example XML */}
          <div className="mt-6 p-4 bg-slate-50 dark:bg-slate-900 rounded-lg border border-slate-200 dark:border-slate-700">
            <p className="text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">
              Example NCPDP XML Format:
            </p>
            <pre className="text-xs text-slate-600 dark:text-slate-400 overflow-x-auto whitespace-pre-wrap">
              {exampleXML}
            </pre>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

