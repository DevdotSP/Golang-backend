package controller

import (
	"backend/custom" // Import your custom utility package
	"backend/model"
	"backend/utils" // Import your JWT utility
	"log"
	"math/rand"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// CreatePerson handles the creation of a new Person
func RegisterUser(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var person model.User

		// Check if the body is empty
		body := c.Body()
		if len(body) == 0 {
			err := custom.NewHttpError("Request body is empty", fiber.StatusBadRequest)
			return custom.SendErrorResponse(c, err)
		}

		// Validate the userAuth struct
		if err := c.Bind().Body(&person); err != nil {
			log.Printf("Validation errors: %+v", err)
			return custom.SendErrorResponse(c, custom.NewHttpError(err.Error(), fiber.StatusBadRequest))
		}
// Validate the person struct using the custom validation
if err := utils.Validator.Validate(&person); err != nil {
	log.Printf("Validation errors: %+v", err)
	return custom.SendErrorResponse(c, custom.NewHttpError(err.Error(), fiber.StatusBadRequest))
}

		// Log the received request body
		log.Printf("Received request body: %s", body)

		log.Printf("Bound userAuth: %+v", person)

		// Check if the user already exists
		var existingUser model.User
		if err := db.Where("email = ?", person.Email).First(&existingUser).Error; err == nil {
			err := custom.NewHttpError("Email already exists", fiber.StatusConflict)
			return custom.SendErrorResponse(c, err)
		}

		// Use the utility function to hash the password
		hashedPassword, err := utils.HashPassword(person.Password)
		if err != nil {
			err := custom.NewHttpError("Could not hash password", fiber.StatusInternalServerError)
			return custom.SendErrorResponse(c, err)
		}
		person.Password = hashedPassword // Store the hashed password

		// Create random balances
		balances := make([]float64, 3)
		val1 := rand.Float64() * 100 // Random value between 0 and 100
		balances[0] = val1

		// Create account detail
		accountDetail := model.AccountDetail{Balance: balances[0]} // Default balance
		person.AccountDetail = accountDetail

		// Insert the new person into the database
		if err := db.Create(&person).Error; err != nil {
			err := custom.NewHttpError("Could not create person", fiber.StatusInternalServerError)
			return custom.SendErrorResponse(c, err)
		}

		// Log the action in the history table
		historyEntry := model.History{
			UserId: person.ID,
			Action: "User created with name: " + person.Name,
		}
		if err := db.Create(&historyEntry).Error; err != nil {
			err := custom.NewHttpError("Could not log history entry", fiber.StatusInternalServerError)
			return custom.SendErrorResponse(c, err)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Registered successfully",
		})
	}
}