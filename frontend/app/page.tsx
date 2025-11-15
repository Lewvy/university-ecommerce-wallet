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