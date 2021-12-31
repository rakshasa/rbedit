package encodings

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/rakshasa/rbedit/data/templates"
	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/types"
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

func NewEncodePrint() types.EncodeFunc {
	return func(metadata types.IOMetadata, object interface{}) (types.IOMetadata, []byte, error) {
		return metadata, []byte(objects.SprintObject(object)), nil
	}
}

func NewEncodePrintList() types.EncodeFunc {
	return func(metadata types.IOMetadata, object interface{}) (types.IOMetadata, []byte, error) {
		stringList, err := SprintListOfStrings(object)
		if err != nil {
			return types.IOMetadata{}, nil, fmt.Errorf("cannot print object: %v", err)
		}

		var str string
		for _, uri := range stringList {
			str += fmt.Sprintf("%s\n", uri)
		}

		return metadata, []byte(strings.TrimSuffix(str, "\n")), nil
	}
}

func NewEncodePrintAsListOfLists() types.EncodeFunc {
	return func(metadata types.IOMetadata, object interface{}) (types.IOMetadata, []byte, error) {
		parentList, ok := objects.AsList(object)
		if !ok {
			return types.IOMetadata{}, nil, fmt.Errorf("cannot print object: not a list")
		}

		var str string
		for idx, childListObject := range parentList {
			stringList, err := SprintListOfStrings(childListObject)
			if err != nil {
				return types.IOMetadata{}, nil, fmt.Errorf("cannot print object: %v", err)
			}

			for _, s := range stringList {
				str += fmt.Sprintf("%d: %s\n", idx, s)
			}
		}

		return metadata, []byte(strings.TrimSuffix(str, "\n")), nil
	}
}

func NewEncodeAsHexString() types.EncodeFunc {
	return func(metadata types.IOMetadata, object interface{}) (types.IOMetadata, []byte, error) {
		str, ok := objects.AsString(object)
		if !ok {
			return types.IOMetadata{}, nil, fmt.Errorf("not a string")
		}

		return metadata, []byte(hex.EncodeToString([]byte(str))), nil
	}
}

func NewEncodePrintTemplate(templateText string) types.EncodeFunc {
	return func(metadata types.IOMetadata, object interface{}) (types.IOMetadata, []byte, error) {
		data, err := templates.ExecuteTemplate(metadata, templateText)
		if err != nil {
			return types.IOMetadata{}, nil, fmt.Errorf("could not print template: %v", err)
		}

		return metadata, []byte(data), nil
	}
}
