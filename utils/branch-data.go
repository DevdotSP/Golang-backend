package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/datatypes"
)

// generateRandomBranchData generates random JSONB data for a branch.
func GenerateRandomBranchData(branchID int) datatypes.JSON {
	data := map[string]interface{}{
		"branch_code": fmt.Sprintf("Branch Code %d", branchID),
		"branch_name": fmt.Sprintf("Branch Name %d", branchID),
		"address":     fmt.Sprintf("Address %d", rand.Intn(100)),
		"latitude":    rand.Float64()*180 - 90,  // Random latitude between -90 and 90
		"longitude":   rand.Float64()*360 - 180, // Random longitude between -180 and 180
		"opened":      time.Now().Format("2006-01-02"),
		"employees":   rand.Intn(100),
	}

	// Convert map to JSON
	jsonData, _ := json.Marshal(data)
	return datatypes.JSON(jsonData)
}
