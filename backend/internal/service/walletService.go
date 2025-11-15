package service

import (
	"context"
	"ecommerce/internal/data"
	db_gen "ecommerce/internal/data/gen"
	"errors"
	"log/slog"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type WalletService struct {
	BaseStore data.WalletStore
	Logger    *slog.Logger
}

func NewWalletService(store data.WalletStore, logger *slog.Logger) *WalletService {
	return &WalletService{BaseStore: store, Logger: logger}
}

func (s *WalletService) GetWalletByUserID(ctx context.Context, userID int32) (db_gen.Wallet, error) {
	s.Logger.Debug("Getting wallet by user ID", "user_id", userID)
	return s.BaseStore.GetWalletByUserID(ctx, userID)
}

func (s *WalletService) CreateWallet(ctx context.Context, userID int32) (db_gen.Wallet, error) {
	s.Logger.Info("Creating new wallet for user", "user_id", userID)
	wallet, err := s.BaseStore.CreateWallet(ctx, userID)
	if err != nil {
		s.Logger.Error("Failed to create wallet", "user_id", userID, "error", err)
		return db_gen.Wallet{}, err
	}
	s.Logger.Info("Successfully created wallet", "wallet_user_id", wallet.UserID)
	return wallet, nil
}

func (s *WalletService) Credit(ctx context.Context, userID int32, amount int64) (db_gen.Wallet, error) {
	s.Logger.Info("Crediting wallet", "user_id", userID, "amount", amount)
	params := db_gen.CreditWalletParams{
		Balance: amount,
		UserID:  userID,
	}
	return s.BaseStore.CreditWallet(ctx, params)
}

func (s *WalletService) Debit(ctx context.Context, userID int32, amount int64) (db_gen.Wallet, error) {
	s.Logger.Info("Debiting wallet", "user_id", userID, "amount", amount)
	params := db_gen.DebitWalletParams{
		Balance: amount,
		UserID:  userID,
	}
	return s.BaseStore.DebitWallet(ctx, params)
}

func (s *WalletService) Transfer(ctx context.Context, txStore data.WalletStore, senderID int32, recipientID int32, amount int64) error {
	s.Logger.Info("Attempting atomic transfer", "sender_id", senderID, "recipient_id", recipientID, "amount", amount)

	senderWallet, err := txStore.GetWalletByUserIDForUpdate(ctx, senderID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			s.Logger.Warn("Sender wallet not found for transfer", "sender_id", senderID)
			return errors.New("sender wallet not found")
		}
		s.Logger.Error("Failed to get sender wallet for update", "sender_id", senderID, "error", err)
		return err
	}

	if senderWallet.Balance < amount {
		s.Logger.Warn("Insufficient funds for transfer", "sender_id", senderID, "balance", senderWallet.Balance, "requested", amount)
		return ErrInsufficientFunds
	}
	s.Logger.Debug("Sender funds sufficient", "sender_id", senderID, "balance", senderWallet.Balance)

	debitParams := db_gen.DebitWalletParams{
		Balance: amount,
		UserID:  senderID,
	}
	if _, err = txStore.DebitWallet(ctx, debitParams); err != nil {
		s.Logger.Error("Failed to debit sender wallet", "sender_id", senderID, "error", err)
		return err
	}
	s.Logger.Debug("Successfully debited sender", "sender_id", senderID, "amount", amount)

	creditParams := db_gen.CreditWalletParams{
		Balance: amount,
		UserID:  recipientID,
	}
	if _, err = txStore.CreditWallet(ctx, creditParams); err != nil {
		s.Logger.Error("Failed to credit recipient wallet", "recipient_id", recipientID, "error", err)

		return err
	}
	s.Logger.Debug("Successfully credited recipient", "recipient_id", recipientID, "amount", amount)

	return nil
}
