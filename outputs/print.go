package outputs

import (
	"fmt"

	"github.com/rakshasa/rbedit/objects"
)

func SprintListOfStrings(obj interface{}) ([]string, error) {
	l, ok := objects.AsList(obj)
	if !ok {
		return nil, fmt.Errorf("cannot print object: not a list")
	}

	var strs []string

	for _, obj := range l {
		s, ok := objects.AsString(obj)
		if !ok {
			return nil, fmt.Errorf("cannot print object: not a string")
		}

		strs = append(strs, s)
	}

	return strs, nil
}
