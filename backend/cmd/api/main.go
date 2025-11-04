package main

import (
	"ecommerce/internal/api"
	"ecommerce/internal/cache"
	"ecommerce/internal/config"
	"ecommerce/internal/data"
	"ecommerce/internal/data/gen"
	"ecommerce/internal/logger"
	"ecommerce/internal/service"
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

	cfg := config.NewConfig()

	cacheClient, err := cache.New()
	if err != nil {
		logger.Error("Error loading cache", "error", err)
		return
	}
	_ = cacheClient

	dbPool, err := data.NewDBPool(cfg.DBString)
	if err != nil {
		logger.Error("Error connecting to the database", "err", err.Error())
		return
	}
	defer dbPool.Close()
	cfg.DB = dbPool
	sqlcQueries := db.New(dbPool)

	userStore := data.NewUserStore(sqlcQueries)
	// walletStore := data.NewWalletStore(sqlcQueries)

	userService := service.NewUserService(logger, userStore)
	// walletService := service.NewWalletService(logger, walletStore)

	api.SetupServer(&cfg, logger, userService)

}
