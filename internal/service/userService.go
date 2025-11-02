package service

import (
	"database/sql"
	"ecommerce/domain"
	"ecommerce/internal/dto"
	"fmt"
)

type UserService struct {
	svc *UserService
	db  *sql.DB
}

func (s UserService) FindUserByEmail(email string) (*domain.User, error) {
	return nil, nil
}

func (s UserService) Signup(input any) (map[string]string, error) {
	switch v := input.(type) {

	case dto.UserSignup:
		validationErrors := domain.ValidateUser(v.Email, v.Password, v.Phone)
		if validationErrors != nil {
			return validationErrors, nil
		}
		return nil, nil

	default:
		return nil, fmt.Errorf("Signup: Unknown data type: %T", v)
	}
}
