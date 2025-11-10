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
	// authStore := data.NewAuthStore(sqlcQueries)

	// walletStore := data.NewWalletStore(sqlcQueries)

	tokenService := service.NewTokenService(tokenStore)
	userService := service.NewUserService(logger, userStore, cacheClient, tokenService)
	// authService := service.NewAuthService(logger, authStore, tokenService)
	// walletService := service.NewWalletService(logger, walletStore)

	api.SetupServer(&cfg, logger, userService, tokenService)

}
