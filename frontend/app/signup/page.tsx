"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import EmailVerificationPage from "../../components/email-verification"
import SignupForm from "../../components/signup-form"

interface FormData {
	email: string;
	name: string;
	phone: string;
	password: string;
}

export default function SignupPage() {
	const router = useRouter()

	const [step, setStep] = useState<"form" | "verification" | "success">("form")
	const [userEmail, setUserEmail] = useState("")
	const [userName, setUserName] = useState("")
	const [userPhone, setUserPhone] = useState("")
	const [password, setPassword] = useState("")
	const [userId, setUserId] = useState<number | null>(null)

	
	const handleFormSubmit = async (data: FormData) => {
		try {
			console.log("Sending registration request:", data)

			const response = await fetch("http://localhost:8088/register", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(data),
			})

			const responseData = await response.json()
			console.log("Registration response:", response.status, responseData)

			if (response.ok) {
				// Save user info
				setUserEmail(data.email)
				setUserName(data.name)
				setUserPhone(data.phone)
				setPassword(data.password)

				// Save userId if backend returns it
				if (responseData.id) {
					setUserId(responseData.id)
				} else if (responseData.user && responseData.user.id) {
					setUserId(responseData.user.id)
				} else {
					console.warn("⚠️ Backend returned no userId")
				}

				setStep("verification") // Move to OTP screen
			} else {
				alert(responseData.message || "Registration failed. Try again.")
			}
		} catch (err) {
			console.error("Network error:", err)
			alert("Network error. Try again.")
		}
	}

	// ------------------------------
	// AFTER SUCCESSFUL VERIFICATION
	// ------------------------------
	const handleVerificationComplete = async () => {
		console.log("Verification complete, preparing to redirect...")
		
		// Store user data in sessionStorage before redirecting
		const userInfo = {
			username: userName,
			email: userEmail,
			phone: userPhone,
			id: userId
		}
		
		console.log("Storing user info:", userInfo)
		localStorage.setItem('user', JSON.stringify(userInfo))
		
		// Show success screen briefly
		setStep("success")

		// Wait a bit then redirect
		await new Promise(resolve => setTimeout(resolve, 1500))
		
		console.log("Redirecting to dashboard...")
		router.push("/login")
	}

	// ------------------------------
	// RENDERING SECTIONS
	// ------------------------------

	if (step === "form") {
		return <SignupForm onSubmit={handleFormSubmit} />
	}

	if (step === "verification") {
		return (
			<EmailVerificationPage
				email={userEmail}
				userId={userId}
				password={password}
				onVerificationComplete={handleVerificationComplete}
			/>
		)
	}

	// SUCCESS PAGE
	return (
		<div className="min-h-screen bg-gradient-to-br from-blue-50 to-white flex items-center justify-center p-4">
			<div className="text-center">
				<div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
					<svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
					</svg>
				</div>
				<h2 className="text-3xl font-bold text-gray-900 mb-2">Signup Successful!</h2>
				<p className="text-gray-600">Redirecting to login...</p>
			</div>
		</div>
	)
}