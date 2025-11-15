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
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupServer(
	config *config.Config,
	logger *slog.Logger,
	userService *service.UserService,
	tokenService *service.TokenService,
	walletService *service.WalletService,
	productService *service.ProductService,
	cloudinaryService *service.CloudinaryService,
	dbPool *pgxpool.Pool,
) {
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
	app.Post("/verify", userHandler.Verify)
	app.Get("/verify", userHandler.GetVerificationCode)

	handlers.TokenRoutes(rh, tokenService)

	protected := app.Group("/", authMiddleware)

	handlers.UserRoutes(rh, userService, protected)
	handlers.WalletRoutes(rh, walletService, dbPool, protected)
	handlers.ProductRoutes(rh, productService, cloudinaryService, dbPool, protected)

	rh.Logger.Info("Starting server", "server", "server")
	err := app.Listen(config.Port)
	rh.Logger.Error("error running server", "err", err)
}
