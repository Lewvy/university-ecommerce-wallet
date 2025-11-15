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
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

var ErrPwdMismatch = errors.New("invalid email or password")

type UserService struct {
	Logger       *slog.Logger
	Store        data.UserStore
	WalletStore  data.WalletStore
	Cache        cache.Cache
	Pool         *pgxpool.Pool
	TokenService *TokenService
}

type UserVerification struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

func NewUserService(
	logger *slog.Logger,
	store data.UserStore,
	walletStore data.WalletStore,
	cache cache.Cache,
	pool *pgxpool.Pool,
	tokenService *TokenService,
) *UserService {
	return &UserService{
		Logger:       logger,
		Store:        store,
		Cache:        cache,
		TokenService: tokenService,
		WalletStore:  walletStore,
		Pool:         pool,
	}
}

func (s *UserService) UpdateEmail(ctx context.Context, input dto.UserEmailUpdate) error {
	return nil
}

func (s *UserService) VerifyUser(ctx context.Context, input *UserVerification) error {
	tokenHash, err := s.Cache.GetTokenHashByUserID(ctx, int64(input.ID))
	if err != nil {
		s.Logger.Error("Error getting token", "error", err)
		return err
	}
	match, err := token.MatchToken(input.Token, tokenHash)
	if err != nil {
		s.Logger.Error("Error decoding token", "error", err)
		return err
	}

	if match {
		err := s.Store.VerifyUserEmail(ctx, input.ID)
		if err != nil {
			s.Logger.Error("Error updating email verification", "error", err, "user_id", input.ID)
			return err
		}
		return nil
	} else {
		return fmt.Errorf("invalid token")
	}
}

func (s UserService) FindUserByEmail(email string) (*domain.User, error) {
	return nil, nil
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

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		s.Logger.Error("Failed to begin transaction", "error", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	txUserStore := s.Store.WithTx(tx)
	txWalletStore := s.WalletStore.WithTx(tx)

	pgTextPwdHash := data.NewPGText(password_hash)
	user := db.CreateUserParams{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: pgTextPwdHash,
	}
	s.Logger.Info("Creating User", "email", user.Email)

	var dbUser db.User
	dbUser, err = txUserStore.CreateUser(ctx, user)
	if err != nil {
		s.Logger.Warn("Failed to create user, rolling back", "error", err)
		return nil, err
	}

	_, err = txWalletStore.CreateWallet(ctx, dbUser.ID)
	if err != nil {
		s.Logger.Warn("Failed to create wallet, rolling back", "user_id", dbUser.ID, "error", err)
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		s.Logger.Error("Failed to commit transaction", "error", err)
		return nil, err
	}

	resUser := &domain.User{
		ID:    uint64(dbUser.ID),
		Email: dbUser.Email,
		Name:  dbUser.Name,
	}

	s.sendToken(ctx, dbUser.ID, dbUser.Email, dbUser.Name)

	s.Logger.Info("User and wallet created successfully", "user_id", resUser.ID, "user_email", resUser.Email)

	return resUser, nil
}

func (s *UserService) Login(ctx context.Context, input dto.UserLogin) (accessToken string, refreshToken string, err error) {
	s.Logger.Info("Attempting user login", "email", input.Email)

	userAuth, err := s.Store.GetUserAuthByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return "", "", data.ErrRecordNotFound
		}
		s.Logger.Error("Error retrieving user auth data", "error", err)
		return "", "", err
	}

	stored_pwd_hash := userAuth.PasswordHash.String
	match, err := password.ComparePasswordAndHash(input.Password, stored_pwd_hash)
	if err != nil || !match {
		if !match {
			s.Logger.Warn("Login failed: password mismatch", "email", input.Email)
		} else {
			s.Logger.Error("Password comparison failed", "error", err)
		}
		return "", "", ErrPwdMismatch
	}

	userID := int64(userAuth.ID)

	newAccessToken, newRefreshToken, err := s.TokenService.CreateNewTokens(ctx, userID)
	if err != nil {
		s.Logger.Error("Failed to generate and save tokens", "error", err, "user_id", userID)
		return "", "", errors.New("failed to generate secure tokens")
	}

	s.Logger.Info("User logged in successfully", "user_id", userID)
	return newAccessToken.Plaintext, newRefreshToken.Plaintext, nil
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
	s.Logger.Info("sending token", "user_email", email, "token", token.Plaintext)
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
