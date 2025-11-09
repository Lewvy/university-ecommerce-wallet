package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/dto"
	"ecommerce/internal/service"
	"ecommerce/internal/validator"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Svc *service.UserService
}

func UserRoutes(rh *rest.RestHandler, userService *service.UserService, protected fiber.Router) {
	app := rh.App

	h := UserHandler{
		Svc: userService,
	}
	app.Post("/register", h.RegisterUserHandler)

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
	input := &service.UserVerification{}
	err := c.BodyParser(&input)
	if err != nil {
		h.Svc.Logger.Error("Error decoding request body", "error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	err = h.Svc.VerifyUser(c.Context(), input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid token",
		})
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "user verified",
	})
}

func (h *UserHandler) CreateProfile(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) CreateCart(c *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) GetCart(c *fiber.Ctx) error {
	h.Svc.Logger.Info("message", "cart", "empty")
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
		h.Svc.Logger.Error("Error decoding", "err", err)
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Code:    "bad_request",
			Message: "Invalid request body",
		})
	}

	user, err := h.Svc.Signup(c.Context(), userSignup)
	if err != nil {
		var validationError *validator.ValidationError
		if errors.As(err, &validationError) {
			h.Svc.Logger.Warn("User Validation Error", "error", validationError.Errors)
			return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
				Code:    "validation_error",
				Message: "invalid details",
				Fields:  validationError.Errors})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "An internal server error has occured",
		})
	}

	return c.Status(http.StatusCreated).JSON(&user, "application/json")
}
