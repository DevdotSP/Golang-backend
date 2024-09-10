package main

import (
	"backend/database"
	"backend/routes"
	"backend/utils"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	//cleanup expired token
	utils.StartCleanupRoutine()

	// Create a new Fiber app with the custom validator
	app := fiber.New(fiber.Config{
		StructValidator: utils.Validator, // Use the initialized custom Validator
	})

	// Initialize the database connection
	db := database.InitDB()

	// // Perform auto migration
	// db.AutoMigrate(
	// 	&model.User{},
	// 	&model.AccountDetail{},
	// 	&model.History{},
	// 	&model.Location{}, // Added Location model
	// 	&model.Manager{},  // Added Manager model
	// 	&model.Branch{},   // Added Branch model
	// )

	// // Insert 50-100 records
	// for i := 1; i <= 100; i++ {
	// 	branchData := utils.GenerateRandomBranchData(i)

	// 	branch := &model.Branch{
	// 		BranchData: branchData,
	// 	}

	// 	// Insert the branch record into the database
	// 	if err := db.Create(&branch).Error; err != nil {
	// 		log.Fatalf("failed to insert branch: %v", err)
	// 	}
	// 	fmt.Printf("Inserted branch ID: %d\n", branch.ID)
	// }

	// fmt.Println("Inserted 100 branches successfully.")

	// Setup routes
	routes.SetupRoutes(app, db)
	routes.ProtectedRoutes(app, db)

	// Start the Fiber app
	log.Fatal(app.Listen(":3000"))
}
