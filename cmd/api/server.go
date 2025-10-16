package main

import (
	"github.com/gofiber/fiber/v2"
)

func (a *application) server() {
	app := a.routes()
	a.logger.Info("starting server", "port", a.cfg.Port)

	app.Listen(a.cfg.Port)
}

func handleRoot(c *fiber.Ctx) error {
	c.SendString("Hello World")
	return nil
}
