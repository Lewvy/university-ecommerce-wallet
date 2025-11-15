package service

import (
	"context"
	"ecommerce/internal/data"
	db_gen "ecommerce/internal/data/gen"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/razorpay/razorpay-go"
	utils "github.com/razorpay/razorpay-go/utils"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrPaymentFailed     = errors.New("payment verification failed")
	ErrOrderMismatch     = errors.New("order/amount mismatch during verification")
)

type WalletService struct {
	BaseStore    data.WalletStore
	Pool         *pgxpool.Pool
	Logger       *slog.Logger
	RzpClient    *razorpay.Client
	RzpKeySecret string
}

type Wallet struct {
	UserID         int32 `json:"user_id"`
	Balance        int64 `json:"balance"`
	LifetimeSpent  int64 `json:"lifetime_spent"`
	LifetimeEarned int64 `json:"lifetime_earned"`
}

func NewWalletService(
	store data.WalletStore,
	pool *pgxpool.Pool,
	logger *slog.Logger,
	rzpKeyID string,
	rzpKeySecret string,
) *WalletService {

	rzpClient := razorpay.NewClient(rzpKeyID, rzpKeySecret)

	return &WalletService{
		BaseStore:    store,
		Pool:         pool,
		Logger:       logger,
		RzpClient:    rzpClient,
		RzpKeySecret: rzpKeySecret,
	}
}

func (s *WalletService) GetWalletByUserID(ctx context.Context, userID int64) (w Wallet, err error) {
	s.Logger.Debug("Getting wallet by user ID", "user_id", userID)
	dbWallet, err := s.BaseStore.GetWalletByUserID(ctx, userID)
	if err != nil {
		return w, err
	}
	w = Wallet{
		UserID:         dbWallet.UserID,
		Balance:        dbWallet.Balance,
		LifetimeSpent:  dbWallet.LifetimeSpent,
		LifetimeEarned: dbWallet.LifetimeEarned,
	}
	return w, nil
}

func (s *WalletService) CreateWallet(ctx context.Context, userID int32) (w Wallet, err error) {
	s.Logger.Info("Creating new wallet for user", "user_id", userID)
	wallet, err := s.BaseStore.CreateWallet(ctx, userID)
	if err != nil {
		s.Logger.Error("Failed to create wallet", "user_id", userID, "error", err)
		return w, err
	}
	s.Logger.Info("Successfully created wallet", "wallet_user_id", wallet.UserID)
	w = Wallet{
		UserID:         wallet.UserID,
		Balance:        wallet.Balance,
		LifetimeSpent:  wallet.LifetimeSpent,
		LifetimeEarned: wallet.LifetimeEarned,
	}
	return w, nil
}

func (s *WalletService) Credit(ctx context.Context, userID int32, amount int64) (Wallet, error) {
	s.Logger.Info("Attempting to credit wallet (ADMIN/TEST)", "user_id", userID, "amount", amount)

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return Wallet{}, err
	}
	defer tx.Rollback(ctx)

	txQueries := db_gen.New(tx)
	txStore := data.NewWalletStore(txQueries)

	wallet, err := s.creditWalletInternal(ctx, txStore, userID, amount, "credit", "completed", nil)
	if err != nil {
		return Wallet{}, err
	}
	w := Wallet{
		UserID:         wallet.UserID,
		Balance:        wallet.Balance,
		LifetimeSpent:  wallet.LifetimeSpent,
		LifetimeEarned: wallet.LifetimeEarned,
	}

	return w, tx.Commit(ctx)
}

