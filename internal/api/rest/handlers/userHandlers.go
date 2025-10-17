package handlers

import (
	"ecommerce/internal/api/helpers"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	app *fiber.App
}

func (u *UserHandler) UserHandlerRoutes() {
	u.app.Post("/v1/users/register", u.RegisterUserHandler)
	u.app.Post("/v1/users/login", u.LoginUserHandler)
}

func (h *UserHandler) RegisterUserHandler(c *fiber.Ctx) error {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	helpers.ReadBody(c, &input)

	return c.Status(http.StatusCreated).JSON(&input)
}

func (h *UserHandler) LoginUserHandler(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	helpers.ReadBody(c, &input)
	return c.Status(http.StatusOK).JSON(&input)
}
