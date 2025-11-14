package data

import (
	"context"
	"database/sql"
	db "ecommerce/internal/data/gen"
	"errors"
)

type WalletStore interface {
	CreateWallet(ctx context.Context, userID int32) (db.Wallet, error)
	GetWalletByUserID(ctx context.Context, userID int32) (db.Wallet, error)
	GetWalletByUserIDForUpdate(ctx context.Context, userID int32) (db.Wallet, error)
	CreditWallet(ctx context.Context, arg db.CreditWalletParams) (db.Wallet, error)
	DebitWallet(ctx context.Context, arg db.DebitWalletParams) (db.Wallet, error)
}

type sqlWalletStore struct {
	q *db.Queries
}

func NewWalletStore(queries *db.Queries) WalletStore {
	return &sqlWalletStore{
		q: queries,
	}
}

func (s *sqlWalletStore) CreateWallet(ctx context.Context, userID int32) (db.Wallet, error) {
	return s.q.CreateWallet(ctx, int32(userID))
}

func (s *sqlWalletStore) GetWalletByUserID(ctx context.Context, userID int32) (db.Wallet, error) {
	wallet, err := s.q.GetBalanceById(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return db.Wallet{}, ErrRecordNotFound
	}
	return wallet, err
}

func (s *sqlWalletStore) GetWalletByUserIDForUpdate(ctx context.Context, userID int32) (db.Wallet, error) {
	wallet, err := s.q.GetWalletByUserIDForUpdate(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return db.Wallet{}, ErrRecordNotFound
	}
	return wallet, err
}

func (s *sqlWalletStore) CreditWallet(ctx context.Context, arg db.CreditWalletParams) (db.Wallet, error) {
	return s.q.CreditWallet(ctx, arg)
}

func (s *sqlWalletStore) DebitWallet(ctx context.Context, arg db.DebitWalletParams) (db.Wallet, error) {
	return s.q.DebitWallet(ctx, arg)
}
