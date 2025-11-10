package middleware

import (
	"context"
	"ecommerce/internal/token"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type contextKey string

const ContextUserIDKey contextKey = "authenticatedUserID"

func NewContextWithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, ContextUserIDKey, userID)
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(ContextUserIDKey).(int64)
	return userID, ok
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}
		accessToken := headerParts[1]

		claims, err := token.VerifyAccessToken(accessToken)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired access token",
			})
		}

		newCtx := NewContextWithUserID(c.Context(), claims.UserID)
		c.SetUserContext(newCtx)
		c.Locals(ContextUserIDKey, claims.UserID)

		return c.Next()
	}
}
