package middleware

import "github.com/gofiber/fiber/v3"

// HeadersMiddleware sets common headers for responses
func HeadersMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Set security headers

		// Handle CORS
		c.Set("Access-Control-Allow-Origin", "*") // Change to your allowed origins
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Set Content-Type header for JSON responses
		c.Set("Content-Type", "application/json")

		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		// Continue to the next handler
		return c.Next()
	}
}
