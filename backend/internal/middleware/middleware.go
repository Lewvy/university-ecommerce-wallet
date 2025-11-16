package middleware

import (
	"ecommerce/internal/data"
	"ecommerce/internal/token"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const LocalsUserIDKey = "authenticatedUserID"

func AuthMiddleware(store data.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {

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

		// u, err := store.GetUserByID(c.Context(), int(claims.UserID))
		// if err != nil {
		// 	return err
		// }

		// if !u.EmailVerified {
		// 	err := c.Redirect("/verify", http.StatusUnauthorized)
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		c.Locals(LocalsUserIDKey, claims.UserID)

		return c.Next()
	}
}
