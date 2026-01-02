package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestLogger adds request ID and structured logging
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate request ID
		requestID := uuid.New().String()
		c.Locals("requestID", requestID)
		c.Set("X-Request-ID", requestID)

		start := time.Now()

		// Process request
		err := c.Next()

		// Log request details
		duration := time.Since(start)
		log.Printf("[%s] %s %s - %d - %v - %s",
			requestID,
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			duration,
			c.IP(),
		)

		return err
	}
}