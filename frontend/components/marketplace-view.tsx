"use client"

import { useState, useEffect } from "react"
import ProductCard from "./product-card"

interface MarketplaceViewProps {
	searchQuery: string
	selectedCategory: string | null
	postedItems: any[]
	onAddToCart: (item: any) => void
}

interface Product {
	id: number
	name: string
	description: string
	price: number
	stock: number
	category?: string
	condition?: string
	seller_id: number
	created_at: string
	updated_at: string
	images?: string[]
}

export default function MarketplaceView({
	searchQuery,
	selectedCategory,
	postedItems,
	onAddToCart,
}: MarketplaceViewProps) {
	const [products, setProducts] = useState<Product[]>([])
	const [isLoading, setIsLoading] = useState(true)
	const [error, setError] = useState<string | null>(null)

	useEffect(() => {
		fetchProducts()
	}, [])

  const fetchProducts = async () => {
    setIsLoading(true)
    setError(null)

    try {
        const token = localStorage.getItem("access_token")

        const response = await fetch("http://localhost:8088/products", {
            headers: {
                "Content-Type": "application/json",
                "Authorization": token ? `Bearer ${token}` : ""
            }
        })

        if (!response.ok) {
            throw new Error("Failed to fetch products")
        }

        const data = await response.json()
        console.log("Fetched products:", data)

        const productList = Array.isArray(data) ? data : (data.products || [])
        setProducts(productList)
    } catch (err: any) {
        console.error("Error fetching products:", err)
        setError(err.message)
    } finally {
        setIsLoading(false)
    }
}


	/*const fetchProducts = async () => {
		setIsLoading(true)
		setError(null)
		try {
			const response = await fetch("http://localhost:8088/products", {
				headers: {
					"Content-Type": "application/json",
				}
			})

			if (!response.ok) {
				throw new Error("Failed to fetch products")
			}

			const data = await response.json()
			console.log("Fetched products:", data)
			
			// Handle both array and object with products array
			const productList = Array.isArray(data) ? data : (data.products || [])
			setProducts(productList)
		} catch (err: any) {
			console.error("Error fetching products:", err)
			setError(err.message)
		} finally {
			setIsLoading(false)
		}
	}*/

	// Transform backend products to match ProductCard interface
	// Transform backend products to match ProductCard interface
const transformedProducts = products.map((product: any) => ({
    id: product.ID,
    name: product.Name,
    category: "General",                 // backend has no category
    price: Number(product.Price),
    condition: "Good",                   // backend has no condition
    seller: `Seller ${product.SellerID}`, 
    phone: "Contact seller",
    image: product.ImageUrl 
        ? product.ImageUrl 
        : "/placeholder.svg?height=200&width=300&query=product",
    description: product.Description,
    stock: Number(product.Stock),
}))



	const allProducts = [...transformedProducts, ...postedItems]

	const filteredProducts = allProducts.filter((product) => {
		const matchesSearch =
			!searchQuery ||
			product.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			(product.category && product.category.toLowerCase().includes(searchQuery.toLowerCase()))

		const matchesCategory = !selectedCategory || product.category === selectedCategory

		return matchesSearch && matchesCategory
	})

	if (isLoading) {
		return (
			<div className="flex items-center justify-center h-64">
				<div className="text-center">
					<div className="w-12 h-12 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
					<p className="text-gray-600">Loading products...</p>
				</div>
			</div>
		)
	}

	if (error) {
		return (
			<div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
				<svg className="w-12 h-12 text-red-600 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>
				<h3 className="text-lg font-semibold text-red-900 mb-2">Error Loading Products</h3>
				<p className="text-red-700 mb-4">{error}</p>
				<button
					onClick={fetchProducts}
					className="bg-red-600 text-white px-6 py-2 rounded-lg font-semibold hover:bg-red-700 transition-colors"
				>
					Try Again
				</button>
			</div>
		)
	}

	return (
		<div>
			<div className="flex items-center justify-between mb-6">
				<h3 className="text-xl font-bold text-gray-900">Available Items</h3>
				<button
					onClick={fetchProducts}
					className="text-blue-600 hover:text-blue-700 font-semibold text-sm flex items-center gap-2"
				>
					<svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
					</svg>
					Refresh
				</button>
			</div>

			{filteredProducts.length === 0 ? (
				<div className="text-center py-12">
					<div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
						<svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
						</svg>
					</div>
					<p className="text-gray-500 text-lg">No items found matching your search.</p>
					{(searchQuery || selectedCategory) && (
						<p className="text-gray-400 text-sm mt-2">Try adjusting your filters</p>
					)}
				</div>
			) : (
				<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 md:gap-6">
					{filteredProducts.map((product, index) => (
						<ProductCard key={`${product.id}-${index}`} product={product} onAddToCart={onAddToCart} />
					))}
				</div>
			)}
		</div>
	)
}

/*"use client"

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
*/