package objects

import (
	"fmt"

	"github.com/rakshasa/rbedit/types"
)

// Add float and uinteger:
func AsInteger(obj interface{}, args ...bool) (int64, bool) {
	if len(args) != 0 && !args[0] {
		return 0, false
	}

	d, ok := obj.(int64)
	return d, ok
}

func AsList(obj interface{}, args ...bool) ([]interface{}, bool) {
	if len(args) != 0 && !args[0] {
		return nil, false
	}

	l, ok := obj.([]interface{})
	return l, ok
}

func AsMap(obj interface{}, args ...bool) (map[string]interface{}, bool) {
	if len(args) != 0 && !args[0] {
		return nil, false
	}

	m, ok := obj.(map[string]interface{})
	return m, ok
}

func AsString(obj interface{}, args ...bool) (string, bool) {
	if len(args) != 0 && !args[0] {
		return "", false
	}

	s, ok := obj.(string)
	return s, ok
}

func AsAbsoluteURI(obj interface{}, args ...bool) (string, bool) {
	if len(args) != 0 && !args[0] {
		return "", false
	}

	s, ok := obj.(string)
	if !ok || !types.VerifyAbsoluteURI(s) {
		return "", false
	}

	return s, true
}

// func LookupKey(obj interface{}, key string) (map[string]interface{}, bool) {
func LookupKey(obj interface{}, key string) (interface{}, bool) {
	m, ok := obj.(map[string]interface{})
	if !ok {
		// return map[string]interface{}, false
		return nil, false
	}

	v, ok := m[key]
	return v, ok
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
