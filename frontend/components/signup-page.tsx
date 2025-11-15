"use client"

import { useState } from "react"
import SignupForm from "./signup-form"

// Define the props this component accepts from app/page.tsx
interface SignupPageProps {
	onSignupSuccess: (data: { userId: number; email: string; password: string; name: string; phone: string }) => void
	onSwitchToLogin: () => void;
}

interface FormData {
	email: string;
	name: string;
	phone: string;
	password: string;
}

export default function SignupPage({ onSignupSuccess, onSwitchToLogin }: SignupPageProps) {
	// This component now manages the API call state
	const [isLoading, setIsLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	// This state is crucial for the fix.
	const [showLoginButton, setShowLoginButton] = useState(false);

	const handleFormSubmit = async (data: FormData) => {
		setIsLoading(true);
		setError(null);
		setShowLoginButton(false);

		try {
			console.log("Sending registration request with data:", data)

			const response = await fetch("http://localhost:8088/register", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(data),
			});

			const responseData = await response.json();
			console.log("Registration response:", response.status, responseData)

			if (response.ok) {
				console.log("Registration Success:", responseData);
				const userId = responseData.id || (responseData.user && responseData.user.id);

				if (!userId) {
					console.warn("No user ID returned from registration");
					setError("Registration succeeded but no User ID was returned.");
					setIsLoading(false);
					return;
				}

				// --- SUCCESS ---
				// Call the prop to notify app/page.tsx to switch to the verification step
				onSignupSuccess({
					userId: userId,
					email: data.email,
					name: data.name,
					phone: data.phone,
					password: data.password // Pass password for auto-login
				});

			} else {
				// --- ERROR HANDLING ---
				console.error("Registration Failed:", responseData);
				let errorMessage = responseData.error || responseData.message || "Registration failed. Please try again.";

				// --- THIS IS THE FIX ---
				// Check for the specific error from your backend
				if (errorMessage.toLowerCase().includes("user already exists")) {
					errorMessage = "This email is already registered. Please log in.";
					setShowLoginButton(true); // Tell the UI to show the "Go to Login" button
				}
				// --- END OF FIX ---

				setError(errorMessage);
			}
		} catch (error) {
			console.error("Network Error:", error);
			setError("A network error occurred. Check your server connection.");
		} finally {
			setIsLoading(false);
		}
	}

	// This page component now just renders the form, passing down
	// the state and the prop functions.
	return (
		<SignupForm
			onSubmit={handleFormSubmit}
			isLoading={isLoading}
			error={error}
			showLoginButton={showLoginButton}
			onSwitchToLogin={onSwitchToLogin} // Pass the prop down
		/>
	)
}
