package helpers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func getCurrentUserID(c *fiber.Ctx) (int32, error) {
	userID64, ok := c.Locals("authenticatedUserID").(int64)

	if !ok || userID64 == 0 {
		return 0, errors.New("unauthenticated or missing user ID in context")
	}
	return int32(userID64), nil
}
