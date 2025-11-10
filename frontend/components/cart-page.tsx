"use client"

interface CartPageProps {
  items: any[]
}

export default function CartPage({ items }: CartPageProps) {
  const total = items.reduce((sum, item) => sum + item.price, 0)

  return (
    <div className="max-w-4xl mx-auto px-4">
      <h3 className="text-2xl font-bold text-gray-900 mb-6">My Cart</h3>

      {items.length === 0 ? (
        <div className="bg-white rounded-lg shadow p-8 md:p-12 text-center">
          <svg
            className="w-12 md:w-16 h-12 md:h-16 text-gray-400 mx-auto mb-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"
            />
          </svg>
          <p className="text-gray-500 text-lg">Your cart is empty</p>
          <p className="text-gray-400 mt-2 text-sm md:text-base">Start adding items to get started!</p>
        </div>
      ) : (
        <>
          <div className="space-y-3 md:space-y-4 mb-6">
            {items.map((item, index) => (
              <div
                key={index}
                className="bg-white rounded-lg shadow p-3 md:p-4 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4"
              >
                <div className="flex items-start sm:items-center space-x-3 md:space-x-4 flex-1 min-w-0">
                  <img
                    src={item.image || "/placeholder.svg?height=80&width=80&query=product"}
                    alt={item.name}
                    className="w-16 md:w-20 h-16 md:h-20 object-cover rounded flex-shrink-0"
                  />
                  <div className="min-w-0">
                    <h4 className="font-semibold text-gray-900 text-sm md:text-base truncate">{item.name}</h4>
                    <p className="text-xs md:text-sm text-gray-600">{item.category}</p>
                    <p className="text-xs md:text-sm text-gray-600">Condition: {item.condition}</p>
                  </div>
                </div>
                <p className="text-base md:text-lg font-bold text-blue-600 flex-shrink-0">
                  ₹{item.price.toLocaleString()}
                </p>
              </div>
            ))}
          </div>

          <div className="bg-white rounded-lg shadow p-4 md:p-6">
            <div className="border-t border-gray-200 pt-4">
              <div className="flex items-center justify-between mb-4">
                <span className="text-base md:text-lg font-semibold text-gray-900">Total:</span>
                <span className="text-xl md:text-2xl font-bold text-blue-600">₹{total.toLocaleString()}</span>
              </div>
              <button className="w-full bg-gradient-to-r from-blue-600 to-blue-500 text-white font-semibold py-3 rounded-lg hover:shadow-lg transition-shadow text-sm md:text-base">
                Proceed to Checkout
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  )
}
