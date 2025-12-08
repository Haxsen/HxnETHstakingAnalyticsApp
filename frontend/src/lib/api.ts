import {
  Token,
  TokensResponse,
  TokenHistoryResponse,
  TokenValuationResponse,
  CacheRefreshResponse,
} from './types'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

async function apiRequest<T>(endpoint: string): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`

  try {
    const response = await fetch(url)

    if (!response.ok) {
      throw new ApiError(response.status, `HTTP ${response.status}: ${response.statusText}`)
    }

    const data = await response.json()
    return data
  } catch (error) {
    if (error instanceof ApiError) {
      throw error
    }

    // Network or parsing error
    throw new ApiError(0, error instanceof Error ? error.message : 'Unknown error')
  }
}

// API Functions

export async function fetchTokens(): Promise<TokensResponse> {
  return apiRequest<TokensResponse>('/api/tokens')
}

export async function fetchTokenHistory(symbol: string): Promise<TokenHistoryResponse> {
  return apiRequest<TokenHistoryResponse>(`/api/token/${symbol}/history`)
}

export async function fetchTokenValuation(symbol: string): Promise<TokenValuationResponse> {
  return apiRequest<TokenValuationResponse>(`/api/token/${symbol}/valuation`)
}

export async function refreshCache(): Promise<CacheRefreshResponse> {
  const url = `${API_BASE_URL}/api/cache/refresh`

  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    })

    if (!response.ok) {
      throw new ApiError(response.status, `HTTP ${response.status}: ${response.statusText}`)
    }

    const data = await response.json()
    return data
  } catch (error) {
    if (error instanceof ApiError) {
      throw error
    }

    throw new ApiError(0, error instanceof Error ? error.message : 'Unknown error')
  }
}

// Utility functions for error handling

export function isApiError(error: unknown): error is ApiError {
  return error instanceof ApiError
}

export function getErrorMessage(error: unknown): string {
  if (isApiError(error)) {
    switch (error.status) {
      case 400:
        return 'Invalid request. Please check your input.'
      case 404:
        return 'Data not found.'
      case 500:
        return 'Server error. Please try again later.'
      case 0:
        return 'Network error. Please check your connection.'
      default:
        return error.message
    }
  }

  return error instanceof Error ? error.message : 'An unexpected error occurred'
}
