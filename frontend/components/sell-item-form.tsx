"use client"

import { useState } from "react"

interface SellItemFormProps {
	onSubmit: (item: any) => void
	userData: {
		username: string
		email: string
		phone: string
		id?: number
	}
}

export default function SellItemForm({ onSubmit, userData }: SellItemFormProps) {
	const [formData, setFormData] = useState({
		name: "",
		category: "Books",
		description: "",
		price: "",
		condition: "Good",
		stock: "1",
	})

	const [images, setImages] = useState<File[]>([])
	const [imagePreviews, setImagePreviews] = useState<string[]>([])
	const [isSubmitting, setIsSubmitting] = useState(false)
	const [submitted, setSubmitted] = useState(false)
	const [error, setError] = useState<string | null>(null)

	const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
		const { name, value } = e.target
		setFormData((prev) => ({ ...prev, [name]: value }))
		setError(null)
	}

	const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		const files = e.target.files
		if (!files) return

		const fileArray = Array.from(files)
		
		if (fileArray.length + images.length > 5) {
			setError("Maximum 5 images allowed")
			return
		}

		const validFiles: File[] = []
		const newPreviews: string[] = []

		fileArray.forEach(file => {
			if (!file.type.startsWith('image/')) {
				setError(`${file.name} is not an image file`)
				return
			}

			if (file.size > 5 * 1024 * 1024) {
				setError(`${file.name} is too large. Max size is 5MB`)
				return
			}

			validFiles.push(file)
			
			const reader = new FileReader()
			reader.onloadend = () => {
				newPreviews.push(reader.result as string)
				setImagePreviews(prev => [...prev, reader.result as string])
			}
			reader.readAsDataURL(file)
		})

		setImages(prev => [...prev, ...validFiles])
	}

	const removeImage = (index: number) => {
		setImages(prev => prev.filter((_, i) => i !== index))
		setImagePreviews(prev => prev.filter((_, i) => i !== index))
	}

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault()
		setError(null)
		setIsSubmitting(true)

		try {
			const token = localStorage.getItem("access_token")
			
			if (!token) {
				throw new Error("You must be logged in to sell items")
			}

			const formDataToSend = new FormData()
			
			formDataToSend.append("name", formData.name)
			formDataToSend.append("description", formData.description)
			formDataToSend.append("price", formData.price)
			formDataToSend.append("stock", formData.stock)
			formDataToSend.append("category", formData.category)
			formDataToSend.append("condition", formData.condition)

			images.forEach((image) => {
				formDataToSend.append("images", image)
			})

			console.log("Submitting product:", {
				name: formData.name,
				price: formData.price,
				stock: formData.stock,
				imageCount: images.length
			})


			const response = await fetch("http://localhost:8088/products", {
				method: "POST",
				headers: {
					"Authorization": `Bearer ${token}`,
				},
				body: formDataToSend,
			})

			if (!response.ok) {
				const errorData = await response.json()
				throw new Error(errorData.error || errorData.message || "Failed to create product")
			}

			const data = await response.json()
			console.log("Product created successfully:", data)

			onSubmit({
				...formData,
				price: parseInt(formData.price),
				seller: userData.username,
				phone: userData.phone,
				image: imagePreviews[0] || "",
			})

			setFormData({
				name: "",
				category: "Books",
				description: "",
				price: "",
				condition: "Good",
				stock: "1",
			})
			setImages([])
			setImagePreviews([])
			setSubmitted(true)

			setTimeout(() => {
				setSubmitted(false)
			}, 3000)

		} catch (err: any) {
			console.error("Error creating product:", err)
			setError(err.message || "Failed to create product. Please try again.")
		} finally {
			setIsSubmitting(false)
		}
	}

	if (submitted) {
		return (
			<div className="max-w-2xl mx-auto px-4">
				<div className="bg-green-50 border border-green-200 rounded-lg p-6 text-center">
					<div className="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
						<svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
						</svg>
					</div>
					<h3 className="text-lg font-semibold text-green-900 mb-2">Product Listed Successfully!</h3>
					<p className="text-green-700 text-sm md:text-base">
						Your product has been published and is now visible to other students.
					</p>
				</div>
			</div>
		)
	}

	return (
		<div className="max-w-2xl mx-auto px-4">
			<h3 className="text-2xl font-bold text-gray-900 mb-6">Sell an Item</h3>

			{error && (
				<div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
					<div className="flex items-center gap-2">
						<svg className="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
						<p className="text-red-800 text-sm">{error}</p>
					</div>
				</div>
			)}

			<form onSubmit={handleSubmit} className="bg-white rounded-lg shadow p-6 md:p-8">
				<div className="space-y-6">
					<div>
						<label className="block text-sm font-medium text-gray-700 mb-2">Product Name *</label>
						<input
							type="text"
							name="name"
							value={formData.name}
							onChange={handleChange}
							placeholder="e.g., Advanced Physics Textbook"
							required
							className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm"
						/>
					</div>

					<div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<div>
							<label className="block text-sm font-medium text-gray-700 mb-2">Category *</label>
							<select
								name="category"
								value={formData.category}
								onChange={handleChange}
								className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm"
							>
								<option>Books</option>
								<option>Electronics</option>
								<option>Furniture</option>
								<option>Clothing</option>
								<option>Others</option>
							</select>
						</div>

						<div>
							<label className="block text-sm font-medium text-gray-700 mb-2">Condition *</label>
							<select
								name="condition"
								value={formData.condition}
								onChange={handleChange}
								className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm"
							>
								<option>New</option>
								<option>Like New</option>
								<option>Good</option>
								<option>Used</option>
							</select>
						</div>
					</div>

					<div>
						<label className="block text-sm font-medium text-gray-700 mb-2">Description *</label>
						<textarea
							name="description"
							value={formData.description}
							onChange={handleChange}
							placeholder="Describe the item condition and any details..."
							rows={4}
							required
							className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm"
						/>
					</div>

					<div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<div>
							<label className="block text-sm font-medium text-gray-700 mb-2">Price (₹) *</label>
							<input
								type="number"
								name="price"
								value={formData.price}
								onChange={handleChange}
								placeholder="0"
								required
								min="1"
								className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm"
							/>
						</div>

						<div>
							<label className="block text-sm font-medium text-gray-700 mb-2">Stock Quantity *</label>
							<input
								type="number"
								name="stock"
								value={formData.stock}
								onChange={handleChange}
								placeholder="1"
								required
								min="1"
								className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm"
							/>
						</div>
					</div>

					{}
					<div>
						<label className="block text-sm font-medium text-gray-700 mb-2">
							Product Images (Max 5)
						</label>
						<div className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-500 transition-colors">
							<input
								type="file"
								accept="image/*"
								multiple
								onChange={handleImageChange}
								className="hidden"
								id="image-upload"
								disabled={images.length >= 5}
							/>
							<label
								htmlFor="image-upload"
								className={`cursor-pointer ${images.length >= 5 ? 'opacity-50 cursor-not-allowed' : ''}`}
							>
								<svg className="w-12 h-12 text-gray-400 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
								</svg>
								<p className="text-sm text-gray-600 mb-1">
									{images.length >= 5 ? 'Maximum images reached' : 'Click to upload images'}
								</p>
								<p className="text-xs text-gray-500">PNG, JPG, WEBP up to 5MB each</p>
							</label>
						</div>

						{}
						{imagePreviews.length > 0 && (
							<div className="grid grid-cols-3 gap-4 mt-4">
								{imagePreviews.map((preview, index) => (
									<div key={index} className="relative group">
										<img
											src={preview}
											alt={`Preview ${index + 1}`}
											className="w-full h-24 object-cover rounded-lg border border-gray-200"
										/>
										<button
											type="button"
											onClick={() => removeImage(index)}
											className="absolute top-1 right-1 bg-red-500 text-white rounded-full p-1 opacity-0 group-hover:opacity-100 transition-opacity"
										>
											<svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
											</svg>
										</button>
									</div>
								))}
							</div>
						)}
					</div>

					<button
						type="submit"
						disabled={isSubmitting || images.length === 0}
						className="w-full bg-gradient-to-r from-blue-600 to-blue-500 text-white font-semibold py-3 rounded-lg hover:shadow-lg transition-shadow text-sm md:text-base disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{isSubmitting ? (
							<span className="flex items-center justify-center gap-2">
								<svg className="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
									<circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
									<path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
								</svg>
								Publishing...
							</span>
						) : (
							'Publish Listing'
						)}
					</button>
				</div>
			</form>
		</div>
	)
}