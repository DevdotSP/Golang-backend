package generic

import (
	"backend/custom"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// Create a generic function for creating a resource
func CreateResource[T any](db *gorm.DB, input *T) fiber.Handler {
	return func(c fiber.Ctx) error {
		if err := c.Bind().Body(input); err != nil {
			log.Printf("Validation errors: %+v", err)
			return custom.SendErrorResponse(c, custom.NewHttpError(err.Error(), fiber.StatusBadRequest))
		}

		if err := db.Create(input).Error; err != nil {
			err := custom.NewHttpError("Could not create resource", fiber.StatusInternalServerError)
			return custom.SendErrorResponse(c, err)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Resource created successfully",
			"data":    input,
		})
	}
}

// Get all resources with optional preload
func GetAllResources[T any](db *gorm.DB, preloads []string) fiber.Handler {
	return func(c fiber.Ctx) error {
		var resources []T

		query := db
		for _, preload := range preloads {
			query = query.Preload(preload)
		}

		if err := query.Find(&resources).Error; err != nil {
			err := custom.NewHttpError("Could not retrieve resources", fiber.StatusInternalServerError)
			return custom.SendErrorResponse(c, err)
		}

		if len(resources) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "No records found",
			})
		}

		return c.JSON(resources)
	}
}

// Get a resource by ID with optional preload
func GetResourceByID[T any](db *gorm.DB, preloads []string) fiber.Handler {
	return func(c fiber.Ctx) error {
		var resource T
		id := c.Params("id")

		resourceID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Invalid ID", fiber.StatusBadRequest))
		}

		query := db
		for _, preload := range preloads {
			query = query.Preload(preload)
		}

		if err := query.First(&resource, resourceID).Error; err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Could not retrieve resource", fiber.StatusNotFound))
		}

		return c.JSON(resource)
	}
}

// Update a resource by ID
func UpdateResource[T any](db *gorm.DB, input *T) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		resourceID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Invalid ID", fiber.StatusBadRequest))
		}

		if err := db.Model(&input).Where("id = ?", resourceID).Updates(input).Error; err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Could not update resource", fiber.StatusInternalServerError))
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Resource updated successfully",
		})
	}
}

// DeleteResource deletes a resource by ID, with optional cascade delete for related records.
func DeleteResource[T any](db *gorm.DB, relatedModels ...interface{}) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		resourceID, err := custom.ParseID(id)
		if err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Invalid ID", fiber.StatusBadRequest))
		}

		// Handle cascading delete of related models if provided
		for _, relatedModel := range relatedModels {
			if err := db.Where("user_id = ?", resourceID).Delete(relatedModel).Error; err != nil {
				return custom.SendErrorResponse(c, custom.NewHttpError("Could not delete related records", fiber.StatusInternalServerError))
			}
		}

		// Delete the main resource
		if err := db.Delete(new(T), resourceID).Error; err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Could not delete resource", fiber.StatusInternalServerError))
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Resource deleted successfully",
		})
	}
}

// ExportToExcel exports any data structure to an Excel file and writes it to the response.
func ExportToExcel[T any](c fiber.Ctx, data []T, sheetName string, headers []string, columnWidths map[string]float64, dataMapper func(T) []interface{}) error {
	// Create a new Excel file
	f := excelize.NewFile()

	// Create a new sheet and handle the return values
	sheetIndex, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	// Set headers for the columns
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1) // Row 1 for headers
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return err
		}
	}

	// Set column width for readability
	for col, width := range columnWidths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return err
		}
	}

	// Fill in the data
	rowIndex := 2 // Start from the second row
	for _, item := range data {
		mappedData := dataMapper(item)
		for colIndex, value := range mappedData {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex) // Convert to Excel cell
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				return err
			}
		}
		rowIndex++
	}

	// Set the active sheet to the newly created sheet
	f.SetActiveSheet(sheetIndex)

	// Set response headers
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=\"data.xlsx\"")

	// Write the file to the response
	if err := f.Write(c); err != nil {
		return err
	}

	return nil
}
