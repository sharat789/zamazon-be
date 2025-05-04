package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be-ms/auth/internal/api/handlers"
)

func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler) {
	// Auth routes
	authGroup := app.Group("/auth")

	authGroup.Post("/hash-password", authHandler.HashPassword)
	authGroup.Post("/verify-password", authHandler.VerifyPassword)
	authGroup.Post("/generate-token", authHandler.GenerateToken)
	authGroup.Post("/verify-token", authHandler.VerifyToken)
	authGroup.Post("/authorize-by-role", authHandler.AuthorizeByRole)
	authGroup.Get("/generate-code", authHandler.GenerateCode)

	// This endpoint can be used by other services to validate tokens
	authGroup.Get("/validate", authHandler.AuthMiddleware, func(c *fiber.Ctx) error {
		user := c.Locals("user")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Valid token",
			"user":    user,
		})
	})
}
