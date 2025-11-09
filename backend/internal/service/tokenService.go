package service

import (
	"context"
	"ecommerce/internal/data"
	db "ecommerce/internal/data/gen"
	"ecommerce/internal/token"

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
	Store data.TokenStore
}

func NewTokenService(store data.TokenStore) *TokenService {
	return &TokenService{Store: store}
}

func (s *TokenService) CreateNewTokens(ctx context.Context, userID int64) (*token.Token, *token.Token, error) {
	accessToken, err := token.GenerateAccessToken(userID, AccessTokenTTL, token.ScopeAuthentication)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := token.GenerateRefreshToken(userID, RefreshTokenTTL, token.ScopeRefresh)
	if err != nil {
		return nil, nil, err
	}

	if err = s.InsertToken(ctx, refreshToken); err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *TokenService) InsertToken(ctx context.Context, t *token.Token) error {
	hashBytes, err := hex.DecodeString(token.GenerateTokenHash(t.Plaintext))
	if err != nil {
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
	return s.Store.InsertToken(ctx, params)
}

func (s *TokenService) RefreshAndRevokeTokens(ctx context.Context, oldPlaintextToken string) (*token.Token, *token.Token, error) {
	hashHex := token.GenerateTokenHash(oldPlaintextToken)
	hashBytes, err := hex.DecodeString(hashHex)
	if err != nil {
		return nil, nil, errors.New("invalid token format provided")
	}

	tokenRecord, err := s.Store.GetTokenByHash(ctx, hashBytes)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return nil, nil, token.ErrTokenNotFound
		}
		return nil, nil, err
	}

	if tokenRecord.Scope != token.ScopeRefresh {
		if err := s.Store.DeleteAllForUserAndScope(ctx, tokenRecord.Scope, tokenRecord.UserID); err != nil {
			return nil, nil, errors.New("internal server error: failed to revoke misused token")
		}

		return nil, nil, token.ErrTokenNotFound
	}

	if tokenRecord.Expiry.Valid && time.Now().After(tokenRecord.Expiry.Time) {
		if err := s.Store.DeleteAllForUserAndScope(ctx, tokenRecord.Scope, tokenRecord.UserID); err != nil {
			return nil, nil, err
		}

		return nil, nil, errors.New("refresh token expired")
	}
	if !tokenRecord.Expiry.Valid {
		return nil, nil, errors.New("token expiry time is invalid in database")
	}

	if err = s.Store.DeleteAllForUserAndScope(ctx, tokenRecord.Scope, tokenRecord.UserID); err != nil {
		return nil, nil, err
	}

	newAccessToken, newRefreshToken, err := s.CreateNewTokens(ctx, tokenRecord.UserID)
	if err != nil {
		return nil, nil, err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *TokenService) RevokeAllUserTokens(ctx context.Context, scope string, userID int64) error {
	return s.Store.DeleteAllForUserAndScope(ctx, scope, userID)
}
