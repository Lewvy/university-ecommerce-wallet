"use client"

interface ProfilePageProps {
  userData: {
    username: string
    email: string
    phone: string
  }
  postedItems: any[]
}

export default function ProfilePage({ userData, postedItems }: ProfilePageProps) {
  return (
    <div className="max-w-4xl mx-auto">
      <h3 className="text-2xl font-bold text-gray-900 mb-6">My Profile</h3>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 md:gap-6 mb-8">
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-600 mb-1">Username</p>
          <p className="text-lg md:text-xl font-bold text-gray-900 break-all">{userData.username}</p>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-600 mb-1">Phone Number</p>
          <p className="text-lg md:text-xl font-bold text-gray-900 break-all">{userData.phone}</p>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-600 mb-1">Email</p>
          <p className="text-lg md:text-xl font-bold text-gray-900 break-all">{userData.email}</p>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow p-6 mb-8">
        <button className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 font-medium text-sm">
          Update Profile
        </button>
      </div>

      <div>
        <h4 className="text-xl font-bold text-gray-900 mb-4">Your Posted Items ({postedItems.length})</h4>
        {postedItems.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-12 text-center">
            <p className="text-gray-500">You haven't posted any items yet.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-4">
            {postedItems.map((item, index) => (
              <div key={index} className="bg-white rounded-lg shadow p-4">
                <h5 className="font-semibold text-gray-900 mb-2">{item.productName}</h5>
                <p className="text-sm text-gray-600 mb-2">{item.category}</p>
                <p className="text-lg font-bold text-blue-600">₹{item.price.toLocaleString()}</p>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
