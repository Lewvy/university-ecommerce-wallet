"use client";

import { useEffect, useState } from "react";

interface PostedItem {
    ID: number;
    SellerID: number;
    Name: string;
    Description: string;
    Price: number;
    Stock: number;
    ImageUrl: string | null;
    CreatedAt: string;
    UpdatedAt: string;
}

interface ProfilePageProps {
    userData: {
        username: string;
        email: string;
        phone: string;
        id?: number;
    };
}

export default function ProfilePage({ userData }: ProfilePageProps) {
    const [postedItems, setPostedItems] = useState<PostedItem[]>([]);
    const [loading, setLoading] = useState(true);
    const [username, setUsername] = useState(userData.username);
    const [phone, setPhone] = useState(userData.phone);

    useEffect(() => {
        fetchUserPostedItems();
    }, []);

    // Fetch products posted by this user
    const fetchUserPostedItems = async () => {
        try {
            const token = localStorage.getItem("access_token");
            if (!token) {
                console.error("No token found");
                return;
            }

            const response = await fetch("http://localhost:8088/products/mine", {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            if (!response.ok) {
                console.error("Failed to fetch posted items");
                return;
            }

            const data = await response.json();
            console.log("Fetched posted items:", data);

            setPostedItems(data || []);
        } catch (error) {
            console.error("Error fetching posted items:", error);
        } finally {
            setLoading(false);
        }
    };

    // Update username + phone number
    const handleUpdateProfile = async () => {
        try {
            const token = localStorage.getItem("access_token");
            if (!token) return;

            const response = await fetch("http://localhost:8088/user/update", {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify({
                    username,
                    phone,
                }),
            });

            if (response.ok) {
                alert("Profile updated successfully!");
            } else {
                alert("Failed to update profile");
            }
        } catch (error) {
            console.error("Profile update error:", error);
            alert("Network error");
        }
    };

    return (
        <div className="p-6">
            <h1 className="text-3xl font-bold text-gray-900 mb-6">My Profile</h1>

            {/* Profile Info */}
            <div className="grid grid-cols-3 gap-6 mb-6">
                <div className="bg-white rounded-lg shadow p-4">
                    <p className="text-sm text-gray-500">Username</p>
                    <input
                        className="mt-1 w-full p-2 border rounded"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                    />
                </div>

                <div className="bg-white rounded-lg shadow p-4">
                    <p className="text-sm text-gray-500">Phone Number</p>
                    <input
                        className="mt-1 w-full p-2 border rounded"
                        value={phone}
                        onChange={(e) => setPhone(e.target.value)}
                    />
                </div>

                <div className="bg-white rounded-lg shadow p-4">
                    <p className="text-sm text-gray-500">Email</p>
                    <p className="mt-1 font-semibold">{userData.email}</p>
                </div>
            </div>

            <button
                onClick={handleUpdateProfile}
                className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition mb-8"
            >
                Update Profile
            </button>

            {/* Posted Items */}
            <h2 className="text-xl font-semibold mb-4">
                Your Posted Items ({postedItems.length})
            </h2>

            {loading ? (
                <p>Loading...</p>
            ) : postedItems.length === 0 ? (
                <p className="text-gray-500">You have not posted any items yet.</p>
            ) : (
                <div className="grid grid-cols-3 gap-6">
                    {postedItems.map((item) => (
                        <div key={item.ID} className="bg-white rounded-lg shadow p-4">

                            {/* Product Image */}
                            {item.ImageUrl ? (
                                <img
                                    src={item.ImageUrl}
                                    alt={item.Name}
                                    className="w-full h-40 object-cover rounded-md mb-3"
                                />
                            ) : (
                                <div className="w-full h-40 bg-gray-200 rounded-md flex items-center justify-center mb-3">
                                    <span className="text-gray-500 text-sm">No Image</span>
                                </div>
                            )}

                            <h5 className="font-semibold text-gray-900 mb-2">{item.Name}</h5>
                            <p className="text-sm text-gray-600 mb-2">{item.Description}</p>

                            <p className="text-lg font-bold text-blue-600">
                                ₹{item.Price.toLocaleString()}
                            </p>

                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
