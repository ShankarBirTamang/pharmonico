import { useParams } from 'react-router-dom'
import { Card, CardHeader, CardTitle, CardContent } from '../components/ui'
import { useEffect, useState } from 'react'
import { api } from '../services/api'

interface EnrollmentData {
  token: string
  patientId?: string
  status?: string
}

export function Enroll() {
  const { token } = useParams<{ token: string }>()
  const [enrollment, setEnrollment] = useState<EnrollmentData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!token) {
      setError('Invalid enrollment token')
      setLoading(false)
      return
    }

    // TODO: Replace with actual enrollment API call when endpoint is available
    // For now, just set the token
    setEnrollment({ token })
    setLoading(false)
  }, [token])

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <div className="mb-4 inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-primary-600 border-r-transparent"></div>
          <p className="text-slate-600 dark:text-slate-400">Loading enrollment...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle className="text-red-600 dark:text-red-400">Enrollment Error</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-slate-600 dark:text-slate-400">{error}</p>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-50 via-blue-50 to-teal-50 dark:from-slate-950 dark:via-slate-900 dark:to-slate-950 px-4 py-12">
      <Card className="w-full max-w-2xl">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-primary-500 to-teal-500">
              <span className="text-2xl font-bold text-white">P</span>
            </div>
            <div>
              <CardTitle>Patient Enrollment</CardTitle>
              <p className="mt-1 text-sm text-slate-600 dark:text-slate-400">
                Complete your enrollment to access prescription services
              </p>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="rounded-lg bg-slate-50 p-4 dark:bg-slate-900">
            <p className="text-sm font-medium text-slate-700 dark:text-slate-300">Enrollment Token</p>
            <p className="mt-1 font-mono text-sm text-slate-900 dark:text-white">{token}</p>
          </div>

          <div className="rounded-lg border border-amber-200 bg-amber-50 p-4 dark:border-amber-800 dark:bg-amber-900/20">
            <p className="text-sm text-amber-800 dark:text-amber-400">
              <strong>Note:</strong> Enrollment API integration is pending. This page will be connected to the backend enrollment endpoint once available.
            </p>
          </div>

          <div className="flex gap-4">
            <button className="btn-primary flex-1">Continue Enrollment</button>
            <button className="btn-secondary">Cancel</button>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

