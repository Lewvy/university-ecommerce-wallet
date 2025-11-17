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
	cfg *config.Config,
	logger *slog.Logger,
	userService *service.UserService,
	tokenService *service.TokenService,
	walletService *service.WalletService,
	walletPaymentService *service.WalletPaymentService,
	productService *service.ProductService,
	dbPool *pgxpool.Pool,
) {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000/",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Authorization, Content-Length",
		AllowCredentials: true,
	}))

	rh := &rest.RestHandler{
		App:    app,
		Logger: logger,
	}

	userHandler := &handlers.UserHandler{Svc: userService}
	ph := &handlers.ProductHandler{Svc: productService}
	wph := &handlers.WalletPaymentHandler{Svc: walletPaymentService}

	rh.App.Post("/register", userHandler.RegisterUserHandler)
	rh.App.Post("/login", userHandler.LoginUserHandler)
	app.Post("/verify", userHandler.Verify)
	app.Get("/verify", userHandler.GetVerificationCode)
	app.Get("/products", ph.GetAllProductsHandler)

	app.Post("/wallet/webhook", wph.RazorpayWebhook)

	authMiddleware := middleware.AuthMiddleware(userService.Store)
	handlers.TokenRoutes(rh, tokenService)

	protected := app.Group("/", authMiddleware)

	handlers.UserRoutes(rh, userService, protected)
	handlers.WalletRoutes(rh, walletService, walletPaymentService, dbPool, protected)
	handlers.ProductRoutes(rh, productService, dbPool, userService, protected)

	rh.Logger.Info("Starting server", "server", "server")
	err := app.Listen(cfg.Port)
	rh.Logger.Error("error running server", "err", err)
}
