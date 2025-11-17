"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";

export default function ProductDetailsPage() {
    const { id } = useParams();
    const [product, setProduct] = useState<any>(null);
    const [loading, setLoading] = useState(true);
    const [imageIndex, setImageIndex] = useState(0);

    useEffect(() => {
        fetchProduct();
    }, []);

    const fetchProduct = async () => {
        try {
            const token = localStorage.getItem("access_token");

            const res = await fetch(`http://localhost:8088/products/${id}`, {
                headers: {
                    "Authorization": token ? `Bearer ${token}` : "",
                    "Content-Type": "application/json",
                },
            });

            if (!res.ok) {
                console.error("API Error:", res.status);
                setProduct(null);
                setLoading(false);
                return;
            }

            const data = await res.json();
            console.log("Loaded Product:", data);
            setProduct(data);

        } catch (err) {
            console.error("Failed to load product:", err);
        } finally {
            setLoading(false);
        }
    };

    if (loading) return <div className="p-10 text-center">Loading...</div>;
    if (!product) return <div className="p-10 text-center">Product not found</div>;

    const images = product.Images?.length > 0 ? product.Images : [product.ImageUrl];

    return (
        <div className="max-w-5xl mx-auto p-6">
            <button
                onClick={() => window.history.back()}
                className="mb-6 bg-gray-200 text-gray-800 px-4 py-2 rounded-lg hover:bg-gray-300"
            >
                ← Back
            </button>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">

                <div>
                    <img
                        src={images[imageIndex]}
                        className="w-full h-96 object-contain border rounded-lg"
                    />

                    <div className="flex gap-3 mt-4 overflow-x-auto">
                        {images.map((img: string, idx: number) => (
                            <img
                                key={idx}
                                src={img}
                                className={`w-20 h-20 object-cover rounded-lg cursor-pointer border-2 ${
                                    idx === imageIndex ? "border-blue-600" : "border-gray-300"
                                }`}
                                onClick={() => setImageIndex(idx)}
                            />
                        ))}
                    </div>
                </div>

                {/* RIGHT — DETAILS */}
                <div>
                    <h1 className="text-3xl font-bold mb-2">{product.Name}</h1>
                    <p className="text-gray-600 mb-4">{product.Category}</p>

                    <p className="text-blue-600 text-2xl font-bold mb-4">
                        ₹{Number(product.Price).toLocaleString()}
                    </p>

                    <p className="text-gray-700 mb-4">
                        <span className="font-semibold">Condition:</span> {product.Condition}
                    </p>

                    <p className="text-gray-700 mb-4">
                        <span className="font-semibold">Stock:</span> {product.Stock}
                    </p>

                    <h2 className="text-xl font-semibold mt-6 mb-2">Description</h2>
                    <p className="text-gray-800 whitespace-pre-line">{product.Description}</p>

                    <h2 className="text-xl font-semibold mt-6 mb-2">Seller Details</h2>
                    <div className="bg-gray-100 p-4 rounded-lg">
                        <p><strong>Name:</strong> {product.SellerName || "Unknown Seller"}</p>
                        <p><strong>Phone:</strong> {product.SellerPhone || "N/A"}</p>
                        <p><strong>Seller ID:</strong> {product.SellerID}</p>
                    </div>

                    <button className="mt-6 w-full bg-blue-600 text-white text-lg py-3 rounded-lg hover:bg-blue-700">
                        Add to Cart
                    </button>
                </div>

            </div>
        </div>
    );
}
