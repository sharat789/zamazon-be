// File: auth/internal/api/handlers/authHandler.go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be-ms/auth/internal/service"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type HashPasswordRequest struct {
	Password string `json:"password"`
}

type HashPasswordResponse struct {
	HashedPassword string `json:"hashed_password"`
}

func (h *AuthHandler) HashPassword(c *fiber.Ctx) error {
	var req HashPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	hashedPassword, err := h.authService.CreateHashPassword(req.Password)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(HashPasswordResponse{
		HashedPassword: hashedPassword,
	})
}

type VerifyPasswordRequest struct {
	PlainPassword  string `json:"plain_password"`
	HashedPassword string `json:"hashed_password"`
}

func (h *AuthHandler) VerifyPassword(c *fiber.Ctx) error {
	var req VerifyPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	err := h.authService.VerifyPassword(req.PlainPassword, req.HashedPassword)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Password verification failed",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Password verified successfully",
	})
}

type GenerateTokenRequest struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type GenerateTokenResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) GenerateToken(c *fiber.Ctx) error {
	var req GenerateTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	token, err := h.authService.GenerateToken(req.ID, req.Email, req.Role)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(GenerateTokenResponse{
		Token: token,
	})
}

type VerifyTokenRequest struct {
	Token string `json:"token"`
}

func (h *AuthHandler) VerifyToken(c *fiber.Ctx) error {
	var req VerifyTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	user, err := h.authService.VerifyToken(req.Token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token verification failed",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Token verified successfully",
		"user":    user,
	})
}

func (h *AuthHandler) GenerateCode(c *fiber.Ctx) error {
	code, err := h.authService.GenerateCode()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate verification code",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"code": code,
	})
}

type AuthorizeByRoleRequest struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func (h *AuthHandler) AuthorizeByRole(c *fiber.Ctx) error {
	var req AuthorizeByRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	user, err := h.authService.VerifyToken(req.Token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token verification failed",
			"error":   err.Error(),
		})
	}

	err = h.authService.AuthorizeByRole(user, req.Role)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"message": "Authorization failed",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User authorized successfully",
		"user":    user,
	})
}

// Middleware function that can be used by other services
func (h *AuthHandler) AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	user, err := h.authService.VerifyToken(authHeader)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authentication failed",
			"error":   err.Error(),
		})
	}

	c.Locals("user", user)
	return c.Next()
}
