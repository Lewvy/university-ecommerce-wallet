package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"ecommerce/internal/data"
	db "ecommerce/internal/data/gen"
	"log/slog"
)

type WalletPaymentService struct {
	Store         data.WalletStore
	Logger        *slog.Logger
	RazorKeyID    string
	RazorSecret   string
	WebhookSecret string
	Client        *http.Client
}

func NewWalletPaymentService(store data.WalletStore, logger *slog.Logger) *WalletPaymentService {
	return &WalletPaymentService{
		Store:         store,
		Logger:        logger,
		RazorKeyID:    os.Getenv("RAZORPAY_KEY_ID"),
		RazorSecret:   os.Getenv("RAZORPAY_KEY_SECRET"),
		WebhookSecret: os.Getenv("RAZORPAY_WEBHOOK_SECRET"),
		Client:        &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *WalletPaymentService) CreateTopupOrder(ctx context.Context, userID int64, amount int64) (string, error) {
	txRow, err := s.Store.CreateTransaction(ctx, db.CreateTransactionParams{
		UserID:            int32(userID),
		Amount:            amount,
		TransactionType:   "razorpay_topup",
		TransactionStatus: "pending",
	})
	if err != nil {
		return "", err
	}

	receipt := fmt.Sprintf("wallet_txn_%d", txRow.ID)

	// 2) Create Razorpay order
	body, _ := json.Marshal(map[string]any{
		"amount":          amount,
		"currency":        "INR",
		"receipt":         receipt,
		"payment_capture": 1,
	})

	req, _ := http.NewRequest("POST",
		"https://api.razorpay.com/v1/orders",
		bytes.NewReader(body),
	)

	req.SetBasicAuth(s.RazorKeyID, s.RazorSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("razorpay error: %s %s", resp.Status, respBody)
	}

	var r struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(respBody, &r)

	// 3) Save Razorpay order ID
	_, err = s.Store.UpdateTransactionOrderID(ctx, db.UpdateTransactionOrderIDParams{
		ID:              txRow.ID,
		RazorpayOrderID: data.NewPGText(r.ID),
	})
	if err != nil {
		return "", err
	}

	return r.ID, nil
}

func (s *WalletPaymentService) VerifySignature(body []byte, sig string) bool {
	secret := s.WebhookSecret
	if secret == "" {
		secret = s.RazorSecret
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	expected := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(sig))
}

func (s *WalletPaymentService) HandleWebhook(ctx context.Context, payload []byte) error {
	var ev struct {
		Event   string `json:"event"`
		Payload struct {
			Payment struct {
				Entity struct {
					ID      string `json:"id"`
					OrderID string `json:"order_id"`
					Status  string `json:"status"`
				} `json:"entity"`
			} `json:"payment"`
		} `json:"payload"`
	}

	if err := json.Unmarshal(payload, &ev); err != nil {
		return err
	}

	p := ev.Payload.Payment.Entity

	if p.Status != "captured" {
		return nil
	}

	// 1) Find the transaction
	txRow, err := s.Store.GetTransactionByOrderID(ctx, p.OrderID)
	if err != nil {
		return err
	}

	// 2) Mark success
	_, err = s.Store.UpdateTransactionStatus(ctx, db.UpdateTransactionStatusParams{
		ID:                txRow.ID,
		TransactionStatus: "success",
		RazorpayPaymentID: data.NewPGText(p.ID),
	})
	if err != nil {
		return err
	}

	// 3) Credit wallet
	err = s.Store.CreditWalletBalance(ctx, txRow.UserID, txRow.Amount)
	if err != nil {
		return err
	}

	s.Logger.Info("Wallet topup successful", "order_id", p.OrderID)
	return nil
}
