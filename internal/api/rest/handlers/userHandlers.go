package handlers

import (
	"ecommerce/internal/api/helpers"
	"ecommerce/internal/api/rest"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	l *slog.Logger
}

func UserRoutes(rh *rest.RestHandler) {
	app := rh.App

	//TODO:: create a user service and pass it to handler

	h := UserHandler{
		l: rh.Logger,
	}
	app.Post("/register", h.RegisterUserHandler)
	app.Post("/login", h.LoginUserHandler)

	app.Post("/verify", h.Verify)
	app.Get("/verify", h.GetVerificationCode)

	app.Post("/profile", h.CreateProfile)
	app.Get("/profile", h.GetProfile)

	app.Post("/cart", h.CreateCart)
	app.Get("/cart", h.GetCart)
	app.Post("/order", h.CreateOrder)
	app.Get("/order/:id", h.GetOrder)

	app.Post("/become-seller", h.BecomeSeller)
}

func (h *UserHandler) GetVerificationCode(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) CreateOrder(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) BecomeSeller(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) GetOrder(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) Verify(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) CreateProfile(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) CreateCart(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) GetCart(c *fiber.Ctx) error {
	h.l.Info("message", "cart", "empty")
	return c.Status(http.StatusOK).JSON(map[string]string{
		"message": "get cart",
	})
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) RegisterUserHandler(c *fiber.Ctx) error {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	helpers.ReadBody(c, &input)

	return c.Status(http.StatusCreated).JSON(&input, "application/text")
}

func (h *UserHandler) LoginUserHandler(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	helpers.ReadBody(c, &input)
	return c.Status(http.StatusOK).JSON(&input)
}
