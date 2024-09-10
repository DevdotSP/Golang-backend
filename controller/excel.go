package controller

import (
	"backend/generic"
	"backend/model"
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// ExportPersons handles the export of users to an Excel file
func ExportPersons(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Retrieve the users from the database
		var persons []model.User
		if err := db.Preload("AccountDetail").Preload("History").Find(&persons).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to fetch users"})
		}
		// Log the number of persons retrieved
log.Printf("Number of persons retrieved: %d", len(persons))

		// Define headers and column widths
		headers := []string{"ID", "Name", "Age", "Email", "Balance", "Action", "Created At"}
		columnWidths := map[string]float64{
			"A": 10, // Width for "ID"
			"B": 25, // Width for "Name"
			"C": 10, // Width for "Age"
			"D": 30, // Width for "Email"
			"E": 15, // Width for "Balance"
			"F": 40, // Width for "Action"
			"G": 20, // Width for "Created At"
		}

		// Call the generic ExportToExcel utility
		if err := generic.ExportToExcel(c, persons, "Users", headers, columnWidths, func(person model.User) []interface{} {
			return []interface{}{
				person.ID,
				person.Name,
				person.Age, // Assuming Age is a pointer, handle nil appropriately
				person.Email,
				person.AccountDetail.Balance,
				person.History.Action,
				person.History.CreatedAt,
			}
		}); err != nil {
			return err
		}

		return nil
	}
}
