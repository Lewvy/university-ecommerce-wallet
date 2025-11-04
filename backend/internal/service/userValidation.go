package service

import (
	"ecommerce/internal/validator"

	"github.com/nyaruka/phonenumbers"
)

func validateUser(email, password, phone string) error {
	v := validator.New()
	validateEmail(email, v)
	validatePassword(password, v)
	validatePhone(phone, v)
	if len(v.Errors) == 0 {
		return nil
	}
	return v
}

func validatePhone(phone string, v *validator.ValidationError) {
	phone_number, err := phonenumbers.Parse(phone, "IN")

	if err != nil {
		v.AddError("phone", err.Error())
		return
	}
	if !phonenumbers.IsValidNumber(phone_number) {
		v.AddError("phone", "add a valid phone number")
	}

}

func validateEmail(email string, v *validator.ValidationError) {
	v.Check(email != "", "email", "must be provided")
	v.Check(v.Matches(email, validator.EmailRX), "email", "invalid email format")
}

func validatePassword(pwd string, v *validator.ValidationError) {
	v.Check(pwd != "", "password", "must be provided")
	v.Check(len(pwd) >= 8, "password", "must be longer than 8 bytes")
	v.Check(len(pwd) <= 50, "password", "must be less than 50 bytes")
}
