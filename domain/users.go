package domain

import (
	"ecommerce/internal/validator"
	"github.com/nyaruka/phonenumbers"
	"time"
)

type User struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password_Hash string    `json:"password"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Verified      bool      `json:"verified"`
	User_Type     string    `json:"user_type"`
}

func ValidateUser(email, password, phone string) map[string]string {
	v := validator.New()
	validateEmail(email, v)
	validatePassword(password, v)
	validatePhone(phone, v)
	return v.Errors
}

func validatePhone(phone string, v *validator.Validator) {
	phone_number, err := phonenumbers.Parse(phone, "IN")

	if err != nil {
		v.AddError("phone", err.Error())
		return
	}
	if !phonenumbers.IsValidNumber(phone_number) {
		v.AddError("phone", "add a valid phone number")
	}

}

func validateEmail(email string, v *validator.Validator) {
	v.Check(email != "", "email", "must be provided")
	v.Check(v.Matches(email, validator.EmailRX), "email", "invalid email format")
}

func validatePassword(pwd string, v *validator.Validator) {
	v.Check(pwd != "", "password", "must be provided")
	v.Check(len(pwd) >= 8, "password", "must be longer than 8 bytes")
	v.Check(len(pwd) <= 50, "password", "must be less than 50 bytes")

}
