"use client"

import { useState } from "react"
import SignupPage from "@/components/signup-page"
import HomePage from "@/components/home-page"

export default function Page() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [userData, setUserData] = useState({
    username: "",
    email: "",
    phone: "",
  })

  if (!isAuthenticated) {
    return (
      <SignupPage
        onSignupSuccess={(name, email, phone) => {
          setUserData({ username: name, email, phone })
          setIsAuthenticated(true)
        }}
      />
    )
  }

  return <HomePage userData={userData} />
}
