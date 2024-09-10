package controller

import (
	"backend/generic"
	"backend/model"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)


// GetAllPersons uses the generic GetAllResources function for retrieving all users
func GetAllPersons(db *gorm.DB) fiber.Handler {
	return generic.GetAllResources[model.User](db, []string{"AccountDetail", "History"})
}

// GetPersonByID uses the generic GetResourceByID function for retrieving a user by ID
func GetPersonByID(db *gorm.DB) fiber.Handler {
	return generic.GetResourceByID[model.User](db, []string{"AccountDetail", "History"})
}

// UpdatePerson uses the generic UpdateResource function for updating a user
func UpdatePerson(db *gorm.DB) fiber.Handler {
	var person model.User
	return generic.UpdateResource(db, &person)
}

func DeletePerson(db *gorm.DB) fiber.Handler {
	return generic.DeleteResource[model.User](db, &model.AccountDetail{}, &model.History{})
}
