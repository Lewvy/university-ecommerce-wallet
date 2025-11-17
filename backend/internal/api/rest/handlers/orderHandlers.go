package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/service"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	Svc    *service.OrderService
	Logger *slog.Logger
}

func OrderRoutes(
	app *rest.RestHandler,
	orderSvc *service.OrderService,
	logger *slog.Logger,
	protected fiber.Router,
) {
	h := &OrderHandler{
		Svc:    orderSvc,
		Logger: logger,
	}

	orderGroup := protected.Group("/orders")

	orderGroup.Post("/", h.CreateOrderFromCartHandler)
	// TODO: Add GET routes for order history
	// orderGroup.Get("/", h.GetMyOrdersHandler)
	// orderGroup.Get("/:id", h.GetOrderDetailsHandler)
}

func (h *OrderHandler) CreateOrderFromCartHandler(c *fiber.Ctx) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	order, err := h.Svc.CreateOrderFromCart(c.Context(), int64(userID))
	if err != nil {
		h.Logger.Error("Failed to create order from cart", "error", err)
		if errors.Is(err, service.ErrCartEmpty) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cart is empty"})
		}
		if errors.Is(err, service.ErrInsufficientFunds) {
			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{"error": "insufficient funds"})
		}
		if err.Error() == "insufficient stock" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create order"})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}
