package objects

import (
	"fmt"

	"github.com/rakshasa/rbedit/types"
)

// Add float and uinteger:
func AsInteger(obj interface{}) (int64, bool) {
	d, ok := obj.(int64)
	return d, ok
}

func AsList(obj interface{}) ([]interface{}, bool) {
	l, ok := obj.([]interface{})
	return l, ok
}

func AsMap(obj interface{}) (map[string]interface{}, bool) {
	m, ok := obj.(map[string]interface{})
	return m, ok
}

func AsString(obj interface{}) (string, bool) {
	s, ok := obj.(string)
	return s, ok
}

func AsAbsoluteURI(obj interface{}) (string, bool) {
	s, ok := obj.(string)
	if !ok || !types.VerifyAbsoluteURI(s) {
		return "", false
	}

	return s, true
}

func CopyObject(src interface{}) (interface{}, error) {
	if d, ok := AsInteger(src); ok {
		return d, nil
	} else if l, ok := AsList(src); ok {
		dst := make([]interface{}, len(l))

		for idx, obj := range l {
			o, err := CopyObject(obj)
			if err != nil {
				return nil, err
			}

			dst[idx] = o
		}

		return dst, nil
	} else if m, ok := AsMap(src); ok {
		dst := make(map[string]interface{})

		for key, obj := range m {
			o, err := CopyObject(obj)
			if err != nil {
				return nil, err
			}

			dst[key] = o
		}
		return dst, nil
	} else if s, ok := AsString(src); ok {
		return s, nil
	} else if src == nil {
		return nil, fmt.Errorf("null object")
	} else {
		return nil, fmt.Errorf("invalid object type")
	}
}
