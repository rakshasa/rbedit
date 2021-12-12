package objects

import (
	"fmt"
	"strconv"
)

func stringToIntInRange(str string, min, max int) (int, error) {
	value, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("not an integer")
	}
	if value < min || value >= max {
		return 0, fmt.Errorf("out-of-bounds list index")
	}

	return value, nil
}
