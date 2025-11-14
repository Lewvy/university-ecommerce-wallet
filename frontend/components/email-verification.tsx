"use client"

import type React from "react"
import { useState, useEffect } from "react"

interface EmailVerificationPageProps {
	email: string
	userId: number | null
	password: string
	onVerificationComplete: () => void
}

export default function EmailVerificationPage({ email, userId, password, onVerificationComplete }: EmailVerificationPageProps) {
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
			// Build the verification payload
			// Try with different combinations based on what backend might expect
			const verificationPayloads = [
				// Most common: just the token (token contains all needed info)
				{ token: token },
				// Or: token + email
				{ token: token, email: email },
				// Or: token + user_id
				userId ? { token: token, user_id: userId } : null,
				// Or: all three
				userId ? { token: token, email: email, user_id: userId } : null,
			].filter(Boolean) // Remove null entries

			let verificationSuccessful = false
			let lastError = null

			// Try each payload format until one works
			for (const payload of verificationPayloads) {
				console.log("Trying verification with payload:", payload)

				try {
					const verifyResponse = await fetch("http://localhost:8088/verify-email", {
						method: "POST",
						headers: {
							"Content-Type": "application/json",
						},
						body: JSON.stringify(payload),
					})

					const verifyData = await verifyResponse.json()
					console.log("Verification response:", verifyResponse.status, verifyData)

					if (verifyResponse.ok) {
						console.log("Email verification successful!")
						verificationSuccessful = true
						
						// After successful verification, automatically log in
						try {
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

								// Store tokens
								if (loginData.access_token) {
									localStorage.setItem('access_token', loginData.access_token)
								}
								if (loginData.refresh_token) {
									localStorage.setItem('refresh_token', loginData.refresh_token)
								}

								onVerificationComplete()
							} else {
								console.error("Login failed:", loginData)
								setError(loginData.message || "Verification successful but login failed. Please try logging in manually.")
								setTimeout(() => {
									onVerificationComplete()
								}, 3000)
							}
						} catch (loginError) {
							console.error("Login error:", loginError)
							setError("Verification successful but login failed. Please try logging in manually.")
							setTimeout(() => {
								onVerificationComplete()
							}, 3000)
						}
						
						break // Exit the loop since verification was successful
					} else {
						lastError = verifyData.message || "Verification failed"
					}
				} catch (err) {
					console.error("Verification attempt error:", err)
					lastError = "Network error"
				}
			}

			if (!verificationSuccessful) {
				setError(lastError || "Invalid verification token. Please check the code and try again.")
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

			const data = await response.json()

			if (response.ok) {
				setTimeLeft(60)
				setError(null)
				alert("Verification code resent successfully! Check your email.")
			} else {
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
					<p className="text-gray-600 mb-2">
						We sent a verification code to
					</p>
					<p className="font-semibold text-gray-900 mb-6">{email}</p>

					{userId && (
						<p className="text-xs text-gray-500 mb-4">User ID: {userId}</p>
					)}

					<form onSubmit={handleVerify} className="space-y-4">
						<div>
							<label className="block text-sm font-medium text-gray-700 mb-2">Verification Code</label>
							<input
								type="text"
								value={token}
								onChange={(e) => {
									setToken(e.target.value.trim())
									setError(null)
								}}
								placeholder="Enter 6-digit code"
								required
								disabled={isVerifying}
								maxLength={6}
								className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-center text-2xl font-mono tracking-wider disabled:opacity-50"
							/>
						</div>

						{error && (
							<div className="bg-red-50 border border-red-200 rounded-lg p-3">
								<p className="text-sm text-red-600 font-medium">{error}</p>
							</div>
						)}

						<button
							type="submit"
							disabled={isVerifying || token.length === 0}
							className="w-full bg-gradient-to-r from-blue-600 to-blue-500 text-white font-semibold py-3 rounded-lg hover:shadow-lg transition-shadow disabled:opacity-50 disabled:cursor-not-allowed"
						>
							{isVerifying ? (
								<span className="flex items-center justify-center">
									<svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
										<circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
										<path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
									</svg>
									Verifying...
								</span>
							) : 'Verify Email'}
						</button>
					</form>

					<div className="mt-6">
						{timeLeft > 0 ? (
							<p className="text-gray-600 text-sm">
								Didn't receive the code? Resend in <span className="font-semibold text-blue-600">{timeLeft}s</span>
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

					<div className="mt-6 p-3 bg-blue-50 rounded-lg">
						<p className="text-xs text-gray-600 text-left">
							<span className="font-semibold">💡 Tips:</span>
							<br />• Check your email inbox for the verification code
							<br />• Look in your spam/junk folder if you don't see it
							<br />• The code is usually 6 digits
						</p>
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
