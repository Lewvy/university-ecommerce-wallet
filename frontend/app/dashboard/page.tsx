"use client"

// We're assuming this is your main authenticated page (HomePage/Dashboard)
// This file had all the broken imports.

import { useState } from "react"
// --- FIX: Corrected all import paths ---
import Sidebar from "../../components/sidebar"
import Header from "../../components/header"
import MarketplaceView from "../../components/marketplace-view"
import SellItemForm from "../../components/sell-item-form"
import ProfilePage from "../../components/profile-page"
import CartPage from "../../components/cart-page"
import CategoriesView from "../../components/categories-view"
// This login page import was also broken
import LoginPage from "../../components/login-page"
// --- END FIX ---

// Define a type for your user data
interface UserData {
	name: string;
	username?: string; // Add optional username
	email: string;
	phone: string;
	postedItems: any[]; // <-- Add postedItems to the interface
}

interface HomePageProps {
	userData: UserData | null; // Allow userData to be null initially
}

// Re-creating the component structure based on the imports
export default function DashboardPage({ userData }: HomePageProps) {
	const [activeView, setActiveView] = useState("marketplace")

	const renderView = () => {
		switch (activeView) {
			case "marketplace":
				// --- FIX ---
				// Pass the 'postedItems' from the userData prop
				// down to the MarketplaceView component.
				// We use `userData?.postedItems || []` as a safeguard.
				return <MarketplaceView postedItems={userData?.postedItems || []} onAddToCart={() => { }} />
			case "sell":
				return <SellItemForm />
			case "profile":
				return <ProfilePage />
			case "cart":
				return <CartPage />
			case "categories":
				return <CategoriesView />
			default:
				// --- FIX ---
				// Also apply the fix to the default case
				return <MarketlaceView postedItems={userData?.postedItems || []} onAddToCart={() => { }} />
		}
	}

	// This is just a sample layout. You can replace this with your actual JSX.
	return (
		<div className="flex h-screen bg-gray-100">
			<Sidebar onNavigate={setActiveView} />
			<div className="flex-1 flex flex-col">
				{/* Fix for name vs username mismatch */}
				<Header username={userData?.username || userData?.name || "User"} />
				<main className="flex-1 p-6 overflow-y-auto">
					{renderView()}
				</main>
			</div>
		</div>
	)
}
