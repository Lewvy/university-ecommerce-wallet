"use client"

import { useEffect } from "react"

export default function ProductModal({ product, onClose }: any) {
  
  if (!product) return null

  const images = product.images?.length > 0 
    ? product.images 
    : [product.image]

  return (
    <div className="fixed inset-0 bg-black bg-opacity-70 z-50 flex items-center justify-center p-4">
      
      {/* Close Button */}
      <button
        onClick={onClose}
        className="absolute top-5 right-5 text-white text-3xl font-bold"
      >
        ✖
      </button>

      <div className="bg-white rounded-lg max-w-2xl w-full p-4 shadow-xl">
        
        {/* IMAGE VIEWER */}
        <div className="flex overflow-x-auto gap-4 scrollbar-hide">
          {images.map((img: string, i: number) => (
            <img
              key={i}
              src={img}
              className="h-96 object-contain rounded-lg"
            />
          ))}
        </div>

        <div className="text-center mt-4">
          <h2 className="text-xl font-semibold">{product.name}</h2>
          <p className="text-gray-600">{product.description}</p>
        </div>

      </div>
    </div>
  )
}
