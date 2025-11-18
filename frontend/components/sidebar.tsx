"use client"

import { useState } from "react"

interface SidebarProps {
	currentView: string
	onViewChange: (view: any) => void
}

const menuItems = [
	{ id: "home", label: "Buy Items", icon: "🛍️" },
	{ id: "sell", label: "Sell an Item", icon: "📤" },
	{ id: "profile", label: "My Profile", icon: "👤" },
	{ id: "categories", label: "Categories", icon: "📂" },
	{ id: "buy", label: "My Cart", icon: "🛒" },
	{ id: "wallet", label: "My Wallet", icon: "💳" }, ]
  
export default function Sidebar({ currentView, onViewChange }: SidebarProps) {
	const [isOpen, setIsOpen] = useState(false)

	return (
		<>
			<button
				onClick={() => setIsOpen(!isOpen)}
				className="lg:hidden fixed top-4 left-4 z-50 p-2 bg-blue-600 text-white rounded-lg"
			>
				<svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
				</svg>
			</button>

			{isOpen && (
				<div className="lg:hidden fixed inset-0 bg-black bg-opacity-50 z-30" onClick={() => setIsOpen(false)} />
			)}

			<aside
				className={`fixed lg:static w-64 bg-white border-r border-gray-200 overflow-y-auto h-screen z-40 transition-transform ${
					isOpen ? "translate-x-0" : "-translate-x-full lg:translate-x-0"
				}`}
			>
				<div className="p-6 border-b border-gray-200">
					<div className="flex items-center space-x-3">
						<div className="w-10 h-10 bg-gradient-to-br from-blue-600 to-blue-400 rounded-lg flex items-center justify-center">
							<span className="text-white font-bold">U</span>
						</div>
						<h1 className="text-xl font-bold text-gray-900">Unimart</h1>
					</div>
				</div>

				<nav className="p-4 space-y-2">
					{menuItems.map((item) => (
						<button
							key={item.id}
							onClick={() => {
								onViewChange(item.id)
								setIsOpen(false)
							}}
							className={`w-full flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors text-left ${
								currentView === item.id ? "bg-blue-50 text-blue-600 font-semibold" : "text-gray-700 hover:bg-gray-50"
							}`}
						>
							<span className="text-xl">{item.icon}</span>
							<span>{item.label}</span>
						</button>
					))}
				</nav>

				<div className="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-200 bg-white w-64">
					<button
						onClick={() => {
							onViewChange("logout")
							setIsOpen(false)
						}}
						className="w-full px-4 py-2 text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors font-medium text-sm"
					>
						Logout
					</button>
				</div>
			</aside>
		</>
	)
}