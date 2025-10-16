package main

import "github.com/gofiber/fiber/v2"

func (a *application) routes() *fiber.App {
	app := fiber.New()
	app.Get("/", handleRoot)
	app.Get("/v1/healthcheck", a.handler.Healthcheck)
	app.Post("/v1/users/register", a.handler.RegisterUserHandler)
	app.Post("/v1/users/login", a.handler.LoginUserHandler)

	return app
}
