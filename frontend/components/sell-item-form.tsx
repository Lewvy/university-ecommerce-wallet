"use client"

import type React from "react"

import { useState } from "react"

interface SellItemFormProps {
  onSubmit: (item: any) => void
  userData: {
    username: string
    email: string
    phone: string
  }
}

export default function SellItemForm({ onSubmit, userData }: SellItemFormProps) {
  const [formData, setFormData] = useState({
    productName: "",
    category: "Books",
    description: "",
    price: "",
    condition: "Good",
    imageUrl: "",
  })

  const [submitted, setSubmitted] = useState(false)

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const newItem = {
      ...formData,
      price: Number.parseInt(formData.price),
      seller: userData.username,
      phone: userData.phone,
    }
    onSubmit(newItem)
    setFormData({
      productName: "",
      category: "Books",
      description: "",
      price: "",
      condition: "Good",
      imageUrl: "",
    })
    setSubmitted(true)
    setTimeout(() => setSubmitted(false), 3000)
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
          <h3 className="text-lg font-semibold text-green-900 mb-2">Item Listed Successfully!</h3>
          <p className="text-green-700 text-sm md:text-base">
            Your item has been published and is now visible to other students.
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-2xl mx-auto px-4">
      <h3 className="text-2xl font-bold text-gray-900 mb-6">Sell an Item</h3>

      <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow p-6 md:p-8">
        <div className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Product Name</label>
            <input
              type="text"
              name="productName"
              value={formData.productName}
              onChange={handleChange}
              placeholder="e.g., Advanced Physics Textbook"
              required
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm"
            />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Category</label>
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
              <label className="block text-sm font-medium text-gray-700 mb-2">Condition</label>
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
            <label className="block text-sm font-medium text-gray-700 mb-2">Description</label>
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
              <label className="block text-sm font-medium text-gray-700 mb-2">Price (₹)</label>
              <input
                type="number"
                name="price"
                value={formData.price}
                onChange={handleChange}
                placeholder="0"
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Image URL</label>
              <input
                type="url"
                name="imageUrl"
                value={formData.imageUrl}
                onChange={handleChange}
                placeholder="https://example.com/image.jpg"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm"
              />
            </div>
          </div>

          <button
            type="submit"
            className="w-full bg-gradient-to-r from-blue-600 to-blue-500 text-white font-semibold py-3 rounded-lg hover:shadow-lg transition-shadow text-sm md:text-base"
          >
            Publish Listing
          </button>
        </div>
      </form>
    </div>
  )
}
