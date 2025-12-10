// API Configuration and Base Client

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'
const API_SERVER_URL = import.meta.env.VITE_API_URL?.replace('/api', '') || 'http://localhost:8080'

interface RequestOptions extends RequestInit {
  params?: Record<string, string>
}

interface HealthResponse {
  status: string
  message: string
}

class ApiClient {
  private baseUrl: string
  private serverUrl: string

  constructor(baseUrl: string, serverUrl: string) {
    this.baseUrl = baseUrl
    this.serverUrl = serverUrl
  }

  private async request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
    const { params, ...fetchOptions } = options
    
    let url = `${this.baseUrl}${endpoint}`
    if (params) {
      const searchParams = new URLSearchParams(params)
      url += `?${searchParams.toString()}`
    }

    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...fetchOptions.headers,
    }

    // Add auth token if available
    const token = localStorage.getItem('auth_token')
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    const response = await fetch(url, {
      ...fetchOptions,
      headers,
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: 'Request failed' }))
      throw new Error(error.message || `HTTP error! status: ${response.status}`)
    }

    return response.json()
  }

  get<T>(endpoint: string, options?: RequestOptions): Promise<T> {
    return this.request<T>(endpoint, { ...options, method: 'GET' })
  }

  post<T>(endpoint: string, data?: unknown, options?: RequestOptions): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  put<T>(endpoint: string, data?: unknown, options?: RequestOptions): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  delete<T>(endpoint: string, options?: RequestOptions): Promise<T> {
    return this.request<T>(endpoint, { ...options, method: 'DELETE' })
  }

  /**
   * Health check endpoint (outside /api/v1)
   * Calls GET /health on the server
   */
  async healthCheck(): Promise<HealthResponse> {
    const response = await fetch(`${this.serverUrl}/health`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })

    if (!response.ok) {
      throw new Error(`Health check failed: ${response.status}`)
    }

    return response.json()
  }

  /**
   * Submit a prescription for intake
   * @param payload - NCPDP XML or JSON string
   * @param format - "xml" or "json" (defaults to "xml")
   */
  async submitPrescription(
    payload: string,
    format: 'xml' | 'json' = 'xml'
  ): Promise<{ prescription_id: string; message?: string }> {
    return this.post<{ prescription_id: string; message?: string }>(
      '/v1/prescriptions/intake',
      {
        payload,
        format,
      }
    )
  }
}

export const api = new ApiClient(API_BASE_URL, API_SERVER_URL)

// Export types
export type { HealthResponse }

