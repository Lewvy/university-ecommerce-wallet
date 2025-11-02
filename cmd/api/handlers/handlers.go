package handlers

import (
	"ecommerce/internal/config"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type envelope map[string]any

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

	var env envelope
	env["uptime"] = time.Since(a.cfg.StartTime).String()
	env["env"] = a.cfg.Env
	env["port"] = a.cfg.Port

	return c.Status(http.StatusOK).JSON(&env)
}

// func (a *Handler) UserSignup(c *fiber.Ctx) error {
// 	userInput := struct {
// 		Name     string `json:"name"`
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}{}
// 	err := c.BodyParser(&userInput)
// 	if err != nil {
// 		return err
// 	}
// 	v := validator.New()
// }
