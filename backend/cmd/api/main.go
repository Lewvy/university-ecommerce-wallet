package main

import (
	"ecommerce/internal/api"
	"ecommerce/internal/cache"
	"ecommerce/internal/config"
	"ecommerce/internal/data"
	db "ecommerce/internal/data/gen"
	"ecommerce/internal/logger"
	"ecommerce/internal/mailer"
	"ecommerce/internal/service"
	"ecommerce/internal/worker"
	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	logger := logger.NewLogger()
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading env", "err", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("Error loading config", "error", err)
		return
	}

	mailer, err := mailer.New(cfg.MailerHost, cfg.MailerPort, cfg.MailerUsername, cfg.MailerPassword, cfg.MailerSender)
	if err != nil {
		logger.Error("Error creating mailer", "error", err)
		return
	}

	cacheClient, err := cache.NewValkeyCache()
	if err != nil {
		logger.Error("Error loading cache", "error", err)
		return
	}

	dbPool, err := data.NewDBPool(cfg.DBString)
	if err != nil {
		logger.Error("Database connection error", "err", err)
		return
	}
	defer dbPool.Close()

	cfg.DB = dbPool
	sqlcQueries := db.New(dbPool)

	workers := worker.NewWorkerPool(mailer, cacheClient, logger, true)
	workers.StartQueueMonitor()
	workers.StartEmailWorkers(1)

	userStore := data.NewUserStore(sqlcQueries)
	tokenStore := data.NewTokenStore(sqlcQueries)
	walletStore := data.NewWalletStore(sqlcQueries)
	productStore := data.NewProductStore(sqlcQueries)

	tokenService := service.NewTokenService(tokenStore, logger)
	walletPaymentService := service.NewWalletPaymentService(dbPool, walletStore, logger)
	cloudService, err := service.NewCloudinaryService(&cfg, logger)
	if err != nil {
		logger.Error("Cloudinary init error", "error", err)
		return
	}

	walletService := service.NewWalletService(walletStore, dbPool, walletPaymentService, logger)
	userService := service.NewUserService(logger, userStore, walletStore, cacheClient, dbPool, tokenService)
	productService := service.NewProductService(productStore, cloudService, dbPool, logger)

	api.SetupServer(
		&cfg,
		logger,
		userService,
		tokenService,
		walletService,
		walletPaymentService,
		productService,
		dbPool,
	)
}
