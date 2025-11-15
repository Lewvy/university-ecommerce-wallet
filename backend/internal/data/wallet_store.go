package data

import (
	"context"
	"database/sql"
	db "ecommerce/internal/data/gen"
	"errors"

	"github.com/jackc/pgx/v5"
)

type WalletStore interface {
	GetWalletByUserID(ctx context.Context, userID int64) (db.Wallet, error)
	GetWalletByUserIDForUpdate(ctx context.Context, userID int32) (db.Wallet, error)
	CreateWallet(ctx context.Context, userID int32) (db.Wallet, error)
	CreditWallet(ctx context.Context, arg db.CreditWalletParams) (db.Wallet, error)
	DebitWallet(ctx context.Context, arg db.DebitWalletParams) (db.Wallet, error)
	WithTx(tx pgx.Tx) WalletStore

	CreateTransaction(ctx context.Context, arg db.CreateTransactionParams) (db.WalletTransaction, error)
	GetTransactionByOrderID(ctx context.Context, rzpOrderID string) (db.WalletTransaction, error)
	UpdateTransactionOrderID(ctx context.Context, arg db.UpdateTransactionOrderIDParams) (db.WalletTransaction, error)
	UpdateTransactionStatus(ctx context.Context, arg db.UpdateTransactionStatusParams) (db.WalletTransaction, error)
}

type sqlWalletStore struct {
	q *db.Queries
}

func NewWalletStore(queries *db.Queries) WalletStore {
	return &sqlWalletStore{
		q: queries,
	}
}
func (s *sqlWalletStore) WithTx(tx pgx.Tx) WalletStore {
	return &sqlWalletStore{
		q: db.New(tx),
	}
}

func (s *sqlWalletStore) CreateWallet(ctx context.Context, userID int32) (db.Wallet, error) {
	return s.q.CreateWallet(ctx, int32(userID))
}

func (s *sqlWalletStore) GetWalletByUserID(ctx context.Context, userID int64) (db.Wallet, error) {
	wallet, err := s.q.GetWalletByUserID(ctx, int32(userID))
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return db.Wallet{}, ErrRecordNotFound
	}
	return wallet, err
}

func (s *sqlWalletStore) GetWalletByUserIDForUpdate(ctx context.Context, userID int32) (db.Wallet, error) {
	wallet, err := s.q.GetWalletByUserID(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
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

func (s *sqlWalletStore) CreateTransaction(ctx context.Context, arg db.CreateTransactionParams) (db.WalletTransaction, error) {
	return s.q.CreateTransaction(ctx, arg)
}

func (s *sqlWalletStore) GetTransactionByOrderID(ctx context.Context, rzpOrderID string) (db.WalletTransaction, error) {
	tx, err := s.q.GetTransactionByOrderID(ctx, NewPGText(rzpOrderID))
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return db.WalletTransaction{}, ErrRecordNotFound
	}
	return tx, err
}

func (s *sqlWalletStore) UpdateTransactionOrderID(ctx context.Context, arg db.UpdateTransactionOrderIDParams) (db.WalletTransaction, error) {
	tx, err := s.q.UpdateTransactionOrderID(ctx, arg)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return db.WalletTransaction{}, ErrRecordNotFound
	}
	return tx, err
}

func (s *sqlWalletStore) UpdateTransactionStatus(ctx context.Context, arg db.UpdateTransactionStatusParams) (db.WalletTransaction, error) {
	tx, err := s.q.UpdateTransactionStatus(ctx, arg)
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return db.WalletTransaction{}, ErrRecordNotFound
	}
	return tx, err
}
