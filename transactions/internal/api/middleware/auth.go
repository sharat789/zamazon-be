package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be-ms/transactions/internal/client"
)

func AuthorizeUser(authClient *client.AuthClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if len(token) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing authorization token",
			})
		}

		// Don't remove the Bearer prefix
		user, err := authClient.VerifyToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid or expired token",
			})
		}

		// Store user in context
		c.Locals("user", user)
		return c.Next()
	}
}

func AuthorizeByRole(authClient *client.AuthClient, role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if len(token) == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing authorization token",
			})
		}

		// Remove "Bearer " if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		user, err := authClient.AuthorizeByRole(token, role)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Insufficient permissions",
			})
		}

		c.Locals("user", user)
		return c.Next()
	}
}
