package routes

import (
	"backend/controller"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// ProtectedRoutes initializes the protected routes for the Fiber app
func ProtectedRoutes(app *fiber.App, db *gorm.DB) {
	protected := app.Group("/api/protected",middleware.AuthMiddleware(),middleware.HeadersMiddleware())

	protected.Get("/single-data", controller.GetBranch(db))     // Get a single branch
	protected.Get("/all-data", controller.GetAllBranches(db)) // Get all branches
	protected.Post("/logout", controller.Logout())
	protected.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to the protected route!"})
	})
}
