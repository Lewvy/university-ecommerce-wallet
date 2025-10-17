package handlers

import (
	"ecommerce/internal/api/helpers"
	"ecommerce/internal/api/rest"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
}

func UserRoutes(rh *rest.RestHandler) {
	app := rh.App

	//TODO:: create a user service and pass it to handler

	h := UserHandler{}
	app.Post("/v1/users/register", h.RegisterUserHandler)
	app.Post("/v1/users/login", h.LoginUserHandler)
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
