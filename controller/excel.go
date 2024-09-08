package controller

import (
	"backend/model"
	"backend/utils"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// AllPersonExcel handles the retrieval of all persons and exports them to an Excel file.
func AllPersonExcel(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var persons []model.User // Change to a slice of User

		// Retrieve all persons from the database, preloading AccountDetail and History
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

		// Export data to Excel
		return utils.ExportToExcel(c, persons) // Pass the context and persons slice
	}
}
