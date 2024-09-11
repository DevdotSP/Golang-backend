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

		personGroup.Get("/verify", controller.VerifyEmail(db))
		personGroup.Post("/", controller.CreatePerson(db))
		personGroup.Get("/", controller.GetAllPersons(db))
		personGroup.Get("/excel", controller.ExportPersons(db))
		personGroup.Get("/:id", controller.GetPersonByID(db))
		personGroup.Put("/:id", controller.UpdatePerson(db))
		personGroup.Delete("/:id", controller.DeletePerson(db))

		personGroup.Post("/register", controller.RegisterUser(db))
		personGroup.Post("/login", controller.Login(db))
		personGroup.Post("/logout", controller.Logout())
	}

	// Group routes for branches under /api/branch
	branchGroup := app.Group("/api/branch")
	{
		branchGroup.Get("/", controller.GetBranch(db))
		branchGroup.Post("/", controller.CreateBranch(db))
		branchGroup.Delete("/:id", controller.DeleteBranch(db))
		branchGroup.Get("/info", controller.GetAllBranches(db))
	}

}
