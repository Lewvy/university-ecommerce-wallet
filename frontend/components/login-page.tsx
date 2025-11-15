"use client"

import { useState } from "react"
// We remove useRouter because the parent (app/page.tsx) handles state switching
// import { useRouter } from 'next/navigation'

// Define the props this component now accepts from its parent (app/page.tsx)
interface LoginPageProps {
	onLoginSuccess: (user: any) => void;
	onSwitchToSignup: () => void;
}

export default function LoginPage({ onLoginSuccess, onSwitchToSignup }: LoginPageProps) {
	const [formData, setFormData] = useState({
		email: "",
		password: "",
	})
	const [isLoading, setIsLoading] = useState(false)
	const [error, setError] = useState<string | null>(null)

	// const router = useRouter() // No longer needed

	const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		const { name, value } = e.target
		setFormData((prev) => ({ ...prev, [name]: value }))
		if (error) setError(null)
	}

	const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()
		setError(null)
		setIsLoading(true)

		try {
			const response = await fetch("http://localhost:8088/login", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(formData),
			})

			const data = await response.json()

			if (response.ok) {
				console.log("Login successful:", data)

				// Store authentication tokens
				if (data.access_token) {
					sessionStorage.setItem('access_token', data.access_token)
					console.log("Access token stored")
				}
				if (data.refresh_token) {
					sessionStorage.setItem('refresh_token', data.refresh_token)
					console.log("Refresh token stored")
				}

				// Store user info if provided
				const user = data.user || { email: formData.email }; // Create a fallback user object
				if (data.user) {
					sessionStorage.setItem('user', JSON.stringify(data.user))
				}

				// --- FIX ---
				// Instead of router.push, we call the onLoginSuccess prop
				// to tell app/page.tsx that we are authenticated.
				onLoginSuccess(user);

			} else {
				console.error("Login failed:", data)
				setError(data.error || "Login failed. Please check your credentials.")
			}
		} catch (err) {
			console.error("Network or server error:", err)
			setError("Could not connect to the server. Please try again later.")
		} finally {
			setIsLoading(false)
		}
	}

	// The rest of your JSX remains the same, but we update the "Sign Up" link
	return (
		<div className="min-h-screen bg-gradient-to-br from-blue-50 to-white flex items-center justify-center p-4">
			<div className="w-full max-w-md">
				<div className="bg-white rounded-lg shadow-lg p-8">
					<div className="flex items-center justify-center mb-8">
						<div className="w-12 h-12 bg-gradient-to-br from-blue-600 to-blue-400 rounded-lg flex items-center justify-center">
							<span className="text-white font-bold text-xl">U</span>
						</div>
						<h1 className="text-2xl font-bold text-gray-900 ml-3">Unimart</h1>
					</div>

					<h2 className="text-2xl font-bold text-gray-900 mb-2">Welcome Back</h2>
					<p className="text-gray-600 mb-6">Sign in to your account to continue.</p>

					<form onSubmit={handleSubmit} className="space-y-4">
						<div>
							<label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">University Email</label>
							<input
								id="email"
								type="email"
								name="email"
								value={formData.email}
								onChange={handleChange}
								placeholder="your.name@example.com"
								required
								className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
							/>
						</div>

						<div>
							<label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">Password</label>
							<input
								id="password"
								type="password"
								name="password"
								value={formData.password}
								onChange={handleChange}
								placeholder="Enter your password"
								required
								className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
							/>
						</div>

						{error && (
							<p className="text-sm text-red-600 font-medium text-center">{error}</p>
						)}

						<button
							type="submit"
							disabled={isLoading}
							className="w-full bg-gradient-to-r from-blue-600 to-blue-500 text-white font-semibold py-2 rounded-lg hover:shadow-lg transition-shadow disabled:opacity-50 disabled:cursor-not-allowed"
						>
							{isLoading ? 'Signing In...' : 'Sign In'}
						</button>
					</form>

					<p className="text-center text-gray-600 mt-4">
						Don't have an account?{" "}
						{/* --- FIX --- 
                           Instead of an <a href>, we use a button that calls the prop
                           to tell app/page.tsx to switch the state to "signup".
                        */}
						<button
							onClick={onSwitchToSignup}
							className="text-blue-600 font-semibold cursor-pointer hover:underline"
						>
							Sign Up
						</button>
					</p>
				</div>
			</div>
		</div>
	)
}
