package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func readBody(c *fiber.Ctx, dst any) {
	json.Unmarshal(c.Body(), dst)
}
