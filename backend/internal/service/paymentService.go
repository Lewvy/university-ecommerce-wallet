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

	"github.com/jackc/pgx/v5/pgtype"
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
		RazorKeyID:    os.Getenv("RAZORPAY_ID"),
		RazorSecret:   os.Getenv("RAZORPAY_SECRET"),
		WebhookSecret: os.Getenv("RAZORPAY_WEBHOOK_SECRET"),
		Client:        &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *WalletPaymentService) CreateTopupOrder(ctx context.Context, userID int64, amount int64) (string, error) {
	s.Logger.Info("Creating Razorpay top-up order", "user_id", userID, "amount", amount)

	txRow, err := s.Store.CreateTransaction(ctx, db.CreateTransactionParams{
		UserID:            int32(userID),
		Amount:            amount,
		TransactionType:   "razorpay_topup",
		TransactionStatus: "pending",

		RelatedUserID:     pgtype.Int4{Int32: 0, Valid: false},
		RazorpayOrderID:   pgtype.Text{String: "", Valid: false},
		RazorpayPaymentID: pgtype.Text{String: "", Valid: false},
	})

	if err != nil {
		s.Logger.Error("Failed to create local transaction record", "error", err)
		return "", fmt.Errorf("db transaction creation failed: %w", err)
	}

	receipt := fmt.Sprintf("wallet_txn_%d", txRow.ID)

	bodyMap := map[string]any{
		"amount":          amount,
		"currency":        "INR",
		"receipt":         receipt,
		"payment_capture": 1,
	}

	bodyBytes, _ := json.Marshal(bodyMap)

	req, err := http.NewRequest("POST",
		"https://api.razorpay.com/v1/orders",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		s.Logger.Error("Failed to build Razorpay request", "error", err)
		return "", err
	}

	req.SetBasicAuth(s.RazorKeyID, s.RazorSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		s.Logger.Error("Failed to call Razorpay Orders API", "error", err)
		return "", fmt.Errorf("razorpay request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		s.Logger.Error("Razorpay returned an error", "status", resp.Status, "response", string(respBody))
		return "", fmt.Errorf("razorpay error %s", resp.Status)
	}

	var r struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBody, &r); err != nil {
		s.Logger.Error("Failed to parse Razorpay order response", "error", err, "body", string(respBody))
		return "", err
	}

	_, err = s.Store.UpdateTransactionOrderID(ctx, db.UpdateTransactionOrderIDParams{
		ID: txRow.ID,
		RazorpayOrderID: pgtype.Text{
			String: r.ID,
			Valid:  true,
		},
	})

	if err != nil {
		s.Logger.Error("Failed to update Razorpay order ID in DB", "error", err)
		return "", err
	}

	s.Logger.Info("Razorpay order created successfully",
		"user_id", userID,
		"transaction_id", txRow.ID,
		"order_id", r.ID)

	return r.ID, nil
}

func (s *WalletPaymentService) VerifySignature(body []byte, provided string) bool {
	secret := os.Getenv("RAZORPAY_WEBHOOK_SECRET")
	if secret == "" {
		s.Logger.Error("Missing webhook secret")
		return false
	}

	s.Logger.Info("Backend received body",
		"length", len(body),
		"hex", hex.EncodeToString(body),
	)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedBytes := mac.Sum(nil)
	expectedHex := hex.EncodeToString(expectedBytes)

	s.Logger.Info("Signature Check",
		"expected", expectedHex,
		"provided", provided,
	)

	return hmac.Equal([]byte(expectedHex), []byte(provided))
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

	txRow, err := s.Store.GetTransactionByOrderID(ctx, p.OrderID)
	if err != nil {
		return err
	}

	_, err = s.Store.UpdateTransactionStatus(ctx, db.UpdateTransactionStatusParams{
		ID:                txRow.ID,
		TransactionStatus: "success",
		RazorpayPaymentID: data.NewPGText(p.ID),
	})
	if err != nil {
		return err
	}

	err = s.Store.CreditWalletBalance(ctx, txRow.UserID, txRow.Amount)
	if err != nil {
		return err
	}

	s.Logger.Info("Wallet topup successful", "order_id", p.OrderID)
	return nil
}
