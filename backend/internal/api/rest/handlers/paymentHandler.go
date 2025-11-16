package handlers

import (
	"context"
	"ecommerce/internal/service"

	"github.com/gofiber/fiber/v2"
)

type WalletPaymentHandler struct {
	Svc *service.WalletPaymentService
}

// POST /wallet/create-topup-order
func (h *WalletPaymentHandler) CreateTopupOrder(c *fiber.Ctx) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	var body struct {
		Amount int64 `json:"amount"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid payload"})
	}

	orderID, err := h.Svc.CreateTopupOrder(c.Context(), int64(userID), body.Amount)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create order"})
	}

	return c.JSON(fiber.Map{
		"order_id": orderID,
		"key_id":   h.Svc.RazorKeyID,
		"amount":   body.Amount,
		"currency": "INR",
	})
}

// POST /wallet/webhook
func (h *WalletPaymentHandler) RazorpayWebhook(c *fiber.Ctx) error {
	sig := c.Get("X-Razorpay-Signature")
	if sig == "" {
		return c.Status(400).SendString("signature missing")
	}

	payload := c.Body()

	if !h.Svc.VerifySignature(payload, sig) {
		return c.Status(400).SendString("invalid signature")
	}

	if err := h.Svc.HandleWebhook(context.Background(), payload); err != nil {
		return c.Status(500).SendString("webhook error")
	}

	return c.SendString("ok")
}
