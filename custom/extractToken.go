package custom

import (
	"github.com/gofiber/fiber/v3"
)

// Helper function to extract the JWT from the Authorization header
func ExtractToken(c fiber.Ctx) (string, error) {
	token := c.Get("Authorization")
	if token == "" || len(token) < 7 || token[:7] != "Bearer " {
		return "", fiber.NewError(fiber.StatusUnauthorized, "No token provided")
	}
	return token[7:], nil
}
