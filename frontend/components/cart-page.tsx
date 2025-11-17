"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

interface CartItem {
	id?: number;
	product_id?: number;
	name?: string;
	image_url?: string;
	category?: string;
	condition?: string;
	price?: number;
	quantity?: number;
	[key: string]: any;
}

export default function CartPage() {
	const router = useRouter();
	const [items, setItems] = useState<CartItem[]>([]);
	const [loading, setLoading] = useState(true);
	const [updating, setUpdating] = useState<number | null>(null);
	const [checkingOut, setCheckingOut] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		fetchCart();
	}, []);

	const fetchCart = async () => {
		try {
			const token = localStorage.getItem("access_token");
			if (!token) {
				setLoading(false);
				return;
			}

			const res = await fetch("http://localhost:8088/cart", {
				headers: {
					Authorization: `Bearer ${token}`,
					"Content-Type": "application/json",
				},
			});

			if (res.ok) {
				const data = await res.json();

				if (Array.isArray(data)) {
					setItems(data);
				} else if (data.items && Array.isArray(data.items)) {
					setItems(data.items);
				} else if (data.Items && Array.isArray(data.Items)) {
					setItems(data.Items);
				} else {
					setItems([]);
				}
			} else {
				const errorText = await res.text();
				setError(`Failed to load cart: ${res.status}`);
			}
		} catch (err) {
			setError("Failed to connect to server");
		} finally {
			setLoading(false);
		}
	};

	const updateQuantity = async (productId: number, newQuantity: number) => {
		if (newQuantity < 1) return;

		setUpdating(productId);
		setError(null);

		try {
			const token = localStorage.getItem("access_token");
			if (!token) throw new Error("Please log in");

			const res = await fetch("http://localhost:8088/cart/update", {
				method: "PUT",
				headers: {
					Authorization: `Bearer ${token}`,
					"Content-Type": "application/json",
				},
				body: JSON.stringify({
					product_id: productId,
					quantity: newQuantity,
				}),
			});

			if (res.ok) {
				await fetchCart();
			} else {
				const errorData = await res.json();
				throw new Error(errorData.error || "Failed to update quantity");
			}
		} catch (err: any) {
			setError(err.message);
		} finally {
			setUpdating(null);
		}
	};

	const removeItem = async (productId: number) => {
		setUpdating(productId);
		setError(null);

		try {
			const token = localStorage.getItem("access_token");
			if (!token) throw new Error("Please log in");

			const res = await fetch(`http://localhost:8088/cart/item/${productId}`, {
				method: "DELETE",
				headers: {
					Authorization: `Bearer ${token}`,
					"Content-Type": "application/json",
				},
			});

			if (res.ok) {
				await fetchCart();
			} else {
				const errorData = await res.json();
				throw new Error(errorData.error || "Failed to remove item");
			}
		} catch (err: any) {
			setError(err.message);
		} finally {
			setUpdating(null);
		}
	};

	const clearCart = async () => {
		if (!confirm("Are you sure you want to clear your cart?")) return;

		setLoading(true);
		setError(null);

		try {
			const token = localStorage.getItem("access_token");
			if (!token) throw new Error("Please log in");

			const res = await fetch("http://localhost:8088/cart/clear", {
				method: "DELETE",
				headers: {
					Authorization: `Bearer ${token}`,
					"Content-Type": "application/json",
				},
			});

			if (res.ok) {
				setItems([]);
			} else {
				const errorData = await res.json();
				throw new Error(errorData.error || "Failed to clear cart");
			}
		} catch (err: any) {
			setError(err.message);
		} finally {
			setLoading(false);
		}
	};

	const proceedToCheckout = async () => {
		setCheckingOut(true);
		setError(null);

		try {
			const token = localStorage.getItem("access_token");
			if (!token) throw new Error("Please log in to checkout");

			const res = await fetch("http://localhost:8088/orders", {
				method: "POST",
				headers: {
					Authorization: `Bearer ${token}`,
					"Content-Type": "application/json",
				},
			});

			if (res.ok) {
				alert("Order created successfully!");
				setItems([]);
			} else {
				const errorData = await res.json();
				throw new Error(errorData.error || "Failed to create order");
			}
		} catch (err: any) {
			setError(err.message);
			alert(`Checkout failed: ${err.message}`);
		} finally {
			setCheckingOut(false);
		}
	};

	const total = items.reduce(
		(sum, item) => sum + (item.price || 0) * (item.quantity || 0),
		0
	);

	if (loading) {
		return (
			<div className="max-w-4xl mx-auto px-4 py-8">
				<div className="text-center">Loading cart...</div>
			</div>
		);
	}

	return (
		<div className="max-w-4xl mx-auto px-4 py-6">
			<div className="flex justify-between items-center mb-6">
				<h3 className="text-2xl font-bold text-gray-900">My Cart</h3>
				{items.length > 0 && (
					<button
						onClick={clearCart}
						className="text-red-600 hover:text-red-700 text-sm font-medium"
					>
						Clear Cart
					</button>
				)}
			</div>

			{error && (
				<div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
					{error}
				</div>
			)}

			{items.length === 0 ? (
				<div className="bg-white rounded-lg shadow p-8 md:p-12 text-center">
					<p className="text-gray-500 text-lg">Your cart is empty</p>
					<button
						onClick={() => router.push("/")}
						className="mt-4 bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700"
					>
						Continue Shopping
					</button>
				</div>
			) : (
				<>
					<div className="space-y-3 md:space-y-4 mb-6">
						{items.map((item) => (
							<div
								key={item.product_id}
								className="bg-white rounded-lg shadow p-3 md:p-4 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
							>
								<div className="flex items-start sm:items-center space-x-3 md:space-x-4 flex-1 min-w-0">
									<img
										src={item.image_url || "/placeholder.svg?height=80&width=80"}
										alt={item.name}
										className="w-16 md:w-20 h-16 md:h-20 object-cover rounded flex-shrink-0"
									/>

									<div className="min-w-0 flex-1">
										<h4 className="font-semibold text-gray-900 text-sm md:text-base">
											{item.name}
										</h4>

										<p className="text-xs md:text-sm text-gray-600">{item.category}</p>

										<p className="text-xs md:text-sm text-gray-600">
											Condition: {item.condition}
										</p>

										<p className="text-sm font-medium text-blue-600 mt-1">
											₹{(item.price || 0).toLocaleString()} each
										</p>
									</div>
								</div>

								<div className="flex items-center gap-4 flex-shrink-0">
									<div className="flex items-center gap-2 border rounded-lg">
										<button
											onClick={() =>
												updateQuantity(
													item.product_id!,
													(item.quantity || 1) - 1
												)
											}
											disabled={
												updating === item.product_id ||
												(item.quantity || 1) <= 1
											}
											className="px-3 py-1 hover:bg-gray-100 disabled:opacity-50"
										>
											−
										</button>

										<span className="px-2 font-medium">
											{item.quantity}
										</span>

										<button
											onClick={() =>
												updateQuantity(
													item.product_id!,
													(item.quantity || 0) + 1
												)
											}
											disabled={updating === item.product_id}
											className="px-3 py-1 hover:bg-gray-100 disabled:opacity-50"
										>
											+
										</button>
									</div>

									<p className="text-base md:text-lg font-bold text-blue-600 w-24 text-right">
										₹{((item.price || 0) * (item.quantity || 0)).toLocaleString()}
									</p>

									<button
										onClick={() => removeItem(item.product_id!)}
										disabled={updating === item.product_id}
										className="text-red-600 hover:text-red-700 disabled:opacity-50"
									>
										<svg
											className="w-5 h-5"
											fill="none"
											stroke="currentColor"
											viewBox="0 0 24 24"
										>
											<path
												strokeLinecap="round"
												strokeLinejoin="round"
												strokeWidth={2}
												d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 00-1 1v3M4 7h16"
											/>
										</svg>
									</button>
								</div>
							</div>
						))}
					</div>

					<div className="bg-white rounded-lg shadow p-4 md:p-6">
						<div className="border-t border-gray-200 pt-4">
							<div className="flex items-center justify-between mb-4">
								<span className="text-base md:text-lg font-semibold text-gray-900">
									Total:
								</span>
								<span className="text-xl md:text-2xl font-bold text-blue-600">
									₹{total.toLocaleString()}
								</span>
							</div>

							<button
								onClick={proceedToCheckout}
								disabled={checkingOut}
								className="w-full bg-gradient-to-r from-blue-600 to-blue-500 text-white font-semibold py-3 rounded-lg hover:shadow-lg transition-shadow text-sm md:text-base disabled:opacity-50 disabled:cursor-not-allowed"
							>
								{checkingOut ? "Processing..." : "Proceed to Checkout"}
							</button>
						</div>
					</div>
				</>
			)}
		</div>
	);
}
