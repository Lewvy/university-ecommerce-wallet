"use client"

import { useState } from "react"

interface CategoriesViewProps {
  onCategorySelect: (category: string) => void
  postedItems: any[]
}

const defaultCategories = ["Books", "Electronics", "Furniture", "Clothing", "Others"]

export default function CategoriesView({ onCategorySelect, postedItems }: CategoriesViewProps) {
  const [hoveredCategory, setHoveredCategory] = useState<string | null>(null)

  // Count items by category
  const mockProductCounts: Record<string, number> = {
    Books: 24,
    Electronics: 18,
    Furniture: 12,
    Clothing: 31,
    Others: 8,
  }

  const getCategoryCount = (category: string) => {
    const postedCount = postedItems.filter((item) => item.category === category).length
    return mockProductCounts[category] + postedCount
  }

  return (
    <div className="px-4">
      <h3 className="text-2xl font-bold text-gray-900 mb-6">Browse by Category</h3>

      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3 md:gap-6">
        {defaultCategories.map((category) => (
          <button
            key={category}
            onMouseEnter={() => setHoveredCategory(category)}
            onMouseLeave={() => setHoveredCategory(null)}
            onClick={() => onCategorySelect(category)}
            className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow p-4 md:p-6 text-center group cursor-pointer active:scale-95"
          >
            <div
              className={`text-4xl md:text-5xl mb-4 transition-transform ${hoveredCategory === category ? "scale-110" : ""}`}
            >
              {category === "Books" && "📚"}
              {category === "Electronics" && "💻"}
              {category === "Furniture" && "🪑"}
              {category === "Clothing" && "👕"}
              {category === "Others" && "📦"}
            </div>
            <h4 className="font-bold text-gray-900 mb-2 text-sm md:text-base">{category}</h4>
            <p className="text-xs md:text-sm text-gray-600">{getCategoryCount(category)} items</p>
          </button>
        ))}
      </div>
    </div>
  )
}
