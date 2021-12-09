package objects

import (
	"fmt"
	"net/url"
	"strconv"
)

func VerifyAbsoluteURI(uri string) bool {
	u, err := url.Parse(uri)
	return err == nil && u.Scheme != "" && u.Host != ""
}

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
