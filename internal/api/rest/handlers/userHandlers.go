package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/dto"
	"ecommerce/internal/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	svc *service.UserService
}

func UserRoutes(rh *rest.RestHandler, userService *service.UserService) {
	app := rh.App

	h := UserHandler{
		svc: userService,
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
	h.svc.Logger.Info("message", "cart", "empty")
	return c.Status(http.StatusOK).JSON(map[string]string{
		"message": "get cart",
	})
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) RegisterUserHandler(c *fiber.Ctx) error {
	var userSignup dto.UserSignup

	err := c.BodyParser(&userSignup)
	if err != nil {
		h.svc.Logger.Error("Error decoding", "err", err)
		return err
	}

	validation_check, err := h.svc.Signup(userSignup)
	if validation_check != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": validation_check,
		})
	}
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(&userSignup, "application/json")
}

func (h *UserHandler) LoginUserHandler(c *fiber.Ctx) error {
	var input dto.UserLogin

	err := c.BodyParser(&input)
	if err != nil {
		h.svc.Logger.Error("Error decoding user", "err", err)
		return c.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid email/password",
			},
		)
	}

	err = h.svc.Login(input)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(&input)
}
