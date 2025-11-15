package main

import (
	"ecommerce/internal/api"
	"ecommerce/internal/cache"
	"ecommerce/internal/config"
	"ecommerce/internal/data"
	"ecommerce/internal/data/gen"
	"ecommerce/internal/logger"
	"ecommerce/internal/mailer"
	"ecommerce/internal/service"
	"ecommerce/internal/worker"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	logger := logger.NewLogger()
	slog.SetDefault(logger)

	err := godotenv.Load()

	if err != nil {
		logger.Error("Error loading env", "err", err)
		return
	}

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("Error loading config", "error", err)
	}

	mailer, err := mailer.New(cfg.MailerHost, cfg.MailerPort, cfg.MailerUsername, cfg.MailerPassword, cfg.MailerSender)
	if err != nil {
		logger.Error("Error creating the mailer service", "error", err)
	}

	cacheClient, err := cache.NewValkeyCache()
	if err != nil {
		logger.Error("Error loading cache", "error", err)
		return
	}

	dbPool, err := data.NewDBPool(cfg.DBString)
	if err != nil {
		logger.Error("Error connecting to the database", "err", err.Error())
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
	userService := service.NewUserService(logger, userStore, walletStore, cacheClient, dbPool, tokenService)
	cloudService, err := service.NewCloudinaryService(&cfg, logger)

	rzrpay_id := os.Getenv("RAZORPAY_ID")
	rzrpay_secret := os.Getenv("RAZORPAY_SECRET")

	if rzrpay_id == "" || rzrpay_secret == "" {
		logger.Error("Error initializing razorpay creds")
	}

	walletService := service.NewWalletService(walletStore, dbPool, logger, rzrpay_id, rzrpay_secret)
	productService := service.NewProductService(productStore, cloudService, dbPool, logger)
	if err != nil {
		logger.Error("Error initializing cloudinaryService", "error", err)
		return
	}
	// authService := service.NewAuthService(logger, authStore, tokenService)

	api.SetupServer(&cfg, logger, userService, tokenService, walletService, productService, dbPool)

}
