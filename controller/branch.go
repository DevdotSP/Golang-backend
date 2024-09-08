package controller

import (
	"backend/model"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// GetBranches retrieves all branches from the database
func GetBranch(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Define a slice to hold the results
		var results []struct {
			BranchCode string `json:"branch_code"`
			BranchName string `json:"branch_name"`
		}

		// Query the database for branch codes and names from the JSONB column
		if err := db.Table("branches").
			Select("branch_data->>'branch_code' as branch_code, branch_data->>'branch_name' as branch_name").
			Scan(&results).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to fetch branches"})
		}

		// Create a simplified response with branch codes and names
		return c.JSON(results)
	}
}

// GetAllPersons handles the retrieval of all persons
func GetAllBranchData(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error { // Correct signature
		var branches []model.Branch

		// Retrieve all persons from the database, preloading AccountDetail
		if err := db.Find(&branches).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not retrieve persons",
			})
		}

		// Check if the persons slice is empty
		if len(branches) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "No records found",
			})
		}

		// Return the list of persons as JSON
		return c.JSON(branches)
	}
}
