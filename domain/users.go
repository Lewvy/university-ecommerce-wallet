package domain

import (
	"time"
)

type User struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Expiry    time.Time `json:"expiry"`
	WalletID  string    `json:"wallet_id"`
	Verified  bool      `json:"verified"`
	UserType  string    `json:"user_type"`
}
