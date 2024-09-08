package middleware

import (
	"backend/utils"
	"log"

	"github.com/gofiber/fiber/v3"
)

// AuthMiddleware checks for a valid JWT token
// AuthMiddleware checks for a valid JWT token
func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get the token from the Authorization header
		token := c.Get("Authorization")

		// Check if the token is present and extract the Bearer token
		if token == "" || len(token) < 7 || token[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No token provided"})
		}

		// Extract the actual token
		jwtToken := token[7:]

		// Log the token being checked
		log.Printf("Checking token: %s", jwtToken)

		// Validate the token
		claims, err := utils.ValidateToken(jwtToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		// Store user claims in context for later use
		c.Locals("claims", claims)

		return c.Next()
	}
}
