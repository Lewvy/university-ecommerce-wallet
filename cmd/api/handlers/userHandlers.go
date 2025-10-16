package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
}

func (h *Handler) RegisterUserHandler(c *fiber.Ctx) error {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	readBody(c, &input)

	return c.Status(http.StatusCreated).JSON(&input)
}

func (h *Handler) LoginUserHandler(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	readBody(c, &input)
	return c.Status(http.StatusOK).JSON(&input)
}
