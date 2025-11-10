package api

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/api/rest/handlers"
	"ecommerce/internal/config"
	"ecommerce/internal/middleware"
	"ecommerce/internal/service"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupServer(config *config.Config, logger *slog.Logger, userService *service.UserService, tokenService *service.TokenService) {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "POST, PATCH, PUT, GET, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	rh := &rest.RestHandler{
		App:    app,
		Logger: logger,
	}

	authMiddleware := middleware.AuthMiddleware()

	userHandler := &handlers.UserHandler{Svc: userService}

	rh.App.Post("/register", userHandler.RegisterUserHandler)
	rh.App.Post("/login", userHandler.LoginUserHandler)

	protected := app.Group("/", authMiddleware)
	handlers.UserRoutes(rh, userService, protected)
	handlers.TokenRoutes(rh, tokenService)

	rh.Logger.Info("Starting server", "server", "server")
	err := app.Listen(config.Port)
	rh.Logger.Error("error running server", "err", err)
}
