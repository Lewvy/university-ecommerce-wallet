package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/data"
	"ecommerce/internal/service"
	"errors"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	Svc    *service.CartService
	Logger *slog.Logger
}

func CartRoutes(
	app *rest.RestHandler,
	cartSvc *service.CartService,
	logger *slog.Logger,
	protected fiber.Router,
) {
	h := &CartHandler{
		Svc:    cartSvc,
		Logger: logger,
	}

	cartGroup := protected.Group("/cart")

	cartGroup.Get("/", h.GetCartHandler)
	cartGroup.Post("/add", h.AddToCartHandler)
	cartGroup.Put("/update", h.UpdateCartItemHandler)
	cartGroup.Delete("/item/:product_id", h.DeleteCartItemHandler)
	cartGroup.Delete("/clear", h.ClearCartHandler)
}

func (h *CartHandler) AddToCartHandler(c *fiber.Ctx) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req struct {
		ProductID int64 `json:"product_id"`
		Quantity  int   `json:"quantity"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.ProductID <= 0 || req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid product_id or quantity"})
	}

	err = h.Svc.AddToCart(c.Context(), int64(userID), req.ProductID, req.Quantity)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
		}
		if err.Error() == "insufficient stock" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "insufficient stock"})
		}
		if errors.Is(err, service.ErrCannotBuyOwnItem) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you cannot add your own item to the cart"})
		}
		h.Logger.Error("Failed to add item to cart", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not add item to cart"})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": "item added to cart"})
}

func (h *CartHandler) GetCartHandler(c *fiber.Ctx) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	cart, err := h.Svc.GetCart(c.Context(), int64(userID))
	if err != nil {
		h.Logger.Error("Failed to get cart", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve cart"})
	}

	return c.Status(fiber.StatusOK).JSON(cart)
}

func (h *CartHandler) UpdateCartItemHandler(c *fiber.Ctx) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req struct {
		ProductID int64 `json:"product_id"`
		Quantity  int   `json:"quantity"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.ProductID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid product_id"})
	}

	err = h.Svc.UpdateCartItemQuantity(c.Context(), int64(userID), req.ProductID, req.Quantity)
	if err != nil {
		h.Logger.Error("Failed to update cart quantity", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update item quantity"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "cart item updated"})
}

func (h *CartHandler) DeleteCartItemHandler(c *fiber.Ctx) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	productID, err := strconv.ParseInt(c.Params("product_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid product ID"})
	}

	err = h.Svc.DeleteCartItem(c.Context(), int64(userID), productID)
	if err != nil {
		h.Logger.Error("Failed to delete cart item", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete item from cart"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "item removed from cart"})
}

func (h *CartHandler) ClearCartHandler(c *fiber.Ctx) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	err = h.Svc.ClearCart(c.Context(), int64(userID))
	if err != nil {
		h.Logger.Error("Failed to clear cart", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not clear cart"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "cart cleared"})
}
