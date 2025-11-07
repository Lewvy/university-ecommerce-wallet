package service

import (
	"context"
	"ecommerce/domain"
	"ecommerce/internal/cache"
	"ecommerce/internal/data"
	db "ecommerce/internal/data/gen"
	"ecommerce/internal/dto"
	"ecommerce/internal/password"
	"ecommerce/internal/token"
	"ecommerce/internal/worker"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

type UserService struct {
	Logger *slog.Logger
	Store  data.UserStore
	Cache  cache.Cache
}

func (s *UserService) VerifyUser(input struct {
	Email string "json:\"email\""
	Token string "json:\"token\""
}) {

}

func NewUserService(logger *slog.Logger, store data.UserStore, cache cache.Cache) *UserService {
	return &UserService{
		Logger: logger,
		Store:  store,
		Cache:  cache,
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
	var password_hash []byte

	err := validateUser(input.Email, input.Password, input.Phone)
	if err != nil {
		return nil, err
	}
	password_hash, err = password.GeneratePasswordHash(input.Password)
	if err != nil {
		return nil, err
	}

	user := db.CreateUserParams{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: password_hash,
	}

	var dbUser db.User
	// if os.Getenv("ENABLE_FAST_VALIDATION") == "true" {
	// 	s.Logger.Info("User created successfully", "user_email", user.Email)
	//
	// 	resUser := &domain.User{
	// 		Email: user.Email,
	// 		Name:  user.Name,
	// 	}
	// 	return resUser, nil
	// }
	dbUser, err = s.Store.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	resUser := &domain.User{
		ID:    uint64(dbUser.ID),
		Email: dbUser.Email,
		Name:  dbUser.Name,
	}

	s.sendToken(ctx, dbUser.ID, dbUser.Email, dbUser.Name)

	s.Logger.Info("User created successfully", "user_id", resUser.ID, "user_email", resUser.Email)

	return resUser, nil

}

func (s UserService) sendToken(ctx context.Context, id int32, email string, name string) {

	expiry := time.Minute * 15
	token, err := token.GenerateVerificationToken(int64(id), expiry, token.ScopeActivation)
	if err != nil {
		s.Logger.Error("Error generating token", "error", err)
		return
	}

	err = s.Cache.SetVerificationToken(ctx, token.Hash, int64(id), expiry)
	if err != nil {
		s.Logger.Error("Error saving token to cache", "error", err)
		return
	}
	if err := s.queueVerificationEmail(ctx, email, name, id, token.Plaintext); err != nil {
		s.Logger.Error("Error saving token to cache", "error", err)
		return
	}
}

type VerificationData struct {
	ID    int32  `json:"ID"`
	Token string `json:"activationToken"`
	Name  string `json:"name"`
}

func (s UserService) queueVerificationEmail(ctx context.Context, userEmail, name string, id int32, tokenPlaintext string) error {
	data := VerificationData{
		ID:    id,
		Token: tokenPlaintext,
		Name:  name,
	}

	job := worker.MailJob{
		Recipient:    userEmail,
		TemplateFile: "user_templates.tmpl",
		TemplateData: data,
	}

	jobJSON, err := json.Marshal(job)
	if err != nil {
		s.Logger.Error("Failed to serialize mail job", "error", err)
		return fmt.Errorf("failed to serialize mail job: %w", err)
	}

	err = s.Cache.AddEmailToQueue(ctx, userEmail, string(jobJSON))
	if err != nil {
		s.Logger.Error("Failed to LPUSH mail job to Valkey", "error", err)
		return fmt.Errorf("failed to enqueue email job to valkey: %w", err)
	}

	return nil
}
