"use client"

import { useState } from "react"
import { useRouter } from 'next/navigation'

function decodeJWT(token: string) {
    try {
        const base64 = token.split(".")[1];
        const decoded = JSON.parse(atob(base64));
        return decoded;
    } catch (e) {
        console.error("JWT decode error:", e);
        return null;
    }
}

export default function LoginPage() {
	const [formData, setFormData] = useState({ email: "", password: "" })
	const [isLoading, setIsLoading] = useState(false)
	const [error, setError] = useState<string | null>(null)

	const router = useRouter()

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
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(formData),
			})

			const data = await response.json()

			if (response.ok) {
				console.log("Login response:", data)

				// Save tokens
				localStorage.setItem("access_token", data.access_token)
				localStorage.setItem("refresh_token", data.refresh_token)

				// Decode JWT to get user_id
				const jwtPayload = decodeJWT(data.access_token)
				console.log("JWT payload:", jwtPayload)

				// IMPORTANT: Use data from backend response, not JWT
				const userInfo = {
					id: jwtPayload?.user_id || null,
					username: data.name || "User",  // Backend returns "name"
					email: formData.email,           // Use login email
					phone: data.phone || "",         // Backend returns "phone"
				}

				console.log("Storing user info:", userInfo)
				localStorage.setItem("user", JSON.stringify(userInfo))

				router.push("/dashboard")
			} else {
				setError(data.message || "Login failed.")
			}
		} catch (err) {
			setError("Could not reach server.")
		} finally {
			setIsLoading(false)
		}
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

					<h2 className="text-2xl font-bold mb-4">Welcome Back</h2>

					<form onSubmit={handleSubmit} className="space-y-4">
						<div>
							<label className="block text-sm font-medium text-gray-700 mb-2">Email</label>
							<input 
								name="email" 
								type="email"
								value={formData.email} 
								onChange={handleChange} 
								className="w-full p-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none" 
								placeholder="Email"
								required
							/>
						</div>

						<div>
							<label className="block text-sm font-medium text-gray-700 mb-2">Password</label>
							<input 
								name="password" 
								type="password" 
								value={formData.password} 
								onChange={handleChange} 
								className="w-full p-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none" 
								placeholder="Password"
								required
							/>
						</div>

						{error && <p className="text-red-600 text-center text-sm">{error}</p>}

						<button 
							type="submit" 
							disabled={isLoading} 
							className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition disabled:opacity-50"
						>
							{isLoading ? "Logging in..." : "Login"}
						</button>
					</form>

					<p className="text-center text-gray-600 mt-4 text-sm">
						Don't have an account? <a href="/signup" className="text-blue-600 font-semibold hover:underline">Sign Up</a>
					</p>
				</div>
			</div>
		</div>
	)
}

/*"use client"

import { useState } from "react"
import { useRouter } from 'next/navigation'

export default function LoginPage() {
	const [formData, setFormData] = useState({
		email: "",
		password: "",
	})
	const [isLoading, setIsLoading] = useState(false)
	const [error, setError] = useState<string | null>(null)

	const router = useRouter()

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
			console.log("Attempting login with:", formData.email)

			const response = await fetch("http://localhost:8088/login", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(formData),
			})

			const data = await response.json()
			console.log("Login response:", data)

			if (response.ok) {
				console.log("Login successful!")
				
				// Store authentication tokens
				if (data.access_token) {
					localStorage.setItem('access_token', data.access_token)
				}
				if (data.refresh_token) {
					localStorage.setItem('refresh_token', data.refresh_token)
				}
				
				// Store user info - normalize the data structure
				const userInfo = {
					username: data.user?.username || data.user?.name || data.username || data.name || "User",
					email: data.user?.email || data.email || formData.email,
					phone: data.user?.phone || data.phone || "",
					id: data.user?.id || data.id
				}
				
				console.log("Storing user info:", userInfo)
				localStorage.setItem('user', JSON.stringify(userInfo))
				
				// Force a small delay to ensure storage completes
				await new Promise(resolve => setTimeout(resolve, 100))
				
				// Redirect to dashboard
				console.log("Redirecting to dashboard...")
				router.push('/dashboard')
			} else {
				console.error("Login failed:", data)
				setError(data.message || "Login failed. Please check your credentials.")
			}
		} catch (err) {
			console.error("Network or server error:", err)
			setError("Could not connect to the server. Please try again later.")
		} finally {
			setIsLoading(false)
		}
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
						Don't have an account? <a href="/signup" className="text-blue-600 font-semibold cursor-pointer hover:underline">Sign Up</a>
					</p>
				</div>
			</div>
		</div>
	)
}*/