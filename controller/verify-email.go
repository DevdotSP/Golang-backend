package controller

import (
	"backend/model"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// VerifyEmail handles the email verification process
func VerifyEmail(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Query("token")
		var user model.User

		// Find the user with the provided verification token
		if err := db.Where("verification_token = ?", token).First(&user).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid or expired verification token",
			})
		}

		// Update user as verified
		user.IsVerified = true
		user.VerificationToken = "" // Optionally, clear the token
		if err := db.Save(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not update user verification status",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Email successfully verified",
		})
	}
}

