"use client"

interface ProductCardProps {
  product: {
    id?: number
    name: string
    category: string
    price: number
    condition: string
    seller: string
    phone: string
    image?: string
  }
  onAddToCart: (item: any) => void
}

export default function ProductCard({ product, onAddToCart }: ProductCardProps) {
  return (
    <div className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow overflow-hidden flex flex-col h-full">
      <img
        src={product.image || "/placeholder.svg?height=200&width=300&query=product"}
        alt={product.name}
        className="w-full h-40 md:h-48 object-cover"
      />
      <div className="p-3 md:p-4 flex flex-col flex-1">
        <div className="flex items-start justify-between mb-2">
          <div className="flex-1 min-w-0">
            <h4 className="font-semibold text-gray-900 text-xs md:text-sm leading-tight mb-1 truncate">
              {product.name}
            </h4>
            <span className="inline-block bg-blue-100 text-blue-700 text-xs px-2 py-0.5 rounded">
              {product.category}
            </span>
          </div>
        </div>

        <div className="mb-3 text-xs">
          <p className="text-gray-600 mb-1">
            Condition: <span className="font-medium">{product.condition}</span>
          </p>
          <p className="text-gray-600 truncate">
            Seller: <span className="font-medium">{product.seller}</span>
          </p>
        </div>

        <div className="border-t border-gray-200 pt-3 mb-3 mt-auto">
          <p className="text-base md:text-lg font-bold text-blue-600 mb-2">₹{Number(product.price).toLocaleString()}</p>
          <p className="text-xs text-gray-500 break-all">{product.phone}</p>
        </div>

        <button
          onClick={() => onAddToCart(product)}
          className="w-full bg-gradient-to-r from-blue-600 to-blue-500 text-white text-xs md:text-sm font-semibold py-2 rounded hover:shadow-lg transition-shadow"
        >
          Add to Cart
        </button>
      </div>
    </div>
  )
}
