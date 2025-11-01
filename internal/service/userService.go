package service

import "ecommerce/domain"

type UserService struct {
}

func (s UserService) FindUserByEmail(email string) (*domain.User, error) {
	return nil, nil
}

func (s UserService) Signup(input any) (string, error) {
	user := struct {
		email    string
		password string
		name     string
	}{}
}
