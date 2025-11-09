package data

import (
	"context"
	"database/sql"
	db "ecommerce/internal/data/gen"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type TokenStore interface {
	InsertToken(ctx context.Context, arg db.InsertTokenParams) error
	DeleteAllForUserAndScope(ctx context.Context, scope string, userID int64) error
	GetTokenByHash(ctx context.Context, hash []byte) (db.Token, error)
}

type sqlTokenStore struct {
	q *db.Queries
}

func NewTokenStore(queries *db.Queries) TokenStore {
	return &sqlTokenStore{
		q: queries,
	}
}

func (s *sqlTokenStore) InsertToken(ctx context.Context, arg db.InsertTokenParams) error {
	return s.q.InsertToken(ctx, arg)
}

func (s *sqlTokenStore) DeleteAllForUserAndScope(ctx context.Context, scope string, userID int64) error {
	params := db.DeleteTokenParams{
		Scope:  scope,
		UserID: userID,
	}
	return s.q.DeleteToken(ctx, params)
}

func (s *sqlTokenStore) GetTokenByHash(ctx context.Context, hash []byte) (db.Token, error) {
	token, err := s.q.GetTokenByHash(ctx, hash)
	if errors.Is(err, sql.ErrNoRows) {
		return db.Token{}, ErrRecordNotFound
	}
	return token, err
}
