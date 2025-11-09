package service

import (
	"context"
	"ecommerce/internal/data"
	"ecommerce/internal/dto"
	"ecommerce/internal/password"
	"errors"
	"log/slog"
)

type AuthService struct {
	Logger       *slog.Logger
	UserStore    data.UserStore
	TokenService *TokenService
}

func NewAuthService(logger *slog.Logger, userStore data.UserStore, tokenService *TokenService) *AuthService {
	return &AuthService{
		Logger:       logger,
		UserStore:    userStore,
		TokenService: tokenService,
	}
}

func (s *AuthService) Login(ctx context.Context, input dto.UserLogin) (accessToken string, refreshToken string, err error) {
	s.Logger.Info("Attempting user login", "email", input.Email)

	userAuth, err := s.UserStore.GetUserAuthByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return "", "", ErrPwdMismatch
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
