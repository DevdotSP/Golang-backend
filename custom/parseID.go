package custom

import (
	"strconv"

)

func ParseID(id string) (uint64, error) {
	personID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, err
	}
	return personID, nil
}