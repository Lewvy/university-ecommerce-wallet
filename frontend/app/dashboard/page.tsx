"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import Header from "../../components/header"
import Sidebar from "../../components/sidebar"
import MarketplaceView from "../../components/marketplace-view"
import SellItemForm from "../../components/sell-item-form"
import ProfilePage from "../../components/profile-page"
import CartPage from "../../components/cart-page"
import CategoriesView from "../../components/categories-view"
import WalletPage from "../../components/wallet-page"

type CurrentView = "home" | "buy" | "sell" | "profile" | "categories" | "wallet" | "logout"

export default function DashboardPage() {
	const router = useRouter()
	const [userData, setUserData] = useState<any>(null)
	const [isLoading, setIsLoading] = useState(true)

	// Check authentication on mount
	useEffect(() => {
		console.log("Dashboard: Checking authentication...")
		
		const storedUser = localStorage.getItem("user")
		const token = localStorage.getItem("access_token")

		console.log("Stored user:", storedUser)
		console.log("Has token:", !!token)

		if (!storedUser || !token) {
			console.log("No user or token found, redirecting to login...")
			router.push("/login")
			return
		}

		try {
			const user = JSON.parse(storedUser)
			console.log("User data loaded:", user)
			setUserData(user)
		} catch (error) {
			console.error("Error parsing user data:", error)
			router.push("/login")
		} finally {
			setIsLoading(false)
		}
	}, [router])

	// Dashboard state
	const [currentView, setCurrentView] = useState<CurrentView>("home")
	const [searchQuery, setSearchQuery] = useState("")
	const [selectedCategory, setSelectedCategory] = useState<string | null>(null)
	const [cartItems, setCartItems] = useState<any[]>([])
	const [postedItems, setPostedItems] = useState<any[]>([])

	const handleAddToCart = (item: any) => setCartItems([...cartItems, item])
	const handlePostItem = (item: any) => setPostedItems([...postedItems, item])
	const handleCategorySelect = (category: string) => {
		setSelectedCategory(category)
		setCurrentView("home")
	}

	const handleLogout = () => {
		console.log("Logging out...")
		localStorage.clear()
		router.push("/login")
	}

	// Handle logout view
	useEffect(() => {
		if (currentView === "logout") {
			handleLogout()
		}
	}, [currentView])

	// Show loading spinner while checking auth
	if (isLoading) {
		return (
			<div className="h-screen flex items-center justify-center bg-gray-50">
				<div className="text-center">
					<div className="w-16 h-16 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
					<p className="text-gray-600">Loading dashboard...</p>
				</div>
			</div>
		)
	}

	// If no user data after loading, don't render
	if (!userData) {
		return null
	}

	return (
		<div className="flex flex-col lg:flex-row h-screen bg-gray-50">
			<Sidebar currentView={currentView} onViewChange={setCurrentView} />

			<div className="flex-1 flex flex-col overflow-hidden">
				{currentView === "home" && (
					<Header
						username={userData.username}
						searchQuery={searchQuery}
						onSearchChange={setSearchQuery}
						selectedCategory={selectedCategory}
						onCategoryChange={setSelectedCategory}
					/>
				)}

				<main className="flex-1 overflow-auto p-4 md:p-6">
					{currentView === "home" && (
						<MarketplaceView
							searchQuery={searchQuery}
							selectedCategory={selectedCategory}
							postedItems={postedItems}
							onAddToCart={handleAddToCart}
						/>
					)}

					{currentView === "buy" && <CartPage items={cartItems} />}
					{currentView === "sell" && <SellItemForm onSubmit={handlePostItem} userData={userData} />}
					{currentView === "profile" && <ProfilePage userData={userData} postedItems={postedItems} />}
					{currentView === "categories" && (
						<CategoriesView onCategorySelect={handleCategorySelect} postedItems={postedItems} />
					)}
					{currentView === "wallet" && <WalletPage userData={userData} />}
				</main>
			</div>
		</div>
	)
}