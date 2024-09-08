package controller

import (
	"backend/custom" // Import your custom utility package
	"backend/model"
	"backend/utils" // Import your JWT utility
	"log"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Login handles user authentication
func Login(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var userAuth model.UserLogin

		// Use BodyParser to bind the request body into the User model
		if err := c.Bind().Body(&userAuth); err != nil {
			log.Printf("Binding error: %s", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
		}

		// Use custom validation utility
		if err := utils.ValidateAndRespond(c, &userAuth); err != nil {
			return err // This already sends the error response
		}

		var user model.User
		if err := db.Where("email = ?", userAuth.Email).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return custom.SendErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
			}
			return custom.SendErrorResponse(c, fiber.StatusInternalServerError, "Could not find user")
		}

		// Compare the provided password with the hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userAuth.Password)); err != nil {
			return custom.SendErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
		}

		// Check if the current token is still active
		if existingToken, err := custom.ExtractToken(c); err == nil {
			if _, err := utils.ValidateToken(existingToken); err == nil {
				return custom.SendErrorResponse(c, fiber.StatusForbidden, "A valid token is already active.")
			}
		}

		// Generate a new JWT token
		token, err := utils.GenerateJWT(user.ID, "")
		if err != nil {
			return custom.SendErrorResponse(c, fiber.StatusInternalServerError, "Could not generate token")
		}

		// Log the generated token for debugging
		log.Printf("Generated token for user ID %d: %s", user.ID, token)

		// Return the generated token to the user
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Login successful",
			"token":   token,
		})
	}
}

// Logout handles user logout
func Logout() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get the token from the Authorization header
		jwtToken, err := custom.ExtractToken(c)
		if err != nil {
			return custom.SendErrorResponse(c, fiber.StatusInternalServerError, "Could not generate token")
		}

		// Delete the token from active tokens
		if err := utils.DeleteToken(jwtToken); err != nil {
			return custom.SendErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		}

		// Log the action for debugging
		log.Printf("Token %s has been marked as inactive.", jwtToken)

		// Invalidate token by instructing the client to remove it
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout successful. Please remove the token from storage."})
	}
}
