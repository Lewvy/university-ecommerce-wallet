"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"

export default function Page() {
	const router = useRouter()

	useEffect(() => {
		const user = localStorage.getItem('user')
		const token = localStorage.getItem('access_token')

		if (user && token) {
			router.push('/dashboard')
		} else {
			router.push('/login')
		}
	}, [router])

	return (
		<div className="min-h-screen bg-gradient-to-br from-blue-50 to-white flex items-center justify-center">
			<div className="text-center">
				<div className="w-16 h-16 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
				<p className="text-gray-600">Loading...</p>
			</div>
		</div>
	)
}
