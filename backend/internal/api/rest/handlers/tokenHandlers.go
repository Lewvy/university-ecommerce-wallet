package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/dto"
	"ecommerce/internal/service"
	"ecommerce/internal/token"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenHandler struct {
	svc *service.TokenService
}

func TokenRoutes(rh *rest.RestHandler, tokenService *service.TokenService) {
	app := rh.App
	h := TokenHandler{
		svc: tokenService,
	}
	app.Post("/refresh", h.RefreshTokenHandler)
}

func (h *TokenHandler) RefreshTokenHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var input refreshTokenRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Code:    "bad_request",
			Message: "Invalid request body",
		})
	}

	if input.RefreshToken == "" {
		return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Code:    "bad_request",
			Message: "Refresh token must be provided",
		})
	}

	newAccessToken, newRefreshToken, err := h.svc.RefreshAndRevokeTokens(ctx, input.RefreshToken)

	if err != nil {
		if errors.Is(err, token.ErrTokenNotFound) || errors.Is(err, errors.New("refresh token expired")) {
			return c.Status(http.StatusForbidden).JSON(dto.ErrorResponse{
				Code:    "invalid_token",
				Message: "Invalid or expired refresh token.",
			})
		}

		if errors.Is(err, errors.New("internal server error: failed to revoke misused token")) {
			return c.Status(http.StatusInternalServerError).JSON(dto.ErrorResponse{
				Code:    "internal_error",
				Message: "Failed to process token securely.",
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "An internal server error has occurred",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"access_token":  newAccessToken.Plaintext,
		"refresh_token": newRefreshToken.Plaintext,
	})
}
