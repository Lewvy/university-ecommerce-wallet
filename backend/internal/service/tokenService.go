package service

import (
	"context"
	"ecommerce/internal/data"
	db "ecommerce/internal/data/gen"
	"ecommerce/internal/token"
	"log/slog"

	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

type TokenService struct {
	Store  data.TokenStore
	Logger *slog.Logger
}

func NewTokenService(store data.TokenStore, logger *slog.Logger) *TokenService {
	return &TokenService{Store: store, Logger: logger}
}

func (s *TokenService) CreateNewTokens(ctx context.Context, userID int64) (*token.Token, *token.Token, error) {

	s.Logger.Info("Generating new token pair", "user_id", userID)

	accessToken, err := token.GenerateAccessToken(userID, AccessTokenTTL, token.ScopeAuthentication)
	if err != nil {
		s.Logger.Error("Failed to generate access token", "user_id", userID, "error", err)
		return nil, nil, err
	}

	refreshToken, err := token.GenerateRefreshToken(userID, RefreshTokenTTL, token.ScopeRefresh)
	if err != nil {
		s.Logger.Error("Failed to generate refresh token", "user_id", userID, "error", err)
		return nil, nil, err
	}

	if err = s.InsertToken(ctx, refreshToken); err != nil {
		s.Logger.Error("Failed to insert refresh token hash into DB", "user_id", userID, "error", err)
		return nil, nil, err
	}
	s.Logger.Info("Successfully persisted refresh token hash", "user_id", userID)

	return accessToken, refreshToken, nil
}

func (s *TokenService) InsertToken(ctx context.Context, t *token.Token) error {
	hashBytes, err := hex.DecodeString(token.GenerateTokenHash(t.Plaintext))
	if err != nil {
		s.Logger.Error("Failed to decode token hash for insertion", "user_id", t.UserID, "error", err)
		return err
	}

	params := db.InsertTokenParams{
		Hash:   hashBytes,
		UserID: t.UserID,
		Expiry: pgtype.Timestamptz{
			Time:  t.Expiry,
			Valid: true,
		},
		Scope: t.Scope,
	}

	s.Logger.Debug("Attempting to insert token into DB", "user_id", t.UserID, "scope", t.Scope)
	return s.Store.InsertToken(ctx, params)
}

func (s *TokenService) RefreshAndRevokeTokens(ctx context.Context, oldPlaintextToken string) (*token.Token, *token.Token, error) {
	s.Logger.Info("Starting token refresh attempt")

	hashHex := token.GenerateTokenHash(oldPlaintextToken)
	hashBytes, err := hex.DecodeString(hashHex)
	if err != nil {
		s.Logger.Warn("Received token with invalid hex format", "error", err)
		return nil, nil, errors.New("invalid token format provided")
	}

	tokenRecord, err := s.Store.GetTokenByHash(ctx, hashBytes)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			s.Logger.Warn("Refresh token not found in DB (revoked or never existed)")
			return nil, nil, token.ErrTokenNotFound
		}
		s.Logger.Error("Database error during token lookup", "error", err)
		return nil, nil, err
	}

	s.Logger.Debug("Token record retrieved", "user_id", tokenRecord.UserID, "scope", tokenRecord.Scope)

	if tokenRecord.Scope != token.ScopeRefresh {
		s.Logger.Warn("Token used for refresh has wrong scope", "expected_scope", token.ScopeRefresh, "actual_scope", tokenRecord.Scope, "user_id", tokenRecord.UserID)

		if err := s.Store.DeleteAllForUserAndScope(ctx, tokenRecord.Scope, tokenRecord.UserID); err != nil {
			s.Logger.Error("CRITICAL: Failed to revoke misused token!", "user_id", tokenRecord.UserID, "error", err)
			return nil, nil, errors.New("internal server error: failed to revoke misused token")
		}
		s.Logger.Info("Successfully revoked misused token", "user_id", tokenRecord.UserID)

		return nil, nil, token.ErrTokenNotFound
	}

	if tokenRecord.Expiry.Valid && time.Now().After(tokenRecord.Expiry.Time) {
		s.Logger.Info("Refresh token expired", "user_id", tokenRecord.UserID, "expiry_time", tokenRecord.Expiry.Time)

		if err := s.Store.DeleteAllForUserAndScope(ctx, tokenRecord.Scope, tokenRecord.UserID); err != nil {
			s.Logger.Error("CRITICAL: Failed to revoke expired token!", "user_id", tokenRecord.UserID, "error", err)
			return nil, nil, err
		}
		s.Logger.Info("Successfully revoked expired token", "user_id", tokenRecord.UserID)

		return nil, nil, errors.New("refresh token expired")
	}

	if !tokenRecord.Expiry.Valid {
		s.Logger.Error("DB error: Token expiry time is NULL/invalid", "user_id", tokenRecord.UserID)
		return nil, nil, errors.New("token expiry time is invalid in database")
	}

	s.Logger.Info("Revoking old refresh token", "user_id", tokenRecord.UserID)
	if err = s.Store.DeleteAllForUserAndScope(ctx, tokenRecord.Scope, tokenRecord.UserID); err != nil {
		s.Logger.Error("CRITICAL: Failed to revoke valid token!", "user_id", tokenRecord.UserID, "error", err)
		return nil, nil, err
	}

	s.Logger.Info("Old refresh token successfully revoked", "user_id", tokenRecord.UserID)

	newAccessToken, newRefreshToken, err := s.CreateNewTokens(ctx, tokenRecord.UserID)
	if err != nil {
		return nil, nil, err
	}

	s.Logger.Info("Token pair successfully refreshed", "user_id", tokenRecord.UserID)

	return newAccessToken, newRefreshToken, nil
}

func (s *TokenService) RevokeAllUserTokens(ctx context.Context, scope string, userID int64) error {
	s.Logger.Info("Attempting to revoke all tokens for user", "user_id", userID, "scope", scope)
	err := s.Store.DeleteAllForUserAndScope(ctx, scope, userID)
	if err != nil {
		s.Logger.Error("Failed to revoke user tokens", "user_id", userID, "scope", scope, "error", err)
	} else {
		s.Logger.Info("Successfully revoked tokens", "user_id", userID, "scope", scope)
	}
	return err
}
