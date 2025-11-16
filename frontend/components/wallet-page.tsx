"use client"

import { useState, useEffect } from "react"

interface WalletData {
	user_id: number
	balance: number
	lifetime_spent: number
	lifetime_earned: number
	created_at: string
	updated_at: string
}

interface WalletPageProps {
	userData: {
		username: string
		email: string
		phone: string
		id?: number
	}
}

// Declare Razorpay on window
declare global {
	interface Window {
		Razorpay: any
	}
}

export default function WalletPage({ userData }: WalletPageProps) {
	const [walletData, setWalletData] = useState<WalletData | null>(null)
	const [isLoading, setIsLoading] = useState(true)
	const [showAddMoney, setShowAddMoney] = useState(false)
	const [showTransfer, setShowTransfer] = useState(false)
	const [addAmount, setAddAmount] = useState("")
	const [recipientId, setRecipientId] = useState("")
	const [transferAmount, setTransferAmount] = useState("")
	const [paymentLoading, setPaymentLoading] = useState(false)

	useEffect(() => {
		fetchWalletData()
		loadRazorpayScript()
	}, [])

	const loadRazorpayScript = () => {
		return new Promise((resolve) => {
			// Check if script already loaded
			if (window.Razorpay) {
				resolve(true)
				return
			}

			const script = document.createElement('script')
			script.src = 'https://checkout.razorpay.com/v1/checkout.js'
			script.onload = () => resolve(true)
			script.onerror = () => resolve(false)
			document.body.appendChild(script)
		})
	}

	const fetchWalletData = async () => {
		setIsLoading(true)
		try {
			const token = localStorage.getItem("access_token")

			console.log("Fetching wallet - Token found:", token ? "Yes" : "No")

			if (!token) {
				console.error("No access token found")
				alert("You are not logged in. Please login first.")
				setIsLoading(false)
				return
			}

			const response = await fetch("http://localhost:8088/wallet/balance", {
				headers: {
					"Authorization": `Bearer ${token}`,
					"Content-Type": "application/json"
				}
			})

			if (response.ok) {
				const data = await response.json()
				console.log("Wallet data:", data)
				setWalletData(data)
			} else {
				const error = await response.json()
				console.error("Failed to fetch wallet:", error)
				alert(error.error || "Failed to load wallet")
			}
		} catch (error) {
			console.error("Error fetching wallet data:", error)
			alert("Network error. Please try again.")
		} finally {
			setIsLoading(false)
		}
	}

	const handleAddMoney = async () => {
		const amountInRupees = parseFloat(addAmount)

		if (!amountInRupees || amountInRupees <= 0) {
			alert("Please enter a valid amount")
			return
		}

		if (amountInRupees < 1) {
			alert("Minimum amount is ₹1")
			return
		}

		if (amountInRupees > 100000) {
			alert("Maximum amount is ₹1,00,000")
			return
		}

		const amountInPaise = Math.round(amountInRupees * 100)

		setPaymentLoading(true)

		try {
			if (!window.Razorpay) {
				const loaded = await loadRazorpayScript()
				if (!loaded) {
					alert("Failed to load payment gateway. Please refresh and try again.")
					setPaymentLoading(false)
					return
				}
			}

			const token = localStorage.getItem("access_token")

			console.log("Token found:", token ? "Yes" : "No")
			console.log("Token length:", token?.length || 0)

			if (!token) {
				alert("You are not logged in. Please login first.")
				setPaymentLoading(false)
				return
			}

			const response = await fetch("http://localhost:8088/wallet/create-topup-order", {
				method: "POST",
				headers: {
					"Authorization": `Bearer ${token}`,
					"Content-Type": "application/json"
				},
				body: JSON.stringify({ amount: amountInPaise })
			})

			console.log("Response status:", response.status)

			if (!response.ok) {
				throw new Error("Failed to create payment order")
			}

			const orderData = await response.json()

			const options = {
				key: orderData.key_id,
				amount: orderData.amount,
				currency: orderData.currency,
				name: "Your Store",
				description: "Wallet Top-up",
				order_id: orderData.order_id,
				handler: async function(response: any) {
					console.log("Payment successful:", response)
					alert("Payment successful! Your wallet will be credited shortly.")
					setAddAmount("")
					setShowAddMoney(false)

					setTimeout(() => {
						fetchWalletData()
					}, 2000)
				},
				prefill: {
					name: userData.username || "User",
					email: userData.email || "",
					contact: userData.phone || ""
				},
				theme: {
					color: "#3b82f6"
				},
				modal: {
					ondismiss: function() {
						console.log("Payment cancelled by user")
						setPaymentLoading(false)
					}
				}
			}

			const razorpay = new window.Razorpay(options)

			razorpay.on('payment.failed', function(response: any) {
				alert(`Payment failed: ${response.error.description}`)
				setPaymentLoading(false)
			})

			razorpay.open()
			setPaymentLoading(false)

		} catch (error) {
			console.error("Error initiating payment:", error)
			alert("Failed to initiate payment. Please try again.")
			setPaymentLoading(false)
		}
	}

	const handleTransfer = async () => {
		const amount = parseInt(transferAmount)
		const recipient = parseInt(recipientId)

		if (!amount || amount <= 0) {
			alert("Please enter a valid amount")
			return
		}

		if (!recipient || recipient <= 0) {
			alert("Please enter a valid recipient user ID")
			return
		}

		try {
			let token = sessionStorage.getItem("access_token")
			if (!token) {
				token = localStorage.getItem("access_token")
			}

			if (!token) {
				alert("You are not logged in. Please login first.")
				return
			}

			const response = await fetch("http://localhost:8088/wallet/transfer", {
				method: "POST",
				headers: {
					"Authorization": `Bearer ${token}`,
					"Content-Type": "application/json"
				},
				body: JSON.stringify({
					recipient_user_id: recipient,
					amount: amount
				})
			})

			if (response.ok) {
				const data = await response.json()
				console.log("Transfer successful:", data)
				setTransferAmount("")
				setRecipientId("")
				setShowTransfer(false)
				alert(data.message || "Transfer successful!")
				fetchWalletData()
			} else {
				const error = await response.json()
				alert(error.error || "Failed to transfer money")
			}
		} catch (error) {
			console.error("Error transferring money:", error)
			alert("Network error. Please try again.")
		}
	}

	const quickAmounts = [5, 10, 20, 50, 100, 500]

	if (isLoading) {
		return (
			<div className="flex items-center justify-center h-full">
				<div className="text-center">
					<div className="w-12 h-12 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
					<p className="text-gray-600">Loading wallet...</p>
				</div>
			</div>
		)
	}

	if (!walletData) {
		return (
			<div className="flex items-center justify-center h-full">
				<div className="text-center">
					<p className="text-gray-600">Failed to load wallet data</p>
					<button
						onClick={fetchWalletData}
						className="mt-4 bg-blue-600 text-white px-6 py-2 rounded-lg font-semibold hover:bg-blue-700"
					>
						Retry
					</button>
				</div>
			</div>
		)
	}

	return (
		<div className="max-w-4xl mx-auto">
			<h1 className="text-3xl font-bold text-gray-900 mb-6">My Wallet</h1>

			{/* Balance Card */}
			<div className="bg-gradient-to-br from-blue-600 to-blue-500 rounded-xl p-8 mb-6 shadow-lg">
				<div className="flex items-center justify-between text-white">
					<div>
						<p className="text-blue-100 text-sm mb-2">Available Balance</p>
						<h2 className="text-4xl font-bold">₹{(walletData.balance / 100).toFixed(2)}</h2>
						<div className="mt-4 space-y-1">
							<p className="text-blue-100 text-xs">Lifetime Earned: ₹{(walletData.lifetime_earned / 100).toFixed(2)}</p>
							<p className="text-blue-100 text-xs">Lifetime Spent: ₹{(walletData.lifetime_spent / 100).toFixed(2)}</p>
						</div>
					</div>
					<div className="w-16 h-16 bg-white/20 rounded-full flex items-center justify-center">
						<svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
						</svg>
					</div>
				</div>
				<div className="mt-6 flex gap-3">
					<button
						onClick={() => {
							setShowAddMoney(!showAddMoney)
							setShowTransfer(false)
						}}
						className="bg-white text-blue-600 px-6 py-2 rounded-lg font-semibold hover:bg-blue-50 transition-colors"
					>
						+ Add Money
					</button>
					<button
						onClick={() => {
							setShowTransfer(!showTransfer)
							setShowAddMoney(false)
						}}
						className="bg-white/20 text-white border border-white/30 px-6 py-2 rounded-lg font-semibold hover:bg-white/30 transition-colors"
					>
						Transfer
					</button>
				</div>
			</div>

			{/* Add Money Form - Razorpay */}
			{showAddMoney && (
				<div className="bg-white rounded-lg shadow-md p-6 mb-6">
					<h3 className="text-xl font-semibold text-gray-900 mb-4">Add Money to Wallet</h3>

					{/* Quick Amount Selection */}
					<div className="mb-4">
						<p className="text-sm text-gray-600 mb-2">Quick Select:</p>
						<div className="grid grid-cols-3 sm:grid-cols-6 gap-2">
							{quickAmounts.map((amt) => (
								<button
									key={amt}
									onClick={() => setAddAmount(amt.toString())}
									className={`px-3 py-2 rounded-lg border-2 transition-all text-sm font-medium ${addAmount === amt.toString()
										? 'border-blue-500 bg-blue-50 text-blue-700'
										: 'border-gray-200 hover:border-blue-300'
										}`}
								>
									₹{amt}
								</button>
							))}
						</div>
					</div>

					{/* Custom Amount Input */}
					<div className="mb-4">
						<p className="text-sm text-gray-600 mb-2">Or enter custom amount:</p>
						<div className="relative">
							<span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">₹</span>
							<input
								type="number"
								value={addAmount}
								onChange={(e) => setAddAmount(e.target.value)}
								placeholder="Enter amount"
								min="1"
								max="100000"
								className="w-full pl-8 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
							/>
						</div>
						<p className="text-xs text-gray-500 mt-1">Min: ₹1 | Max: ₹1,00,000</p>
					</div>

					{/* Action Buttons */}
					<div className="flex gap-4">
						<button
							onClick={handleAddMoney}
							disabled={paymentLoading || !addAmount}
							className="flex-1 bg-blue-600 text-white px-6 py-2 rounded-lg font-semibold hover:bg-blue-700 transition-colors disabled:bg-gray-300 disabled:cursor-not-allowed flex items-center justify-center gap-2"
						>
							{paymentLoading ? (
								<>
									<div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
									Processing...
								</>
							) : (
								'Proceed to Payment'
							)}
						</button>
						<button
							onClick={() => setShowAddMoney(false)}
							className="bg-gray-200 text-gray-700 px-6 py-2 rounded-lg font-semibold hover:bg-gray-300 transition-colors"
						>
							Cancel
						</button>
					</div>

					{/* Security Note */}
					<div className="mt-4 flex items-center gap-2 text-xs text-gray-500">
						<svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
						</svg>
						<span>Payments are processed securely via Razorpay</span>
					</div>
				</div>
			)}

			{/* Transfer Money Form */}
			{showTransfer && (
				<div className="bg-white rounded-lg shadow-md p-6 mb-6">
					<h3 className="text-xl font-semibold text-gray-900 mb-4">Transfer Money</h3>
					<div className="space-y-4">
						<input
							type="number"
							value={recipientId}
							onChange={(e) => setRecipientId(e.target.value)}
							placeholder="Recipient User ID"
							className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
						/>
						<input
							type="number"
							value={transferAmount}
							onChange={(e) => setTransferAmount(e.target.value)}
							placeholder="Amount to transfer"
							className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
						/>
						<div className="flex gap-4">
							<button
								onClick={handleTransfer}
								className="flex-1 bg-blue-600 text-white px-6 py-2 rounded-lg font-semibold hover:bg-blue-700 transition-colors"
							>
								Transfer
							</button>
							<button
								onClick={() => setShowTransfer(false)}
								className="flex-1 bg-gray-200 text-gray-700 px-6 py-2 rounded-lg font-semibold hover:bg-gray-300 transition-colors"
							>
								Cancel
							</button>
						</div>
					</div>
				</div>
			)}

			{/* Quick Actions */}
			<div className="grid grid-cols-3 gap-4 mb-6">
				<button
					onClick={() => {
						setShowAddMoney(true)
						setShowTransfer(false)
					}}
					className="bg-white rounded-lg shadow-md p-4 hover:shadow-lg transition-shadow"
				>
					<div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-2">
						<svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
						</svg>
					</div>
					<p className="text-sm font-semibold text-gray-700">Add Money</p>
				</button>
				<button
					onClick={() => {
						setShowTransfer(true)
						setShowAddMoney(false)
					}}
					className="bg-white rounded-lg shadow-md p-4 hover:shadow-lg transition-shadow"
				>
					<div className="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center mx-auto mb-2">
						<svg className="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
						</svg>
					</div>
					<p className="text-sm font-semibold text-gray-700">Transfer</p>
				</button>
				<button
					onClick={fetchWalletData}
					className="bg-white rounded-lg shadow-md p-4 hover:shadow-lg transition-shadow"
				>
					<div className="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center mx-auto mb-2">
						<svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
						</svg>
					</div>
					<p className="text-sm font-semibold text-gray-700">Refresh</p>
				</button>
			</div>

			{/* Wallet Statistics */}
			<div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
				<div className="bg-white rounded-lg shadow-md p-6">
					<div className="flex items-center justify-between">
						<div>
							<p className="text-gray-500 text-sm mb-1">Current Balance</p>
							<p className="text-2xl font-bold text-gray-900">₹{(walletData.balance / 100).toFixed(2)}</p>
						</div>
						<div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
							<svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
							</svg>
						</div>
					</div>
				</div>

				<div className="bg-white rounded-lg shadow-md p-6">
					<div className="flex items-center justify-between">
						<div>
							<p className="text-gray-500 text-sm mb-1">Total Earned</p>
							<p className="text-2xl font-bold text-green-600">₹{(walletData.lifetime_earned / 100).toFixed(2)}</p>
						</div>
						<div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
							<svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
							</svg>
						</div>
					</div>
				</div>

				<div className="bg-white rounded-lg shadow-md p-6">
					<div className="flex items-center justify-between">
						<div>
							<p className="text-gray-500 text-sm mb-1">Total Spent</p>
							<p className="text-2xl font-bold text-red-600">₹{(walletData.lifetime_spent / 100).toFixed(2)}</p>
						</div>
						<div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
							<svg className="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6" />
							</svg>
						</div>
					</div>
				</div>
			</div>

			{/* Info Section */}
			<div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
				<div className="flex items-start gap-3">
					<svg className="w-5 h-5 text-blue-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
					<div className="text-sm text-blue-800">
						<p className="font-semibold mb-1">Note:</p>
						<p>• Add money securely using Razorpay (UPI, Cards, Net Banking)</p>
						<p>• Transfer money to other users using their User ID</p>
						<p>• All amounts are in Indian Rupees (₹)</p>
						<p>• Your wallet will be credited within seconds after successful payment</p>
					</div>
				</div>
			</div>
		</div>
	)
}
