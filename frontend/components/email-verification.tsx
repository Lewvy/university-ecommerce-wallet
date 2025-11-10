"use client"

import type React from "react"
import { useState, useEffect } from "react"

interface EmailVerificationPageProps {
  email: string
  password: string
  onVerificationComplete: () => void
}

export default function EmailVerificationPage({ email, password, onVerificationComplete }: EmailVerificationPageProps) {
  const [token, setToken] = useState("")
  const [timeLeft, setTimeLeft] = useState(60)
  const [isVerifying, setIsVerifying] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft((prev) => (prev > 0 ? prev - 1 : 0))
    }, 1000)
    return () => clearInterval(timer)
  }, [])

  const handleVerify = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (token.length === 0) {
      setError("Please enter the verification token")
      return
    }

    setIsVerifying(true)
    setError(null)

    try {
      // First, verify the email with the token
      const verifyResponse = await fetch("http://localhost:8088/verify-email", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: email,
          token: token
        }),
      })

      const verifyData = await verifyResponse.json()

      if (verifyResponse.ok) {
        console.log("Email verification successful:", verifyData)
        
        // After successful verification, automatically log in
        const loginResponse = await fetch("http://localhost:8088/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            email: email,
            password: password
          }),
        })

        const loginData = await loginResponse.json()

        if (loginResponse.ok) {
          console.log("Auto-login successful:", loginData)
          
          // Store tokens in memory (you can also use context or state management)
          if (loginData.access_token) {
            sessionStorage.setItem('access_token', loginData.access_token)
          }
          if (loginData.refresh_token) {
            sessionStorage.setItem('refresh_token', loginData.refresh_token)
          }
          
          onVerificationComplete()
        } else {
          setError(loginData.message || "Login failed after verification. Please try logging in manually.")
          setTimeout(() => {
            onVerificationComplete()
          }, 2000)
        }
      } else {
        console.error("Email verification failed:", verifyData)
        setError(verifyData.message || "Invalid verification token. Please try again.")
      }
    } catch (error) {
      console.error("Verification error:", error)
      setError("Network error. Please check your connection and try again.")
    } finally {
      setIsVerifying(false)
    }
  }

  const handleResend = async () => {
    try {
      const response = await fetch("http://localhost:8088/resend-verification", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email }),
      })

      if (response.ok) {
        setTimeLeft(60)
        setError(null)
        alert("Verification code resent successfully!")
      } else {
        const data = await response.json()
        setError(data.message || "Failed to resend verification code")
      }
    } catch (error) {
      console.error("Resend error:", error)
      setError("Failed to resend verification code")
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-white flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="bg-white rounded-lg shadow-lg p-8 text-center">
          <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <svg className="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
              />
            </svg>
          </div>

          <h2 className="text-2xl font-bold text-gray-900 mb-2">Verify Your Email</h2>
          <p className="text-gray-600 mb-6">
            We sent a verification token to <span className="font-semibold">{email}</span>
          </p>

          <form onSubmit={handleVerify} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Verification Token</label>
              <input
                type="text"
                value={token}
                onChange={(e) => {
                  setToken(e.target.value)
                  setError(null)
                }}
                placeholder="Enter verification code"
                required
                disabled={isVerifying}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-center text-lg tracking-widest disabled:opacity-50"
              />
            </div>

            {error && (
              <p className="text-sm text-red-600 font-medium">{error}</p>
            )}

            <button
              type="submit"
              disabled={isVerifying}
              className="w-full bg-gradient-to-r from-blue-600 to-blue-500 text-white font-semibold py-2 rounded-lg hover:shadow-lg transition-shadow disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isVerifying ? 'Verifying...' : 'Verify Email'}
            </button>
          </form>

          <div className="mt-4">
            {timeLeft > 0 ? (
              <p className="text-gray-600 text-sm">
                Resend code in <span className="font-semibold text-blue-600">{timeLeft}s</span>
              </p>
            ) : (
              <button
                onClick={handleResend}
                className="text-blue-600 font-semibold text-sm hover:underline"
              >
                Resend verification code
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

/*"use client"

import type React from "react"

import { useState, useEffect } from "react"

interface EmailVerificationPageProps {
  email: string
  onVerificationComplete: () => void
}

export default function EmailVerificationPage({ email, onVerificationComplete }: EmailVerificationPageProps) {
  const [token, setToken] = useState("")
  const [timeLeft, setTimeLeft] = useState(60)

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft((prev) => (prev > 0 ? prev - 1 : 0))
    }, 1000)
    return () => clearInterval(timer)
  }, [])

  const handleVerify = (e: React.FormEvent) => {
    e.preventDefault()
    if (token.length > 0) {
      onVerificationComplete()
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-white flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="bg-white rounded-lg shadow-lg p-8 text-center">
          <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <svg className="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
              />
            </svg>
          </div>

          <h2 className="text-2xl font-bold text-gray-900 mb-2">Verify Your Email</h2>
          <p className="text-gray-600 mb-6">
            We sent a verification token to <span className="font-semibold">{email}</span>
          </p>

          <form onSubmit={handleVerify} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Verification Token</label>
              <input
                type="text"
                value={token}
                onChange={(e) => setToken(e.target.value)}
                placeholder="Enter verification code"
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-center text-lg tracking-widest"
              />
            </div>

            <button
              type="submit"
              className="w-full bg-gradient-to-r from-blue-600 to-blue-500 text-white font-semibold py-2 rounded-lg hover:shadow-lg transition-shadow"
            >
              Verify Email
            </button>
          </form>

          <p className="text-gray-600 mt-4 text-sm">
            Resend in <span className="font-semibold text-blue-600">{timeLeft}s</span>
          </p>
        </div>
      </div>
    </div>
  )
}*/
