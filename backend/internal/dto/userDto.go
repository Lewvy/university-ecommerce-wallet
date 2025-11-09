package dto

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserEmailUpdate struct {
	ID           int    `json:"id"`
	UpdatedEmail string `json:"updated_email"`
}

type UserSignup struct {
	UserLogin
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
