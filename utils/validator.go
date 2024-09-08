package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// NotBlank validates that a string is not empty and does not consist only of whitespace
func NotBlank(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return strings.TrimSpace(value) != ""
}

// StructValidator is the custom validator struct
type StructValidator struct {
	Validator *validator.Validate
}

// NewStructValidator initializes a new StructValidator
func NewStructValidator() *StructValidator {
	// Create a new validator instance
	validator := validator.New()

	// Register the custom NotBlank validation
	validator.RegisterValidation("notblank", NotBlank)

	return &StructValidator{
		Validator: validator,
	}
}

// Validate implements the fiber.StructValidator interface
func (v *StructValidator) Validate(out interface{}) error {
	return v.Validator.Struct(out)
}
