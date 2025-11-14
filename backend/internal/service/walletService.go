package service

import (
	"context"
	"ecommerce/internal/data"
	db "ecommerce/internal/data/gen"
	"errors"
	"log/slog"
)

type WalletService struct {
	Store  data.WalletStore
	Logger *slog.Logger
	// dbConn *sql.DB // A real-world service might hold the DB connection to manage transactions
}

func NewWalletService(store data.WalletStore, logger *slog.Logger) *WalletService {
	return &WalletService{Store: store, Logger: logger}
}

func (s *WalletService) CreateWallet(ctx context.Context, userID int32) (db.Wallet, error) {
	s.Logger.Info("Attempting to create wallet", "user_id", userID)
	wallet, err := s.Store.CreateWallet(ctx, userID)
	if err != nil {
		s.Logger.Error("Failed to create wallet", "user_id", userID, "error", err)
	} else {
		s.Logger.Info("Successfully created wallet", "user_id", userID, "wallet_id", wallet.UserID)
	}
	return wallet, err
}

func (s *WalletService) GetWalletByUserID(ctx context.Context, userID int32) (db.Wallet, error) {
	s.Logger.Debug("Fetching wallet details", "user_id", userID)
	wallet, err := s.Store.GetWalletByUserID(ctx, userID)

	if errors.Is(err, data.ErrRecordNotFound) {
		s.Logger.Warn("Wallet not found for user", "user_id", userID)
		return db.Wallet{}, data.ErrRecordNotFound
	}

	if err != nil {
		s.Logger.Error("Database error during wallet lookup", "user_id", userID, "error", err)
		return db.Wallet{}, err
	}

	return wallet, nil
}

func (s *WalletService) Credit(ctx context.Context, userID int32, amount int64) (db.Wallet, error) {
	s.Logger.Info("Crediting wallet", "user_id", userID, "amount", amount)
	params := db.CreditWalletParams{
		Balance: amount,
		UserID:  userID,
	}
	wallet, err := s.Store.CreditWallet(ctx, params)
	if err != nil {
		s.Logger.Error("Failed to credit wallet", "user_id", userID, "error", err)
	}
	return wallet, err
}

func (s *WalletService) Debit(ctx context.Context, userID int32, amount int64) (db.Wallet, error) {
	s.Logger.Info("Debiting wallet", "user_id", userID, "amount", amount)
	params := db.DebitWalletParams{
		Balance: amount,
		UserID:  userID,
	}
	wallet, err := s.Store.DebitWallet(ctx, params)
	if err != nil {
		s.Logger.Error("Failed to debit wallet", "user_id", userID, "error", err)
	}
	return wallet, err
}

func (s *WalletService) Transfer(ctx context.Context, fromUserID, toUserID int32, amount int64) error {
	s.Logger.Info("Starting atomic transfer", "from_user", fromUserID, "to_user", toUserID, "amount", amount)

	debitParams := db.DebitWalletParams{
		Balance: amount,
		UserID:  fromUserID,
	}
	_, err := s.Store.DebitWallet(ctx, debitParams)
	if err != nil {
		s.Logger.Error("Transfer failed: Debit failed for sender", "sender_id", fromUserID, "error", err)
		return errors.Join(err, errors.New("debit failed"))
	}
	s.Logger.Debug("Debit successful", "user_id", fromUserID)

	creditParams := db.CreditWalletParams{
		Balance: amount,
		UserID:  toUserID,
	}
	_, err = s.Store.CreditWallet(ctx, creditParams)
	if err != nil {
		s.Logger.Error("Transfer failed: Credit failed for recipient", "recipient_id", toUserID, "error", err)
		return errors.Join(err, errors.New("credit failed"))
	}
	s.Logger.Debug("Credit successful", "user_id", toUserID)

	s.Logger.Info("Atomic transfer successfully completed", "from_user", fromUserID, "to_user", toUserID, "amount", amount)
	return nil
}

func (s *WalletService) TransactionalTransfer(ctx context.Context, tx data.WalletStore, fromUserID, toUserID int32, amount int64) error {
	s.Logger.Info("Starting transactional transfer (via provided tx store)", "from_user", fromUserID, "to_user", toUserID, "amount", amount)

	debitParams := db.DebitWalletParams{Balance: amount, UserID: fromUserID}
	if _, err := tx.DebitWallet(ctx, debitParams); err != nil {
		s.Logger.Error("Transactional Debit failed", "sender_id", fromUserID, "error", err)
		return err
	}

	creditParams := db.CreditWalletParams{Balance: amount, UserID: toUserID}
	if _, err := tx.CreditWallet(ctx, creditParams); err != nil {
		s.Logger.Error("Transactional Credit failed", "recipient_id", toUserID, "error", err)
		return err
	}

	s.Logger.Info("Transactional transfer operations successful", "from_user", fromUserID)
	return nil
}