func (s *WalletService) Debit(ctx context.Context, userID int32, amount int64) (Wallet, error) {
	s.Logger.Info("Attempting to debit wallet", "user_id", userID, "amount", amount)

	if amount <= 0 {
		return Wallet{}, errors.New("debit amount must be positive")
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return Wallet{}, err
	}
	defer tx.Rollback(ctx)

	txQueries := db_gen.New(tx)
	txStore := data.NewWalletStore(txQueries)

	wallet, err := txStore.GetWalletByUserIDForUpdate(ctx, userID)
	if err != nil {
		return Wallet{}, err
	}

	if wallet.Balance < amount {
		return Wallet{}, ErrInsufficientFunds
	}

	updatedWallet, err := s.creditWalletInternal(ctx, txStore, userID, -amount, "debit", "completed", nil)
	if err != nil {
		return Wallet{}, err
	}
	w := Wallet{
		UserID:         updatedWallet.UserID,
		Balance:        updatedWallet.Balance,
		LifetimeSpent:  updatedWallet.LifetimeSpent,
		LifetimeEarned: updatedWallet.LifetimeEarned,
	}

	return w, tx.Commit(ctx)
}

func (s *WalletService) Transfer(ctx context.Context, txStore data.WalletStore, senderID int32, recipientID int32, amount int64) error {
	s.Logger.Info("Attempting atomic transfer", "sender_id", senderID, "recipient_id", recipientID, "amount", amount)

	senderWallet, err := txStore.GetWalletByUserIDForUpdate(ctx, senderID)
	if err != nil {
		return err
	}

	if senderWallet.Balance < amount {
		s.Logger.Warn("Insufficient funds for transfer", "sender_id", senderID, "balance", senderWallet.Balance, "requested", amount)
		return ErrInsufficientFunds
	}
	s.Logger.Debug("Sender funds sufficient", "sender_id", senderID)

	_, err = s.creditWalletInternal(ctx, txStore, senderID, -amount, "transfer_out", "completed", &recipientID)
	if err != nil {
		s.Logger.Error("Failed to debit sender", "sender_id", senderID, "error", err)
		return err
	}
	s.Logger.Debug("Successfully debited sender", "sender_id", senderID, "amount", amount)

	_, err = s.creditWalletInternal(ctx, txStore, recipientID, amount, "transfer_in", "completed", &senderID)
	if err != nil {
		s.Logger.Error("Failed to credit recipient", "recipient_id", recipientID, "error", err)
		return err
	}
	s.Logger.Debug("Successfully credited recipient", "recipient_id", recipientID, "amount", amount)

	return nil
}

