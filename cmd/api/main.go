package main

import (
	"database/sql"
	"ecommerce/cmd/api/handlers"
	"ecommerce/internal/api"
	"ecommerce/internal/config"
	"ecommerce/internal/logger"

	// "ecommerce/internal/logger"
	"log/slog"

	"github.com/joho/godotenv"
)

type application struct {
	Cfg     config.Config
	Logger  *slog.Logger
	Handler *handlers.Handler
	DB      *sql.DB
}

func main() {
	godotenv.Load()

	cfg := config.NewConfig()
	api.SetupServer(&cfg, logger.NewLogger("dev"))

	// cfg := config.NewConfig()
	// l := logger.NewLogger(cfg.Env)
	// h := handlers.New(&cfg, l)

	// app := &Application{
	// 	Cfg:     cfg,
	// 	Logger:  l,
	// 	Handler: h,
	// }

}
