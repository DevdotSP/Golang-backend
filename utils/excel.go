package utils

import (
	"backend/model"

	"github.com/gofiber/fiber/v3"
	"github.com/xuri/excelize/v2"
)

// ExportToExcel exports the given data to an Excel file and writes it to the response.
func ExportToExcel(c fiber.Ctx, data interface{}) error {
	// Create a new Excel file
	f := excelize.NewFile()

	// Create a new sheet
	sheetName := "Sheet1"
	f.NewSheet(sheetName)

	// Set headers for the columns
	headers := []string{"ID", "Name", "Age", "Email", "Balance", "Action", "Created At"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1) // Row 1 for headers
		f.SetCellValue(sheetName, cell, header)
	}

	// Fill in the data
	rowIndex := 2
	persons, ok := data.([]model.User) // Assuming the data is of type []User
	if ok {
		for _, person := range persons {
			cellID, _ := excelize.CoordinatesToCellName(1, rowIndex)
			cellName, _ := excelize.CoordinatesToCellName(2, rowIndex)
			cellAge, _ := excelize.CoordinatesToCellName(3, rowIndex)
			cellEmail, _ := excelize.CoordinatesToCellName(4, rowIndex)
			cellBalance, _ := excelize.CoordinatesToCellName(5, rowIndex)
			cellAction, _ := excelize.CoordinatesToCellName(6, rowIndex)
			cellCreatedAt, _ := excelize.CoordinatesToCellName(7, rowIndex)

			f.SetCellValue(sheetName, cellID, person.ID)
			f.SetCellValue(sheetName, cellName, person.Name)
			if person.Age != nil {
				f.SetCellValue(sheetName, cellAge, *person.Age)
			}
			f.SetCellValue(sheetName, cellEmail, person.Email)
			f.SetCellValue(sheetName, cellBalance, person.AccountDetail.Balance)
			f.SetCellValue(sheetName, cellAction, person.History.Action)
			f.SetCellValue(sheetName, cellCreatedAt, person.History.CreatedAt)
			rowIndex++
		}
	}

	// Save the file to the server
	excelFilePath := `D:\Mobile Development\Golang\Backend\excel\persons.xlsx` // Specify the path
	if err := f.SaveAs(excelFilePath); err != nil {
		return err
	}

	// Set response header
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=\"persons.xlsx\"")

	// Write the file to the response
	if err := f.Write(c); err != nil {
		return err
	}

	return nil
}
