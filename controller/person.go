package controller

import (
	"backend/custom"
	"backend/model"
	"backend/utils"
	"log"
	"math/rand"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// CreatePerson handles the creation of a new Person
func CreatePerson(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {

		var person model.User
		// Use BodyParser to bind the request body into the User model
		if err := c.Bind().Body(&person); err != nil {
			log.Printf("Binding error: %s", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
		}

		// Use custom validation utility
		if err := utils.ValidateAndRespond(c, &person); err != nil {
			return err // This already sends the error response
		}

		// Check if the user already exists
		var existingUser model.User
		if err := db.Where("email = ?", person.Email).First(&existingUser).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
		}

		// Use the utility function to hash the password
		hashedPassword, err := utils.HashPassword(person.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not create person",
			})
		}

		// Log the action in the history table
		historyEntry := model.History{
			UserId: person.ID,
			Action: "User created with name: " + person.Name,
		}
		if err := db.Create(&historyEntry).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Registered successfully",
		})
	}
}

// GetAllPersons handles the retrieval of all persons
func GetAllPersons(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error { // Correct signature
		var persons []model.User

		// Retrieve all persons from the database, preloading AccountDetail
		if err := db.Preload("AccountDetail").Preload("History").Find(&persons).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not retrieve persons",
			})
		}

		// Check if the persons slice is empty
		if len(persons) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "No records found",
			})
		}

		// Return the list of persons as JSON
		return c.JSON(persons)
	}
}

// GetAllPersons handles the retrieval of all persons
func GetPerson(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error { // Correct signature
		var persons []model.User

		// Get the person ID from the URL parameters
		id := c.Params("id")
		println(id)

		personID, err := custom.ParseID(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}

		// Retrieve all persons from the database
		if err := db.Preload("AccountDetail").Preload("History").Find(&persons, personID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not retrieve persons",
			})
		}

		// Check if the persons slice is empty
		if len(persons) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "No records found",
			})
		}

		// Return the list of persons as JSON
		return c.JSON(persons)
	}
}

// DeletePerson handles the deletion of a person by ID
func DeletePerson(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get the person ID from the URL parameters
		id := c.Params("id")
		println(id)

		// Convert the ID to an integer
		personID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}

		// Delete the person from the database
		if err := db.Delete(&model.User{}, personID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not delete person",
			})
		}

		// Return a success message
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Person deleted successfully",
		})
	}
}
