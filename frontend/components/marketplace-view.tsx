"use client";

import { useState, useEffect } from "react";
import ProductCard from "./product-card";

interface MarketplaceViewProps {
    searchQuery: string;
    selectedCategory: string | null;
    postedItems: any[];
    onAddToCart: (item: any) => void;
}

interface Product {
    ID: number;
    Name: string;
    Description: string;
    Price: number;
    Stock: number;
    Category?: string;
    Condition?: string;
    SellerID: number;
    SellerName?: string;
    SellerPhone?: string;
    ImageUrl?: string;
    Images?: string[];
}

export default function MarketplaceView({
    searchQuery,
    selectedCategory,
    postedItems,
    onAddToCart,
}: MarketplaceViewProps) {
    const [products, setProducts] = useState<Product[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // ---------------- FETCH PRODUCTS ----------------
    useEffect(() => {
        fetchProducts();
    }, []);

    const fetchProducts = async () => {
        setIsLoading(true);
        setError(null);

        try {
            const response = await fetch("http://localhost:8088/products", {
                headers: {
                    "Content-Type": "application/json",
                },
            });

            if (!response.ok) {
                throw new Error("Failed to fetch products");
            }

            const data = await response.json();
            console.log("Fetched products:", data);

            if (!data) {
                setProducts([]);
                return;
            }

            if (Array.isArray(data)) {
                setProducts(data);
                return;
            }

            const productList = Array.isArray(data.products) ? data.products : [];
            setProducts(productList);
        } catch (err: any) {
            console.error("Error fetching products:", err);
            setError(err.message);
        } finally {
            setIsLoading(false);
        }
    };

    // ---------------- TRANSFORM DATA ----------------
    const transformedProducts = (products || []).map((product: any) => ({
        id: product.ID,
        name: product.Name,
        category: product.Category || "Others",
        price: Number(product.Price),
        condition: product.Condition || "Good",
        seller: product.SellerName || `Seller ${product.SellerID}`,
        phone: product.SellerPhone || "Contact seller",
        image: product.ImageUrl || "/placeholder.jpg",
        images: product.Images || [],
        description: product.Description,
        stock: Number(product.Stock),
    }));

    const allProducts = [...transformedProducts, ...postedItems];

    const filteredProducts = allProducts.filter((product) => {
        const matchesSearch =
            !searchQuery ||
            product.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            (product.category &&
                product.category.toLowerCase().includes(searchQuery.toLowerCase()));

        const matchesCategory =
            !selectedCategory || product.category === selectedCategory;

        return matchesSearch && matchesCategory;
    });

    // ---------------- UI STATES ----------------
    if (isLoading) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="text-center">
                    <div className="w-12 h-12 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
                    <p className="text-gray-600">Loading products...</p>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
                <h3 className="text-lg font-semibold text-red-900 mb-2">
                    Error Loading Products
                </h3>
                <p className="text-red-700 mb-4">{error}</p>
                <button
                    onClick={fetchProducts}
                    className="bg-red-600 text-white px-6 py-2 rounded-lg font-semibold hover:bg-red-700 transition-colors"
                >
                    Try Again
                </button>
            </div>
        );
    }

    // ---------------- MAIN RETURN ----------------
    return (
        <div>
            <div className="flex items-center justify-between mb-6">
                <h3 className="text-xl font-bold text-gray-900">Available Items</h3>
                <button
                    onClick={fetchProducts}
                    className="text-blue-600 hover:text-blue-700 font-semibold text-sm flex items-center gap-2"
                >
                    ⟳ Refresh
                </button>
            </div>

            {filteredProducts.length === 0 ? (
                <div className="text-center py-12">
                    <p className="text-gray-500 text-lg">
                        No items found matching your search.
                    </p>
                </div>
            ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 md:gap-6">
                    {filteredProducts.map((product, index) => (
                        <ProductCard
                            key={`${product.id}-${index}`}
                            product={product}
                            onAddToCart={onAddToCart}
                        />
                    ))}
                </div>
            )}
        </div>
    );
}
