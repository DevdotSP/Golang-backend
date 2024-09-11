package controller

import (
	"backend/custom" // Import your custom utility package
	"backend/model"
	"backend/utils" // Import your email utility
	"log"
	"math/rand"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// RegisterUser handles the registration of a new user
func RegisterUser(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var user model.User

		// Check if the body is empty
		body := c.Body()
		if len(body) == 0 {
			err := custom.NewHttpError("Request body is empty", fiber.StatusBadRequest)
			return custom.SendErrorResponse(c, err)
		}

		// Validate the user struct
		if err := c.Bind().Body(&user); err != nil {
			log.Printf("Validation errors: %+v", err)
			return custom.SendErrorResponse(c, custom.NewHttpError(err.Error(), fiber.StatusBadRequest))
		}

		// Validate the user struct using the custom validation
		if err := utils.Validator.Validate(&user); err != nil {
			log.Printf("Validation errors: %+v", err)
			return custom.SendErrorResponse(c, custom.NewHttpError(err.Error(), fiber.StatusBadRequest))
		}

		// Check if the user already exists
		var existingUser model.User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			err := custom.NewHttpError("Email already exists", fiber.StatusConflict)
			return custom.SendErrorResponse(c, err)
		}

		// Generate a verification token
		verificationToken := utils.GenerateVerificationToken() // Implement this function to create a secure token
		user.VerificationToken = verificationToken

		// Use the utility function to hash the password
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			err := custom.NewHttpError("Could not hash password", fiber.StatusInternalServerError)
			return custom.SendErrorResponse(c, err)
		}
		user.Password = hashedPassword // Store the hashed password

		// Create random balances
		balances := make([]float64, 3)
		val1 := rand.Float64() * 100 // Random value between 0 and 100
		balances[0] = val1

		// Create account detail
		accountDetail := model.AccountDetail{Balance: balances[0]} // Default balance
		user.AccountDetail = accountDetail

		// Insert the new user into the database
		if err := db.Create(&user).Error; err != nil {
			err := custom.NewHttpError("Could not create user", fiber.StatusInternalServerError)
			return custom.SendErrorResponse(c, err)
		}

		// Construct the verification link
		verificationLink := "http://127.0.0.1:3000/api/person/verify?token=" + verificationToken // Replace with your actual domain

		// Send the verification email
		emailBody := "Please verify your email by clicking the following link: " + verificationLink
		if err := utils.GoogleSendEmail(user.Email, "Email Verification", emailBody, verificationLink); err != nil {
			log.Printf("Could not send verification email: %v", err)
			return custom.SendErrorResponse(c, custom.NewHttpError("Could not send verification email", fiber.StatusInternalServerError))
		}

		// Log the action in the history table
		historyEntry := model.History{
			UserID: user.ID,
			Action: "User created with name: " + user.Name,
		}
		if err := db.Create(&historyEntry).Error; err != nil {
			err := custom.NewHttpError("Could not log history entry", fiber.StatusInternalServerError)
			return custom.SendErrorResponse(c, err)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Registered successfully, please check your email to verify your account",
		})
	}
}
