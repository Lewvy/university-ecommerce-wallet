package helpers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func ReadBody(c *fiber.Ctx, dst any) {
	json.Unmarshal(c.Body(), dst)
}
