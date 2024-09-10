package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

// CustomValidator wraps the go-playground validator to implement Fiber's StructValidator
type CustomValidator struct {
	validator *validator.Validate
}

// Validate performs the validation and returns an error if validation fails
func (cv *CustomValidator) Validate(obj any) error {
	fmt.Printf("Input data: %+v\n", obj)

	if err := cv.validator.Struct(obj); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return fiber.NewError(fiber.StatusBadRequest, cv.formatValidationErrors(validationErrors))
		}
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input data")
	}
	return nil
}

// formatValidationErrors processes validation errors and generates a single string message
func (cv *CustomValidator) formatValidationErrors(errors validator.ValidationErrors) string {
	var validationMessages []string
	for _, err := range errors {
		validationMessages = append(validationMessages, cv.generateErrorMessage(err))
	}
	return strings.Join(validationMessages, ", ")
}

// generateErrorMessage creates a user-friendly error message based on the field and tag
func (cv *CustomValidator) generateErrorMessage(e validator.FieldError) string {
	fieldMessages := map[string]map[string]string{
		"Name": {
			"required": "Name is required",
			"min":      "Name must be at least " + e.Param() + " characters long",
			"max":      "Name must be no more than " + e.Param() + " characters long",
		},
		"Email": {
			"required": "Email is required",
			"email":    "Email must be a valid email address",
			"min":      "Password must be at least " + e.Param() + " characters long",
			"max":      "Password must be no more than " + e.Param() + " characters long",
		},
		"Password": {
			"required": "Password is required",
			"min":      "Password must be at least " + e.Param() + " characters long",
			"max":      "Password must be no more than " + e.Param() + " characters long",
		},
		"Age": {
			"required": "Age is required",
			"gte":      "Age must be greater than or equal to 18",
			"lte":      "Age must be less than or equal to 65",
		},
	}

	if field, exists := fieldMessages[e.Field()]; exists {
		if message, found := field[e.Tag()]; found {
			return message
		}
	}
	return fmt.Sprintf("Invalid value for %s", e.Field())
}

// NewCustomValidator creates a new instance of CustomValidator
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

// Initialize the validator
var Validator = NewCustomValidator()
