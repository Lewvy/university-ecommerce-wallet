package api

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/api/rest/handlers"
	"ecommerce/internal/config"

	"github.com/gofiber/fiber/v2"
)

func SetupServer(config *config.Config) {
	app := fiber.New()
	rh := &rest.RestHandler{
		App: app,
	}

	setupRoutes(rh)
	app.Listen(config.Port)
}

func setupRoutes(rh *rest.RestHandler) {
	handlers.UserRoutes(rh)
}
