package rbeditCmd

import (
	"fmt"
	"strconv"
)

func categoryIndexFromArgs(args []string) (int, error) {
	if len(args) != 1 {
		return 0, fmt.Errorf("command requires a category index argument")
	}

	category, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("category index argument is not an integer")
	}
	if category < 0 {
		return 0, fmt.Errorf("category index argument cannot be negative")
	}

	return category, nil
}

func categoryAndUrlIndicesFromArgs(args []string) (int, int, error) {
	if len(args) != 2 {
		return 0, 0, fmt.Errorf("command requires category and URI index arguments")
	}

	category, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, 0, fmt.Errorf("category index argument is not an integer")
	}
	if category < 0 {
		return 0, 0, fmt.Errorf("category index argument cannot be negative")
	}

	uri, err := strconv.Atoi(args[1])
	if err != nil {
		return 0, 0, fmt.Errorf("URI index argument is not an integer")
	}
	if uri < 0 {
		return 0, 0, fmt.Errorf("URI index argument cannot be negative")
	}

	return category, uri, nil
}
