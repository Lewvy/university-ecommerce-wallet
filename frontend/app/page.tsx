"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"

export default function Page() {
	const router = useRouter()

	useEffect(() => {
		// Check if user is already logged in
		const user = sessionStorage.getItem('user')
		const token = sessionStorage.getItem('access_token')

		if (user && token) {
			// User is logged in, go to dashboard
			router.push('/dashboard')
		} else {
			// User is not logged in, go to login
			router.push('/login')
		}
	}, [router])

	// Show loading while redirecting
	return (
		<div className="min-h-screen bg-gradient-to-br from-blue-50 to-white flex items-center justify-center">
			<div className="text-center">
				<div className="w-16 h-16 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
				<p className="text-gray-600">Loading...</p>
			</div>
		</div>
	)
}

/*"use client"

import { useState } from "react"
// --- FIX: Update import paths to the new component locations ---
import LoginPage from "../components/login-page"
import SignupPage from "../components/signup-page"
// --- END FIX ---
import EmailVerificationPage from "../components/email-verification"
import DashboardPage from "../app/dashboard/page" // Assuming this is your authenticated page

// This is the main controller for your application's authentication state.
// It decides which of the child components (Login, Signup, Verify) to show.

type AuthState = "login" | "signup" | "verifying" | "authenticated"

interface VerificationData {
	userId: number;
	email: string;
	password: string; // Used for auto-login after verification
	name: string;
	phone: string;
}

export default function Page() {
	const [authState, setAuthState] = useState<AuthState>("login")
	const [userData, setUserData] = useState<any>(null) // Store user data after login
	const [verificationData, setVerificationData] = useState<VerificationData | null>(null)

	// This is called from SignupPage when the /register API is successful
	const handleSignupSuccess = (data: VerificationData) => {
		setVerificationData(data)
		setAuthState("verifying") // Switch to the verification page
	}

	// This is called from LoginPage when the /login API is successful
	const handleLoginSuccess = (user: any) => {
		setUserData(user)
		setAuthState("authenticated")
	}

	// This is called from EmailVerificationPage when the /verify API is successful
	const handleVerificationComplete = () => {
		setAuthState("authenticated") // User is now fully logged in
		setVerificationData(null) // Clear temp data
	}

	// Render components based on the current auth state

	if (authState === "authenticated") {
		// You would render your main application here
		return <DashboardPage />
	}

	if (authState === "verifying" && verificationData) {
		// Render the EmailVerificationPage
		return (
			<EmailVerificationPage
				email={verificationData.email}
				userId={verificationData.userId}
				password={verificationData.password}
				onVerificationComplete={handleVerificationComplete}
			/>
		)
	}

	if (authState === "signup") {
		// Render the SignupPage
		return (
			<SignupPage
				onSignupSuccess={handleSignupSuccess}
				onSwitchToLogin={() => setAuthState("login")}
			/>
		)
	}

	// The default state is "login"
	return (
		<LoginPage
			onLoginSuccess={handleLoginSuccess}
			onSwitchToSignup={() => setAuthState("signup")}
		/>
	)
}
*/