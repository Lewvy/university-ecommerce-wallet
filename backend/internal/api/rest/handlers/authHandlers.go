package handlers

//
// import (
// 	"ecommerce/internal/api/rest"
// 	"ecommerce/internal/dto"
// 	"ecommerce/internal/service"
// 	"errors"
// 	"net/http"
//
// 	"github.com/gofiber/fiber/v2"
// )
//
// type AuthHandler struct {
// 	Svc *service.AuthService
// }
//
// func AuthRoutes(rh *rest.RestHandler, authService *service.AuthService) {
// 	app := rh.App
// 	h := AuthHandler{Svc: authService}
//
// 	app.Post("/login", h.LoginUserHandler)
// 	// app.Post("/refresh", h.RefreshTokenHandler)
// }
//
// func (h *AuthHandler) LoginUserHandler(c *fiber.Ctx) error {
// 	ctx := c.Context()
// 	var input dto.UserLogin
//
// 	err := c.BodyParser(&input)
// 	if err != nil {
// 		h.Svc.Logger.Error("Error decoding user", "err", err)
// 		return c.Status(http.StatusBadRequest).JSON(
// 			fiber.Map{
// 				"error": "invalid email/password format",
// 			},
// 		)
// 	}
//
// 	accessToken, refreshToken, err := h.Svc.Login(ctx, input)
//
// 	if err != nil {
// 		if errors.Is(err, service.ErrPwdMismatch) {
// 			h.Svc.Logger.Warn("User Validation Error", "error", service.ErrPwdMismatch.Error())
// 			return c.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
// 				Code:    "validation_error",
// 				Message: "invalid email or password",
// 			})
// 		}
//
// 		h.Svc.Logger.Error("Internal error during login", "err", err)
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "An internal server error has occurred",
// 		})
// 	}
//
// 	return c.Status(http.StatusOK).JSON(fiber.Map{
// 		"access_token":  accessToken,
// 		"refresh_token": refreshToken,
// 	})
// }
