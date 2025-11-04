package service

import (
	"context"
	"ecommerce/domain"
	"ecommerce/internal/data"
	db "ecommerce/internal/data/gen"
	"ecommerce/internal/dto"
	"ecommerce/internal/password"
	"log/slog"
	"time"
)

type UserService struct {
	Logger *slog.Logger
	Store  data.UserStore
}

func (s UserService) FindUserByEmail(email string) (*domain.User, error) {
	return nil, nil
}

func NewUserService(logger *slog.Logger, store data.UserStore) *UserService {
	return &UserService{
		Logger: logger,
		Store:  store,
	}
}

func (s UserService) Login(input dto.UserLogin) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := s.Store.GetUserAuthByEmail(ctx, input.Email)

	s.Logger.Info("user logged in", "user", map[string]any{"name": user.Name, "password_hash": string(user.PasswordHash)})
	if err != nil {
		return err
	}

	return nil

}

func (s UserService) Signup(input dto.UserSignup) (map[string]string, error) {
	validationErrors := domain.ValidateUser(input.Email, input.Password, input.Phone)
	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	password_hash, err := password.GeneratePasswordHash(input.Password)
	s.Logger.Info("Password hash generated", "password", password_hash)
	if err != nil {
		return nil, err
	}
	user := db.CreateUserParams{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: password_hash,
	}
	s.Logger.Info("creating user", "user", user)

	dbUser, err := s.Store.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	s.Logger.Info("User created successfully", "user", dbUser)

	return nil, nil

}
