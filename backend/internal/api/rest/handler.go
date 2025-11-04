package rest

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type RestHandler struct {
	App    *fiber.App
	Logger *slog.Logger
}
