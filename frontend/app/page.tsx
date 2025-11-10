"use client"

import { useState } from "react"
// import SellItemForm from "./sell-item-form"
//
import ProfilePage from "../components/profile-page"
//
// import CartPage from "./cart-page"
//
// import CategoriesView from "../app./categories-view"

// Line 11
import LoginPage from "../app/login/page" // Relative path (depends on file location)

export default function Page() {
	const [isAuthenticated, setIsAuthenticated] = useState(false)
	const [authView, setAuthView] = useState<'login' | 'signup'>('login')

	const [userData, setUserData] = useState({
		username: "",
		email: "",
		phone: "",
	})

	const handleSignupSuccess = (name: string, email: string, phone: string) => {
		setUserData({ username: name, email, phone })
		setIsAuthenticated(true)
	}

	const handleLoginSuccess = (name: string, email: string, phone: string) => {
		setUserData({ username: name, email, phone })
		setIsAuthenticated(true)
	}

	if (isAuthenticated) {
		return <HomePage userData={userData} />
	}

	if (authView === 'signup') {
		return (
			<SignupPage
				onSignupSuccess={handleSignupSuccess}
				onSwitchToLogin={() => setAuthView('login')}
			/>
		)
	}

	return (
		<LoginPage
			onLoginSuccess={handleLoginSuccess}
			onSwitchToSignup={() => setAuthView('signup')}
		/>
	)
}
