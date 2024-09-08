package custom

import (
	"github.com/gofiber/fiber/v3"
)


// Helper function to handle errors and send JSON responses
func SendErrorResponse(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}