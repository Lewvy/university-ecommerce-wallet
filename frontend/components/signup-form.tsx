"use client"

import type React from "react"
import { useState } from "react"

interface SignupFormProps {
	onSubmit: (data: { email: string, username: string, phone: string, password: string }) => void
}

export default function SignupForm({ onSubmit }: SignupFormProps) {
	const [formData, setFormData] = useState({
		username: "",
		email: "",
		password: "",
		phone: "",
	})

	const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		const { name, value } = e.target
		setFormData((prev) => ({ ...prev, [name]: value }))
	}

	const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()

		onSubmit(formData)
	}

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

					<h2 className="text-2xl font-bold text-gray-900 mb-2">Create Account</h2>
					<p className="text-gray-600 mb-6">Join your university marketplace</p>

					<form onSubmit={handleSubmit} className="space-y-4">
						<div>
							<label className="block text-sm font-medium text-gray-700 mb-2">Username</label>
							<input
								type="text"
								name="username"
								value={formData.username}
								onChange={handleChange}
								placeholder="Enter your username"
								required
								className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
							/>
						</div>
						<div>
							<label className="block text-sm font-medium text-gray-700 mb-2">University Email</label>
							<input
								type="email"
								name="email"
								value={formData.email}
								onChange={handleChange}
								placeholder="your.name@gmail.com"
								required
								className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
							/>
						</div>
						<div>
							<label className="block text-sm font-medium text-gray-700 mb-2">Password</label>
							<input
								type="password"
								name="password"
								value={formData.password}
								onChange={handleChange}
								placeholder="Enter your password"
								required
								className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
							/>
						</div>
						<div>
							<label className="block text-sm font-medium text-gray-700 mb-2">Phone Number</label>
							<input
								type="tel"
								name="phone"
								value={formData.phone}
								onChange={handleChange}
								placeholder="+91 9876543210"
								required
								className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
							/>
						</div>

						<button
							type="submit"
							className="w-full bg-gradient-to-r from-blue-600 to-blue-500 text-white font-semibold py-2 rounded-lg hover:shadow-lg transition-shadow"
						>
							Create Account
						</button>
					</form>

					<p className="text-center text-gray-600 mt-4">
						Already have an account? <span className="text-blue-600 font-semibold cursor-pointer">Login</span>
					</p>
				</div>
			</div>
		</div>
	)
}
