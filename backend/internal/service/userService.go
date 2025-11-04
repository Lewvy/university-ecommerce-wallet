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

func NewUserService(logger *slog.Logger, store data.UserStore) *UserService {
	return &UserService{
		Logger: logger,
		Store:  store,
	}
}

func (s UserService) FindUserByEmail(email string) (*domain.User, error) {
	return nil, nil
}

func (s UserService) Login(input dto.UserLogin) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := s.Store.GetUserAuthByEmail(ctx, input.Email)

	s.Logger.Info("user logged in", "user", map[string]any{"name": user.Name})
	if err != nil {
		return err
	}

	return nil

}

func (s UserService) Signup(ctx context.Context, input dto.UserSignup) (*domain.User, error) {
	err := validateUser(input.Email, input.Password, input.Phone)
	if err != nil {
		return nil, err
	}

	password_hash, err := password.GeneratePasswordHash(input.Password)

	if err != nil {
		return nil, err
	}

	user := db.CreateUserParams{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: password_hash,
	}

	dbUser, err := s.Store.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	resUser := &domain.User{
		Email: dbUser.Email,
		Name:  dbUser.Name,
	}

	s.Logger.Info("User created successfully", "user_id", resUser.ID, "user_email", resUser.Email)

	return resUser, nil

}
