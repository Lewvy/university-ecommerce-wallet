package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/data"
	db_gen "ecommerce/internal/data/gen"
	"ecommerce/internal/service"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type walletTransferRequest struct {
	RecipientUserID int32 `json:"recipient_user_id" validate:"required"`
	Amount          int64 `json:"amount" validate:"required,min=1"`
}

type walletAmountRequest struct {
	Amount int64 `json:"amount" validate:"required,min=1"`
}

type WalletHandler struct {
	Svc  *service.WalletService
	Pool *pgxpool.Pool
}

func WalletRoutes(rh *rest.RestHandler, walletService *service.WalletService, dbConn *pgxpool.Pool, protected fiber.Router) {
	h := WalletHandler{
		Svc:  walletService,
		Pool: dbConn,
	}

	protected.Get("/wallet/balance", h.GetBalanceHandler)
	protected.Post("/wallet/transfer", h.TransferHandler)
	protected.Post("/wallet/credit", h.CreditHandler)
	protected.Post("/wallet/debit", h.DebitHandler)
}

func getCurrentUserID(c *fiber.Ctx) (int32, error) {
	userID64, ok := c.Locals("authenticatedUserID").(int64)

	if !ok || userID64 == 0 {
		return 0, errors.New("unauthenticated or missing user ID in context")
	}
	return int32(userID64), nil
}

func (h *WalletHandler) GetBalanceHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	userID, err := getCurrentUserID(c)
	if err != nil {
		h.Svc.Logger.Warn("Auth error", "error", err)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	wallet, err := h.Svc.GetWalletByUserID(ctx, int64(userID))
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "wallet not found"})
		}
		h.Svc.Logger.Error("Internal error getting balance", "user_id", userID, "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(http.StatusOK).JSON(wallet)
}

func (h *WalletHandler) CreditHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var input walletAmountRequest

	userID, err := getCurrentUserID(c)
	if err != nil {
		h.Svc.Logger.Warn("Auth error", "error", err)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.Amount <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "amount must be greater than zero"})
	}

	wallet, err := h.Svc.Credit(ctx, int32(userID), input.Amount)
	if err != nil {
		h.Svc.Logger.Error("Failed to credit wallet", "user_id", userID, "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "credit operation failed"})
	}

	return c.Status(http.StatusOK).JSON(wallet)
}

func (h *WalletHandler) DebitHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var input walletAmountRequest

	userID, err := getCurrentUserID(c)
	if err != nil {
		h.Svc.Logger.Warn("Auth error", "error", err)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if input.Amount <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "amount must be greater than zero"})
	}

	wallet, err := h.Svc.Debit(ctx, int32(userID), input.Amount)
	if err != nil {
		h.Svc.Logger.Error("Failed to debit wallet", "user_id", userID, "error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "debit operation failed (e.g., insufficient funds or constraint violation)"})
	}

	return c.Status(http.StatusOK).JSON(wallet)
}

func (h *WalletHandler) TransferHandler(c *fiber.Ctx) error {
	var input walletTransferRequest
	ctx := c.Context()

	senderID, err := getCurrentUserID(c)
	if err != nil {
		h.Svc.Logger.Warn("Auth error", "error", err)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if senderID == input.RecipientUserID {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "cannot transfer to yourself"})
	}
	if input.Amount <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "amount must be greater than zero"})
	}

	tx, err := h.Pool.Begin(ctx)
	if err != nil {
		h.Svc.Logger.Error("Failed to begin transaction", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to initiate transfer"})
	}

	defer tx.Rollback(ctx)

	txQueries := db_gen.New(tx)

	txStore := data.NewWalletStore(txQueries)

	err = h.Svc.Transfer(ctx, txStore, int32(senderID), input.RecipientUserID, input.Amount)
	if err != nil {

		h.Svc.Logger.Warn("Atomic Transfer failed, transaction rolled back", "sender_id", senderID, "recipient_id", input.RecipientUserID, "error", err.Error())

		if errors.Is(err, service.ErrInsufficientFunds) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "insufficient funds"})
		}

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "transfer failed: " + err.Error()})
	}

	if err = tx.Commit(ctx); err != nil {
		h.Svc.Logger.Error("Failed to commit transaction", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "transfer committed but database error occurred"})
	}

	h.Svc.Logger.Info("Atomic Transfer successful, transaction committed", "sender_id", senderID, "recipient_id", input.RecipientUserID)
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "transfer successful"})
}
