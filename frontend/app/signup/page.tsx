"use client"

import { useState } from "react"
import EmailVerificationPage from "./email-verification"
import SignupForm from "./signup-form"

interface SignupPageProps {
	onSignupSuccess: (username: string, email: string, phone: string) => void
}

interface FormData {
	email: string;
	username: string;
	phone: string;
	password: string;
}

export default function SignupPage({ onSignupSuccess }: SignupPageProps) {
	const [step, setStep] = useState<"form" | "verification" | "success">("form")
	const [userEmail, setUserEmail] = useState("")
	const [username, setUsername] = useState("")
	const [userPhone, setUserPhone] = useState("")
	const [password, setPassword] = useState("")

	const handleFormSubmit = async (data: FormData) => {
		try {
			const response = await fetch("http://localhost:8088/register", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(data),
			});

			const responseData = await response.json();

			if (response.ok) {
				console.log("Registration Success:", responseData);

				setUserEmail(data.email)
				setUsername(data.username)
				setUserPhone(data.phone)
				setPassword(data.password)

				setStep("verification")
			} else {
				console.error("Registration Failed:", responseData);

				let errorMessage = "Registration failed. Please try again.";
				if (responseData.message) {
					errorMessage = responseData.message;
				} else if (responseData.fields) {
					errorMessage = `Validation Error: ${Object.values(responseData.fields).join(', ')}`;
				}

				alert(errorMessage);
			}
		} catch (error) {
			console.error("Network Error:", error);
			alert("A network error occurred. Check your server connection.");
		}
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
