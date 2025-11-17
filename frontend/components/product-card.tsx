"use client"

import { useRouter } from "next/navigation"

interface ProductCardProps {
  product: {
    id: number
    name: string
    category?: string
    price: number
    condition?: string
    seller?: string
    image?: string
  }
  onAddToCart: (item: any) => void
}

export default function ProductCard({ product, onAddToCart }: ProductCardProps) {
  const router = useRouter()

  return (
    <div className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow overflow-hidden flex flex-col h-full">

      {/* CLICKABLE IMAGE */}
      <img
        onClick={() => router.push(`/product/${product.id}`)}
        src={product.image || "/placeholder.jpg"}
        alt={product.name}
        className="w-full h-40 md:h-48 object-cover cursor-pointer"
      />

      <div className="p-3 md:p-4 flex flex-col flex-1">
        <h4 className="font-semibold text-gray-900 text-sm leading-tight mb-1 truncate">
          {product.name}
        </h4>

        <p className="text-gray-500 text-xs mb-2">
          {product.category || "General"}
        </p>

        <p className="text-base font-bold text-blue-600 mb-2">
          ₹{Number(product.price).toLocaleString()}
        </p>

        <button
          onClick={() => onAddToCart(product)}
          className="w-full bg-blue-600 text-white text-sm font-semibold py-2 rounded hover:shadow-lg transition-shadow"
        >
          Add to Cart
        </button>
      </div>
    </div>
  )
}
