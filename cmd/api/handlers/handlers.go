package handlers

import (
	"ecommerce/internal/config"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	cfg    *config.Config
	logger *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) *Handler {
	return &Handler{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *Handler) Healthcheck(c *fiber.Ctx) error {
	envelope := make(map[string]any)
	envelope["uptime"] = time.Since(a.cfg.StartTime).String()
	envelope["env"] = a.cfg.Env
	envelope["port"] = a.cfg.Port

	return c.Status(http.StatusOK).JSON(&envelope)
}
