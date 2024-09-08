package routes

import (
	"backend/controller"
	"backend/middleware"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// SetupRoutes initializes the routes for the Fiber app
func SetupRoutes(app *fiber.App, db *gorm.DB) {

	// Group routes for persons under /api/person
	personGroup := app.Group("/api/person", middleware.HeadersMiddleware())
	{

		personGroup.Post("/", controller.CreatePerson(db))
		personGroup.Get("/", controller.GetAllPersons(db))
		personGroup.Get("/excel", controller.AllPersonExcel(db))
		personGroup.Get("/:id", controller.GetPerson(db))
		personGroup.Delete("/:id", controller.DeletePerson(db))
		personGroup.Post("/login", controller.Login(db))
		personGroup.Post("/logout", controller.Logout())
	}

	// Group routes for branches under /api/branch
	branchGroup := app.Group("/api/branch")
	{
		branchGroup.Get("/", controller.GetBranch(db))            // Get all branches
		branchGroup.Get("/info", controller.GetAllBranchData(db)) // Get all branches

	}

}
