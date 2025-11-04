package dto

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignup struct {
	UserLogin
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
