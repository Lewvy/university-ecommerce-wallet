package api

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/api/rest/handlers"
	"ecommerce/internal/config"
	"ecommerce/internal/service"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func SetupServer(config *config.Config, logger *slog.Logger, userService *service.UserService) {
	app := fiber.New()
	rh := &rest.RestHandler{
		App:    app,
		Logger: logger,
	}

	handlers.UserRoutes(rh, userService)
	rh.Logger.Info("Starting server", "server", "server")
	err := app.Listen(config.Port)
	rh.Logger.Error("error running server", "err", err)
}
