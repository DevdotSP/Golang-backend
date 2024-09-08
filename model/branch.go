package model

import (
	"gorm.io/datatypes"
)


type Location struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Manager struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

type Branch struct {
	ID         uint            `gorm:"primaryKey" json:"branch_id"`
	BranchData datatypes.JSON  `json:"branch_data"` // Store JSONB data
}