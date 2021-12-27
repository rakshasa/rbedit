package outputs

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/rakshasa/bencode-go"
	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/types"
)

func NewEncodeBencode() types.EncodeFunc {
	return func(object interface{}) ([]byte, error) {
		var buf bytes.Buffer

		if err := bencode.Marshal(&buf, object); err != nil {
			return nil, fmt.Errorf("failed to encode data: %v", err)
		}

		return buf.Bytes(), nil
	}
}

func NewEncodePrint() types.EncodeFunc {
	return func(object interface{}) ([]byte, error) {
		return []byte(objects.SprintObject(object)), nil
	}
}

func NewEncodePrintList() types.EncodeFunc {
	return func(object interface{}) ([]byte, error) {
		stringList, err := SprintListOfStrings(object)
		if err != nil {
			return nil, fmt.Errorf("cannot print object: %v", err)
		}

		var str string
		for _, uri := range stringList {
			str += fmt.Sprintf("%s\n", uri)
		}

		return []byte(strings.TrimSuffix(str, "\n")), nil
	}
}

func NewEncodePrintAsListOfLists() types.EncodeFunc {
	return func(object interface{}) ([]byte, error) {
		parentList, ok := objects.AsList(object)
		if !ok {
			return nil, fmt.Errorf("cannot print object: not a list")
		}

		var str string
		for idx, childListObject := range parentList {
			stringList, err := SprintListOfStrings(childListObject)
			if err != nil {
				return nil, fmt.Errorf("cannot print object: %v", err)
			}

			for _, s := range stringList {
				str += fmt.Sprintf("%d: %s\n", idx, s)
			}
		}

		return []byte(strings.TrimSuffix(str, "\n")), nil
	}
}

func NewEncodeAsHexString() types.EncodeFunc {
	return func(object interface{}) ([]byte, error) {
		str, ok := objects.AsString(object)
		if !ok {
			return nil, fmt.Errorf("not a string")
		}

		return []byte(hex.EncodeToString([]byte(str))), nil
	}
}
