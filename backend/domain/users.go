package domain

import (
	"time"
)

type User struct {
	ID            uint64    `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Email         string    `json:"email,omitempty"`
	Password_Hash string    `json:"password_hash,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Verified      bool      `json:"verified,omitempty"`
	User_Type     string    `json:"user_type,omitempty"`
}
