"use client"

import { useState } from "react"
import EmailVerificationPage from "./email-verification"
import SignupForm from "./signup-form"

interface SignupPageProps {
  onSignupSuccess: (username: string, email: string, phone: string) => void
}

export default function SignupPage({ onSignupSuccess }: SignupPageProps) {
  const [step, setStep] = useState<"form" | "verification" | "success">("form")
  const [userEmail, setUserEmail] = useState("")
  const [username, setUsername] = useState("")
  const [userPhone, setUserPhone] = useState("")

  const handleFormSubmit = (email: string, name: string, phone: string) => {
    setUserEmail(email)
    setUsername(name)
    setUserPhone(phone)
    setStep("verification")
  }

  const handleVerificationComplete = () => {
    setStep("success")
    setTimeout(() => {
      onSignupSuccess(username, userEmail, userPhone)
    }, 2000)
  }

  if (step === "form") {
    return <SignupForm onSubmit={handleFormSubmit} />
  }

  if (step === "verification") {
    return <EmailVerificationPage email={userEmail} onVerificationComplete={handleVerificationComplete} />
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-white flex items-center justify-center p-4">
      <div className="text-center">
        <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h2 className="text-3xl font-bold text-gray-900 mb-2">Signup Successful!</h2>
        <p className="text-gray-600">Redirecting to home...</p>
      </div>
    </div>
  )
}
