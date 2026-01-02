package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// APIKeyAuth middleware for simple API key authentication
func APIKeyAuth(apiKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if apiKey == "" {
			// Skip auth if no API key is configured
			return c.Next()
		}

		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		// Expected format: "Bearer <api-key>"
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}

		if parts[1] != apiKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		return c.Next()
	}
}

// RateLimiter middleware for basic rate limiting
func RateLimiter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Simple rate limiting logic can be implemented here
		// For production, consider using Redis-based rate limiting
		return c.Next()
	}
}