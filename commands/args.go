package commands

import (
	"flag"
	"strconv"
)

func categoryIndexFromArgs(f *flag.FlagSet) (int, ExitStatusError) {
	if f.NArg() != 1 {
		return 0, &exitUsageError{msg: "command requires a category index argument"}
	}

	category, err := strconv.Atoi(f.Arg(0))
	if err != nil {
		return 0, &exitUsageError{msg: "category index argument is not an integer"}
	}
	if category < 0 {
		return 0, &exitUsageError{msg: "category index argument cannot be negative"}
	}

	return category, nil
}

func categoryAndUrlIndicesFromArgs(f *flag.FlagSet) (int, int, error) {
	if f.NArg() != 2 {
		return 0, 0, &exitUsageError{msg: "command requires category and URI index arguments"}
	}

	category, err := strconv.Atoi(f.Arg(0))
	if err != nil {
		return 0, 0, &exitUsageError{msg: "category index argument is not an integer"}
	}
	if category < 0 {
		return 0, 0, &exitUsageError{msg: "category index argument cannot be negative"}
	}

	uri, err := strconv.Atoi(f.Arg(1))
	if err != nil {
		return 0, 0, &exitUsageError{msg: "URI index argument is not an integer"}
	}
	if uri < 0 {
		return 0, 0, &exitUsageError{msg: "URI index argument cannot be negative"}
	}

	return category, uri, nil
}
