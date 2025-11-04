package service

import (
	"context"
	"ecommerce/domain"
	db "ecommerce/internal/data"
	"ecommerce/internal/dto"
	"ecommerce/internal/password"
	"fmt"
	"log/slog"
	"time"
)

type UserService struct {
	Logger  *slog.Logger
	Queries *db.Queries
}

func (s UserService) FindUserByEmail(email string) (*domain.User, error) {
	return nil, nil
}

func (s UserService) Login(input any) error {
	switch v := input.(type) {
	case dto.UserLogin:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user, err := s.Queries.GetUserAuthByEmail(ctx, v.Email)

		s.Logger.Info("user logged in", "user", map[string]any{"name": user.Name, "password_hash": string(user.PasswordHash)})
		if err != nil {
			return err
		}

		return nil
	default:
		return fmt.Errorf("Signup: Unknown data type: %T", v)
	}

}

func (s UserService) Signup(input any) (map[string]string, error) {
	switch v := input.(type) {

	case dto.UserSignup:
		validationErrors := domain.ValidateUser(v.Email, v.Password, v.Phone)
		if len(validationErrors) > 0 {
			return validationErrors, nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		password_hash, err := password.GeneratePasswordHash(v.Password)
		s.Logger.Info("Password hash generated", "password", password_hash)
		if err != nil {
			return nil, err
		}
		user := db.CreateUserParams{
			Name:         v.Name,
			Email:        v.Email,
			PasswordHash: password_hash,
		}
		s.Logger.Info("creating user", "user", user)

		dbUser, err := s.Queries.CreateUser(ctx, user)
		if err != nil {
			return nil, err
		}
		s.Logger.Info("User created successfully", "user", dbUser)

		return nil, nil

	default:
		return nil, fmt.Errorf("Signup: Unknown data type: %T", v)
	}
}