func (s *WalletService) CreatePaymentOrder(ctx context.Context, userID int32, amount int64) (map[string]interface{}, error) {
	s.Logger.Info("Creating Razorpay order", "user_id", userID, "amount", amount)

	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	amountInPaise := amount * 100

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	txQueries := db_gen.New(tx)
	txStore := data.NewWalletStore(txQueries)

	txParams := db_gen.CreateTransactionParams{
		UserID:            userID,
		Amount:            amount,
		TransactionType:   "credit_pending",
		TransactionStatus: "pending",
	}

	dbTx, err := txStore.CreateTransaction(ctx, txParams)
	if err != nil {
		s.Logger.Error("Failed to create pending transaction", "user_id", userID, "error", err)
		return nil, err
	}

	orderData := map[string]interface{}{
		"amount":          amountInPaise,
		"currency":        "INR",
		"receipt":         dbTx.ID,
		"payment_capture": 1,
	}

	rzpOrder, err := s.RzpClient.Order.Create(orderData, nil)
	if err != nil {
		s.Logger.Error("Failed to create Razorpay order", "user_id", userID, "error", err)
		return nil, err
	}

	orderID := rzpOrder["id"].(string)

	orderIdtext := data.NewPGText(orderID)
	_, err = txStore.UpdateTransactionOrderID(ctx, db_gen.UpdateTransactionOrderIDParams{
		ID:              dbTx.ID,
		RazorpayOrderID: orderIdtext,
	})
	if err != nil {
		s.Logger.Error("Failed to link transaction to order_id", "tx_id", dbTx.ID, "error", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	s.Logger.Info("Razorpay order created", "user_id", userID, "order_id", orderID)
	response := map[string]interface{}{
		"order_id": orderID,
		"amount":   amountInPaise,
		"currency": "INR",
	}
	return response, nil
}

func (s *WalletService) VerifyPayment(ctx context.Context, rzpOrderID, rzpPaymentID, rzpSignature string) (db_gen.Wallet, error) {
	s.Logger.Info("Verifying Razorpay payment", "order_id", rzpOrderID, "payment_id", rzpPaymentID)

	params := map[string]any{
		"razorpay_order_id":   rzpOrderID,
		"razorpay_payment_id": rzpPaymentID,
	}
	ok := utils.VerifyPaymentSignature(params, rzpSignature, s.RzpKeySecret)
	if !ok {
		s.Logger.Warn("Razorpay signature verification failed", "order_id", rzpOrderID)
		return db_gen.Wallet{}, ErrPaymentFailed
	}
	s.Logger.Debug("Razorpay signature verified", "order_id", rzpOrderID)

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return db_gen.Wallet{}, err
	}
	defer tx.Rollback(ctx)

	txQueries := db_gen.New(tx)
	txStore := data.NewWalletStore(txQueries)

	dbTx, err := txStore.GetTransactionByOrderID(ctx, rzpOrderID)
	if err != nil {
		s.Logger.Error("Failed to find transaction by order_id", "order_id", rzpOrderID, "error", err)
		return db_gen.Wallet{}, err
	}

	if dbTx.TransactionStatus != "pending" {
		s.Logger.Warn("Transaction already processed", "order_id", rzpOrderID, "status", dbTx.TransactionStatus)
		return s.BaseStore.GetWalletByUserID(ctx, int64(dbTx.UserID))
	}

	rgid := data.NewPGText(rzpPaymentID)
	_, err = txStore.UpdateTransactionStatus(ctx, db_gen.UpdateTransactionStatusParams{
		ID:                dbTx.ID,
		TransactionStatus: "completed",
		RazorpayPaymentID: rgid,
	})
	if err != nil {
		s.Logger.Error("Failed to update transaction status", "tx_id", dbTx.ID, "error", err)
		return db_gen.Wallet{}, err
	}

	wallet, err := s.creditWalletInternal(ctx, txStore, dbTx.UserID, dbTx.Amount, "credit_payment", "completed", nil)
	if err != nil {
		s.Logger.Error("Failed to credit wallet after verification", "user_id", dbTx.UserID, "error", err)
		return db_gen.Wallet{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return db_gen.Wallet{}, err
	}

	s.Logger.Info("Payment verified and wallet credited", "user_id", wallet.UserID)
	return wallet, nil
}

func (s *WalletService) creditWalletInternal(
	ctx context.Context,
	txStore data.WalletStore,
	userID int32,
	amount int64,
	txType string,
	txStatus string,
	relatedUserID *int32,
) (db_gen.Wallet, error) {

	txParams := db_gen.CreateTransactionParams{
		UserID:            userID,
		Amount:            amount,
		TransactionType:   txType,
		TransactionStatus: txStatus,
	}

	if relatedUserID != nil {
		txParams.RelatedUserID = pgtype.Int4{Int32: *relatedUserID, Valid: true}
	}
	if _, err := txStore.CreateTransaction(ctx, txParams); err != nil {
		s.Logger.Error("Failed to create transaction record", "user_id", userID, "error", err)
		return db_gen.Wallet{}, err
	}

	var wallet db_gen.Wallet
	var err error

	if amount > 0 {
		params := db_gen.CreditWalletParams{
			Balance: amount,
			UserID:  userID,
		}
		wallet, err = txStore.CreditWallet(ctx, params)
	} else {
		params := db_gen.DebitWalletParams{
			Balance: -amount,
			UserID:  userID,
		}
		wallet, err = txStore.DebitWallet(ctx, params)
	}

	if err != nil {
		s.Logger.Error("Failed to update wallet balance", "user_id", userID, "error", err)
		return db_gen.Wallet{}, err
	}

	return wallet, nil
}
