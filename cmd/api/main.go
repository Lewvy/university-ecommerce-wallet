package main

import (
	"ecommerce/cmd/api/handlers"
	"ecommerce/internal/config"
	"ecommerce/internal/logger"
	"log/slog"

	"github.com/joho/godotenv"
)

type application struct {
	cfg     config.Config
	logger  *slog.Logger
	handler *handlers.Handler
}

func main() {
	godotenv.Load()

	cfg := config.NewConfig()
	l := logger.NewLogger(cfg.Env)
	h := handlers.New(&cfg, l)

	app := &application{
		cfg:     cfg,
		logger:  l,
		handler: h,
	}

	app.server()
}
