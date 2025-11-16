package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletPaymentService struct {
	Pool          *pgxpool.Pool
	Logger        *slog.Logger
	RazorKeyID    string
	RazorSecret   string
	WebhookSecret string
	Client        *http.Client
}

func NewWalletPaymentService(pool *pgxpool.Pool, logger *slog.Logger) *WalletPaymentService {
	return &WalletPaymentService{
		Pool:          pool,
		Logger:        logger,
		RazorKeyID:    os.Getenv("RAZORPAY_ID"),
		RazorSecret:   os.Getenv("RAZORPAY_SECRET"),
		WebhookSecret: os.Getenv("RAZORPAY_WEBHOOK_SECRET"),
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Create Topup transaction + create Razorpay order
func (s *WalletPaymentService) CreateTopupOrder(ctx context.Context, userID int64, amount int64) (string, error) {
	// Insert db row
	var txnID int64
	err := s.Pool.QueryRow(ctx, `
		INSERT INTO wallet_transactions (user_id, amount, transaction_type, transaction_status)
		VALUES ($1, $2, 'razorpay_topup', 'pending')
		RETURNING id
	`, userID, amount).Scan(&txnID)
	if err != nil {
		s.Logger.Error("Insert txn failed", "error", err)
		return "", err
	}

	receipt := fmt.Sprintf("wal_txn_%d", txnID)

	// Prepare Razorpay order
	body, _ := json.Marshal(map[string]any{
		"amount":          amount,
		"currency":        "INR",
		"receipt":         receipt,
		"payment_capture": 1,
	})

	req, _ := http.NewRequest("POST", "https://api.razorpay.com/v1/orders", bytes.NewReader(body))
	req.SetBasicAuth(s.RazorKeyID, s.RazorSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		s.Logger.Error("Razorpay request failed", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		s.Logger.Error("Razorpay error", "status", resp.StatusCode, "body", string(respBody))
		return "", fmt.Errorf("razorpay error %d", resp.StatusCode)
	}

	var r struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(respBody, &r)

	// Save Razorpay order ID
	_, err = s.Pool.Exec(ctx, `
		UPDATE wallet_transactions
		SET razorpay_order_id = $1
		WHERE id = $2
	`, r.ID, txnID)
	if err != nil {
		s.Logger.Error("Failed to update razorpay_order_id", "error", err)
	}

	return r.ID, nil
}

// Verify webhook signature
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

// Handle Razorpay webhook
func (s *WalletPaymentService) HandleWebhook(ctx context.Context, payload []byte) error {
	var data struct {
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

	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}

	payment := data.Payload.Payment.Entity
	if payment.OrderID == "" {
		return errors.New("missing order_id")
	}

	if payment.Status == "captured" {
		// Update DB
		_, err := s.Pool.Exec(ctx, `
			UPDATE wallet_transactions
			SET transaction_status = 'success',
			    razorpay_payment_id = $1
			WHERE razorpay_order_id = $2
		`, payment.ID, payment.OrderID)
		if err != nil {
			s.Logger.Error("Failed DB update on webhook", "error", err)
			return err
		}

		// TODO: credit wallet here
		s.Logger.Info("Wallet topup successful", "order_id", payment.OrderID)
	}

	return nil
}
