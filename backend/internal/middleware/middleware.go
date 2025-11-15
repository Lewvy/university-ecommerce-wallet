package middleware

import (
	"ecommerce/internal/token"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const LocalsUserIDKey = "authenticatedUserID"

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		log.Printf("[AuthMiddleware] Running for path: %s", c.Path())

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Println("[AuthMiddleware] FAILED: Authorization header required")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			log.Println("[AuthMiddleware] FAILED: Invalid Authorization header format")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}

		accessToken := headerParts[1]

		claims, err := token.VerifyAccessToken(accessToken)
		if err != nil {
			log.Printf("[AuthMiddleware] FAILED: Invalid or expired access token. Error: %v", err)
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired access token",
			})
		}

		log.Printf("[AuthMiddleware] SUCCESS: Token verified for UserID: %d (Type: %T)", claims.UserID, claims.UserID)

		c.Locals(LocalsUserIDKey, claims.UserID)

		val := c.Locals(LocalsUserIDKey)
		log.Printf("[AuthMiddleware] Set c.Locals('%s'). Value is now: %v (Type: %T)", LocalsUserIDKey, val, val)

		return c.Next()
	}
}
