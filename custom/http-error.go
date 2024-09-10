package custom

import (
	"github.com/gofiber/fiber/v3"
)

// HttpError represents a custom error type with a message and an error code.
type HttpError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NewHttpError creates a new instance of HttpError.
func NewHttpError(message string, code int) *HttpError {
	return &HttpError{
		Message: message,
		Code:    code,
	}
}

// SendErrorResponse sends a JSON error response.
func SendErrorResponse(c fiber.Ctx, err *HttpError) error {
	return c.Status(err.Code).JSON(fiber.Map{"error": err.Message})
}
