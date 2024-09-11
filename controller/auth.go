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

		// Check if the body is empty
		body := c.Body()
		if len(body) == 0 {
			err := custom.NewHttpError("Request body is empty", fiber.StatusBadRequest)
			log.Printf("Validation errors: %+v", err)
			return custom.SendErrorResponse(c, err)
		}

		// Validate the userAuth struct
		if err := c.Bind().Body(&userAuth); err != nil {
			log.Printf("Validation errors: %+v", err)
			return custom.SendErrorResponse(c, custom.NewHttpError(err.Error(), fiber.StatusBadRequest))
		}

		// Validate the userAuth struct
		if err := utils.Validator.Validate(&userAuth); err != nil {
			log.Printf("Validation errors: %+v", err)
			return custom.SendErrorResponse(c, custom.NewHttpError(err.Error(), fiber.StatusBadRequest))
		}

		// Log the received request body
		log.Printf("Received request body: %s", body)

		log.Printf("Bound userAuth: %+v", userAuth)

		var user model.User
		if err := db.Where("email = ?", userAuth.Email).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err := custom.NewHttpError("Invalid invalid email", fiber.StatusBadRequest)
				return custom.SendErrorResponse(c, err)
			}
			err := custom.NewHttpError("Could not find user", fiber.StatusInternalServerError)
			return custom.SendErrorResponse(c, err)
		}

		// Compare the provided password with the hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userAuth.Password)); err != nil {
			err := custom.NewHttpError("Invalid password", fiber.StatusBadRequest)
			return custom.SendErrorResponse(c, err)
		}

		// Check if the current token is still active
		if existingToken, err := custom.ExtractToken(c); err == nil {
			if _, err := utils.ValidateToken(existingToken); err == nil {
				err := custom.NewHttpError("A valid token is already active", fiber.StatusForbidden)
				return custom.SendErrorResponse(c, err)
			}
		}

		// Generate a new JWT token
		token, err := utils.GenerateJWT(user.ID, "")
		if err != nil {
			err := custom.NewHttpError("Could not generate token", fiber.StatusInternalServerError)
			return custom.SendErrorResponse(c, err)
		}

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
			err := custom.NewHttpError("Could not extract token", fiber.StatusInternalServerError)
			return custom.SendErrorResponse(c, err)
		}

		// Delete the token from active tokens
		if err := utils.DeleteToken(jwtToken); err != nil {
			err := custom.NewHttpError("no token found", fiber.StatusInternalServerError)
			return custom.SendErrorResponse(c, err)
		}

		// Log the action for debugging
		log.Printf("Token %s has been marked as inactive.", jwtToken)

		// Invalidate token by instructing the client to remove it
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout successful. Please remove the token from storage."})
	}
}
