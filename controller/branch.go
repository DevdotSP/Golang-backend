package controller

import (
	"backend/generic"
	"backend/model"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// GetBranches retrieves all branches from the database
func GetBranch(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Define a slice to hold the results
		var results []struct {
			BranchCode string  `json:"branch_code"`
			BranchName string `json:"branch_name"`
		}

		// Query the database for branch codes and names from the JSONB column
		if err := db.Table("branches").
			Select("branch_data->>'branch_code' as branch_code,branch_data->>'branch_name' as branch_name").
			Scan(&results).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to fetch branches"})
		}

		// Create a simplified response with branch codes and names
		return c.JSON(results)
	}
}

// GetAllBranches uses the generic function to fetch all branches
func GetAllBranches(db *gorm.DB) fiber.Handler {
	return generic.GetAllResources[model.Branch](db, []string{})
}

// CreateBranch uses the generic CreateResource function to create a Branch
func CreateBranch(db *gorm.DB) fiber.Handler {
	var branch model.Branch
	return generic.CreateResource(db, &branch)
}
