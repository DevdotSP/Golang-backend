package main

import (
	"log"

	"backend/database"
	"backend/utils"
	"backend/routes"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// Initialize the custom validator
	validator := utils.NewStructValidator()

	// Create a new Fiber app with the custom validator
	app := fiber.New(fiber.Config{
		StructValidator: validator, // Pass the initialized validator here
	})

	// Initialize the database connection
	db := database.InitDB()

	// Setup routes
	routes.SetupRoutes(app, db)
	routes.ProtectedRoutes(app, db)

	// Start the Fiber app
	log.Fatal(app.Listen(":3000"))
}
