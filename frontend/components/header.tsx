"use client"

import { useState } from "react"

interface HeaderProps {
  username: string
  searchQuery: string
  onSearchChange: (query: string) => void
  selectedCategory: string | null
  onCategoryChange: (category: string | null) => void
}

const categories = ["Books", "Electronics", "Furniture", "Clothing", "Others"]

export default function Header({
  username,
  searchQuery,
  onSearchChange,
  selectedCategory,
  onCategoryChange,
}: HeaderProps) {
  const [showCategoryDropdown, setShowCategoryDropdown] = useState(false)

  return (
    <header className="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-10">
      <div className="px-4 md:px-6 py-4">
        <h2 className="text-xl md:text-2xl font-bold text-gray-900 mb-4">
          Hi {username}! What are you looking for today?
        </h2>
        <div className="flex flex-col sm:flex-row gap-3">
          <div className="flex-1 relative">
            <input
              type="text"
              placeholder="Search items..."
              value={searchQuery}
              onChange={(e) => onSearchChange(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-sm md:text-base"
            />
            <svg
              className="absolute right-4 top-3.5 w-5 h-5 text-gray-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              />
            </svg>
          </div>

          <div className="relative">
            <button
              onClick={() => setShowCategoryDropdown(!showCategoryDropdown)}
              className="w-full sm:w-auto px-4 py-3 border border-gray-300 rounded-lg bg-white text-gray-700 hover:bg-gray-50 font-medium text-sm md:text-base flex items-center justify-between gap-2"
            >
              <span>{selectedCategory || "All Categories"}</span>
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
              </svg>
            </button>

            {showCategoryDropdown && (
              <div className="absolute top-full left-0 right-0 mt-2 bg-white border border-gray-300 rounded-lg shadow-lg z-20">
                <button
                  onClick={() => {
                    onCategoryChange(null)
                    setShowCategoryDropdown(false)
                  }}
                  className="w-full text-left px-4 py-2 hover:bg-gray-100 text-gray-700 text-sm"
                >
                  All Categories
                </button>
                {categories.map((cat) => (
                  <button
                    key={cat}
                    onClick={() => {
                      onCategoryChange(cat)
                      setShowCategoryDropdown(false)
                    }}
                    className={`w-full text-left px-4 py-2 text-sm ${
                      selectedCategory === cat ? "bg-blue-100 text-blue-700" : "text-gray-700 hover:bg-gray-100"
                    }`}
                  >
                    {cat}
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  )
}
