package service

import (
	"context"
	"ecommerce/domain"
	"ecommerce/internal/data"
	db "ecommerce/internal/data/gen"
	"ecommerce/internal/dto"
	"ecommerce/internal/password"
	"ecommerce/internal/validator"
	"log/slog"
	"time"

	"github.com/nyaruka/phonenumbers"
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
