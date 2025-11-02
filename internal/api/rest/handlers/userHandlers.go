package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/dto"
	"ecommerce/internal/service"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	logger *slog.Logger
	svc    service.UserService
}

func UserRoutes(rh *rest.RestHandler) {
	app := rh.App

	svc := service.UserService{}

	h := UserHandler{
		logger: rh.Logger,
		svc:    svc,
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
	h.logger.Info("message", "cart", "empty")
	return c.Status(http.StatusOK).JSON(map[string]string{
		"message": "get cart",
	})
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) RegisterUserHandler(c *fiber.Ctx) error {
	var input dto.UserSignup
	err := c.BodyParser(&input)
	if err != nil {
		h.logger.Error("Error decoding", "err", err)
		return err
	}

	validation_check, err := h.svc.Signup(input)
	if validation_check != nil {
		c.Status(http.StatusBadRequest).JSON(validation_check)
		return nil
	}
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(err)
		return err
	}
	return c.Status(http.StatusCreated).JSON(&input, "application/text")
}

func (h *UserHandler) LoginUserHandler(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.BodyParser(&input)
	return c.Status(http.StatusOK).JSON(&input)
}
