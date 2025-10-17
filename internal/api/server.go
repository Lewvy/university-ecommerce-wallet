package api

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/api/rest/handlers"
	"ecommerce/internal/config"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func SetupServer(config *config.Config, logger *slog.Logger) {
	app := fiber.New()
	rh := &rest.RestHandler{
		App:    app,
		Logger: logger,
	}

	setupRoutes(rh)
	app.Listen(config.Port)
}

func setupRoutes(rh *rest.RestHandler) {
	handlers.UserRoutes(rh)
}
