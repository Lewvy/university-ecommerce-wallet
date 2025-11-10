// lib/auth.ts
// Authentication utility functions

const API_BASE_URL = 'http://localhost:8088'

export interface AuthTokens {
  access_token: string
  refresh_token: string
}

export interface User {
  id: number
  name: string
  email: string
  phone: string
}

// Token management
export const getAccessToken = (): string | null => {
  if (typeof window === 'undefined') return null
  return sessionStorage.getItem('access_token')
}

export const getRefreshToken = (): string | null => {
  if (typeof window === 'undefined') return null
  return sessionStorage.getItem('refresh_token')
}

export const setTokens = (tokens: AuthTokens): void => {
  if (typeof window === 'undefined') return
  sessionStorage.setItem('access_token', tokens.access_token)
  sessionStorage.setItem('refresh_token', tokens.refresh_token)
}

export const clearTokens = (): void => {
  if (typeof window === 'undefined') return
  sessionStorage.removeItem('access_token')
  sessionStorage.removeItem('refresh_token')
  sessionStorage.removeItem('user')
}

export const isAuthenticated = (): boolean => {
  return !!getAccessToken()
}

// API fetch wrapper with authentication
export const authenticatedFetch = async (
  endpoint: string,
  options: RequestInit = {}
): Promise<Response> => {
  const token = getAccessToken()
  
  const headers = {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
    ...options.headers,
  }

  let response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  })

  // If token expired, try to refresh
  if (response.status === 401 && getRefreshToken()) {
    const newToken = await refreshAccessToken()
    if (newToken) {
      // Retry the original request with new token
      headers.Authorization = `Bearer ${newToken}`
      response = await fetch(`${API_BASE_URL}${endpoint}`, {
        ...options,
        headers,
      })
    }
  }

  return response
}

// Refresh access token
export const refreshAccessToken = async (): Promise<string | null> => {
  const refreshToken = getRefreshToken()
  if (!refreshToken) return null

  try {
    const response = await fetch(`${API_BASE_URL}/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ refresh_token: refreshToken }),
    })

    if (response.ok) {
      const data = await response.json()
      if (data.access_token) {
        sessionStorage.setItem('access_token', data.access_token)
        return data.access_token
      }
    }
  } catch (error) {
    console.error('Token refresh failed:', error)
  }

  // If refresh fails, clear tokens
  clearTokens()
  return null
}

// Get current user
export const getCurrentUser = (): User | null => {
  if (typeof window === 'undefined') return null
  const userStr = sessionStorage.getItem('user')
  return userStr ? JSON.parse(userStr) : null
}

// Logout
export const logout = (): void => {
  clearTokens()
  if (typeof window !== 'undefined') {
    window.location.href = '/login'
  }
}