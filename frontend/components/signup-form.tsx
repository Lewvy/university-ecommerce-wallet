"use client"

import { useState } from "react"

// This is the actual UI for the signup form.
// It is a "dumb" component that just displays state and reports events.

interface SignupFormProps {
	onSubmit: (data: { email: string; name: string; phone: string; password: string }) => void;
	isLoading: boolean;
	error: string | null;
	showLoginButton: boolean; // <-- Prop for the fix
	onSwitchToLogin: () => void; // <-- Prop for the fix
}

export default function SignupForm({ onSubmit, isLoading, error, showLoginButton, onSwitchToLogin }: SignupFormProps) {
	// Local state for form fields
	const [name, setName] = useState("")
	const [email, setEmail] = useState("")
	const [phone, setPhone] = useState("")
	const [password, setPassword] = useState("")
	const [confirmPassword, setConfirmPassword] = useState("")
	const [localError, setLocalError] = useState<string | null>(null);

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		setLocalError(null);

		// Local form validation
		if (password !== confirmPassword) {
			setLocalError("Passwords do not match.");
			return;
		}
		if (password.length < 8) {
			setLocalError("Password must be at least 8 characters.");
			return;
		}
		// If validation passes, call the onSubmit prop to trigger the API call
		onSubmit({ name, email, phone, password });
	}

	return (
		<div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
			<div className="max-w-md w-full bg-white rounded-lg shadow-lg p-8">
				<h2 className="text-3xl font-bold text-center text-gray-900 mb-6">Create your account</h2>
				<form onSubmit={handleSubmit} className="space-y-4">
					<div>
						<label className="block text-sm font-medium text-gray-700">Full Name</label>
						<input
							type="text"
							value={name}
							onChange={(e) => setName(e.target.value)}
							required
							className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
						/>
					</div>
					<div>
						<label className="block text-sm font-medium text-gray-700">Email</label>
						<input
							type="email"
							value={email}
							onChange={(e) => setEmail(e.target.value)}
							required
							className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
						/>
					</div>
					<div>
						<label className="block text-sm font-medium text-gray-700">Phone</label>
						<input
							type="tel"
							value={phone}
							onChange={(e) => setPhone(e.target.value)}
							required
							className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
						/>
					</div>
					<div>
						<label className="block text-sm font-medium text-gray-700">Password</label>
						<input
							type="password"
							value={password}
							onChange={(e) => setPassword(e.target.value)}
							required
							className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
						/>
					</div>
					<div>
						<label className="block text-sm font-medium text-gray-700">Confirm Password</label>
						<input
							type="password"
							value={confirmPassword}
							onChange={(e) => setConfirmPassword(e.target.value)}
							required
							className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
						/>
					</div>

					{/* --- ERROR DISPLAY --- */}
					{(error || localError) && (
						<div className="bg-red-50 border border-red-200 rounded-lg p-3 text-center">
							<p className="text-sm text-red-600 font-medium">{error || localError}</p>

							{/* --- THE FIX: Conditional Login Button --- */}
							{showLoginButton && (
								<button
									type="button"
									onClick={onSwitchToLogin} // This calls the prop
									className="mt-2 w-full bg-green-500 text-white font-semibold py-2 rounded-lg hover:bg-green-600"
								>
									Go to Login
								</button>
							)}
						</div>
					)}

					<button
						type="submit"
						disabled={isLoading}
						className="w-full bg-blue-600 text-white font-semibold py-3 rounded-lg hover:shadow-lg transition-shadow disabled:opacity-50"
					>
						{isLoading ? "Creating Account..." : "Create Account"}
					</button>
				</form>

				<p className="text-center text-sm text-gray-600 mt-6">
					Already have an account?{" "}
					<button
						onClick={onSwitchToLogin} // This calls the prop
						className="font-semibold text-blue-600 hover:underline"
					>
						Log in
					</button>
				</p>
			</div>
		</div>
	)
}
