package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

// ValidateAndRespond validates the provided data and sends a structured error response if validation fails.
func ValidateAndRespond(c fiber.Ctx, data any) error {
	if err := c.App().Config().StructValidator.Validate(data); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			fieldErrors := make(map[string]string)
			for _, fieldErr := range validationErrors {
				// Customize the error message for each field
				// Prefix the field name with its model name
				modelName := fieldErr.StructNamespace() // Get the full namespace (e.g., "User.Name")
				fieldErrors[modelName] = generateErrorMessage(fieldErr) // Use a helper function to generate error messages
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fieldErrors})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation error: " + err.Error()})
	}
	return nil
}

// Helper function to generate error messages
func generateErrorMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return fieldErr.Field() + " is required"
	case "email":
		return fieldErr.Field() + " must be a valid email address"
	case "notblank":
		return fieldErr.Field() + " cannot be blank"
	default:
		return "Invalid value for " + fieldErr.Field()
	}
}
