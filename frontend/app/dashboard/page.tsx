"use client"

import { useState } from "react"
import Sidebar from "./sidebar"
import Header from "./header"
import MarketplaceView from "./marketplace-view"
import SellItemForm from "./sell-item-form"
import ProfilePage from "./profile-page"
import CartPage from "./cart-page"
import CategoriesView from "./categories-view"
import LoginPage from "./login-page"

interface HomePageProps {
	userData: {
		username: string
		email: string
		phone: string
	}
}

type CurrentView = "home" | "buy" | "sell" | "profile" | "categories" | "logout"

export default function HomePage({ userData }: HomePageProps) {
	const [currentView, setCurrentView] = useState<CurrentView>("home")
	const [searchQuery, setSearchQuery] = useState("")
	const [selectedCategory, setSelectedCategory] = useState<string | null>(null)
	const [cartItems, setCartItems] = useState<any[]>([])
	const [postedItems, setPostedItems] = useState<any[]>([])

	const handleAddToCart = (item: any) => {
		setCartItems([...cartItems, item])
	}

	const handlePostItem = (item: any) => {
		setPostedItems([...postedItems, item])
	}

	const handleCategorySelect = (category: string) => {
		setSelectedCategory(category)
		setCurrentView("home")
	}

	if (currentView === "logout") {
		return <LoginPage />
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
				</main>
			</div>
		</div>
	)
}
