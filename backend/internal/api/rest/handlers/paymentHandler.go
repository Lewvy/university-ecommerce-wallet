package handlers

import (
	"context"
	"ecommerce/internal/service"

	"github.com/gofiber/fiber/v2"
)

type WalletPaymentHandler struct {
	Svc *service.WalletPaymentService
}

func NewWalletPaymentHandler(svc *service.WalletPaymentService) *WalletPaymentHandler {
	return &WalletPaymentHandler{Svc: svc}
}

func (h *WalletPaymentHandler) CreateTopupOrder(c *fiber.Ctx) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req struct{ Amount int64 }
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	orderID, err := h.Svc.CreateTopupOrder(c.Context(), int64(userID), req.Amount)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create order"})
	}

	return c.JSON(fiber.Map{
		"order_id": orderID,
		"key_id":   h.Svc.RazorKeyID,
		"amount":   req.Amount,
		"currency": "INR",
	})
}

func (h *WalletPaymentHandler) RazorpayWebhook(c *fiber.Ctx) error {
	sig := c.Get("X-Razorpay-Signature")
	body := c.Body()

	if !h.Svc.VerifySignature(body, sig) {
		return c.Status(400).SendString("invalid signature")
	}

	if err := h.Svc.HandleWebhook(context.Background(), body); err != nil {
		return c.Status(500).SendString("webhook processing error")
	}

	return c.SendString("ok")
}
