"use client"

import ProductCard from "./product-card"

interface MarketplaceViewProps {
  searchQuery: string
  selectedCategory: string | null
  postedItems: any[]
  onAddToCart: (item: any) => void
}

const mockProducts = [
  {
    id: 1,
    name: "Advanced Calculus Textbook",
    category: "Books",
    price: 450,
    condition: "Like New",
    seller: "Rajesh Kumar",
    phone: "+91 9876543210",
    image: "/calculus-textbook.png",
  },
  {
    id: 2,
    name: "Gaming Laptop",
    category: "Electronics",
    price: 35000,
    condition: "Good",
    seller: "Priya Singh",
    phone: "+91 9123456789",
    image: "/gaming-laptop.png",
  },
  {
    id: 3,
    name: "Wooden Study Desk",
    category: "Furniture",
    price: 5000,
    condition: "Used",
    seller: "Arjun Patel",
    phone: "+91 8765432109",
    image: "/wooden-study-desk.jpg",
  },
  {
    id: 4,
    name: "Winter Jacket",
    category: "Clothing",
    price: 1200,
    condition: "Like New",
    seller: "Sneha Gupta",
    phone: "+91 7654321098",
    image: "/winter-jacket.png",
  },
  {
    id: 5,
    name: "Physics Reference Guide",
    category: "Books",
    price: 350,
    condition: "Good",
    seller: "Rohit Sharma",
    phone: "+91 6543210987",
    image: "/physics-reference.jpg",
  },
  {
    id: 6,
    name: "USB-C Hub",
    category: "Electronics",
    price: 800,
    condition: "New",
    seller: "Ananya Verma",
    phone: "+91 5432109876",
    image: "/usb-hub.png",
  },
  {
    id: 7,
    name: "Bookshelf",
    category: "Furniture",
    price: 3500,
    condition: "Good",
    seller: "Vikram Singh",
    phone: "+91 4321098765",
    image: "/cozy-bookshelf.png",
  },
  {
    id: 8,
    name: "Formal Shirt",
    category: "Clothing",
    price: 600,
    condition: "Used",
    seller: "Pooja Sharma",
    phone: "+91 3210987654",
    image: "/formal-shirt.png",
  },
]

export default function MarketplaceView({
  searchQuery,
  selectedCategory,
  postedItems,
  onAddToCart,
}: MarketplaceViewProps) {
  const allProducts = [...mockProducts, ...postedItems]

  const filteredProducts = allProducts.filter((product) => {
    const matchesSearch =
      !searchQuery ||
      product.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      product.category.toLowerCase().includes(searchQuery.toLowerCase())

    const matchesCategory = !selectedCategory || product.category === selectedCategory

    return matchesSearch && matchesCategory
  })

  return (
    <div>
      <h3 className="text-xl font-bold text-gray-900 mb-6">Available Items</h3>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 md:gap-6">
        {filteredProducts.map((product, index) => (
          <ProductCard key={`${product.id}-${index}`} product={product} onAddToCart={onAddToCart} />
        ))}
      </div>
      {filteredProducts.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-500 text-lg">No items found matching your search.</p>
        </div>
      )}
    </div>
  )
}
