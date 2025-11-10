"use client"

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
}
