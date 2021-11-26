package objects

import (
	"fmt"
	"strconv"
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

func LookupKeyPath(parentObj interface{}, keys []string) (interface{}, error) {
	if len(keys) == 0 {
		return parentObj, nil
	}

	childKey := keys[0]
	if len(childKey) == 0 {
		return nil, fmt.Errorf("empty key path element")
	}

	var childObj interface{}

	if m, ok := AsMap(parentObj); ok {
		if childObj, ok = m[childKey]; !ok {
			return nil, fmt.Errorf("could not find key path object")
		}
	} else if l, ok := AsList(parentObj); ok {
		idx, err := strconv.Atoi(childKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert key path element to list index")
		}
		if idx < 0 || idx >= len(l) {
			return nil, fmt.Errorf("invalid list index number")
		}

		if childObj = l[idx]; !ok {
			return nil, fmt.Errorf("could not find key path object")
		}
	} else {
		return nil, fmt.Errorf("key path objects except the last must be a map or a list")
	}

	return LookupKeyPath(childObj, keys[1:])
}

// Returns the root object with the modified key path object.
func SetObject(parentObj, setObj interface{}, keys []string) (interface{}, error) {
	if len(keys) == 0 {
		return setObj, nil
	}

	childKey := keys[0]
	if len(childKey) == 0 {
		return nil, fmt.Errorf("empty key path element")
	}

	var childObj interface{}

	if m, ok := AsMap(parentObj); ok {
		if childObj, ok = m[childKey]; !ok {
			return nil, fmt.Errorf("could not find key path object")
		}

		childObj, err := SetObject(childObj, setObj, keys[1:])
		if err != nil {
			return nil, err
		}

		m[childKey] = childObj
		parentObj = m

	} else if l, ok := AsList(parentObj); ok {
		idx, err := strconv.Atoi(childKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert key path element to list index")
		}
		if idx < 0 || idx >= len(l) {
			return nil, fmt.Errorf("invalid list index number")
		}

		if childObj = l[idx]; !ok {
			return nil, fmt.Errorf("could not find key path object")
		}

		childObj, err := SetObject(childObj, setObj, keys[1:])
		if err != nil {
			return nil, err
		}

		l[idx] = childObj
		parentObj = l

	} else {
		return nil, fmt.Errorf("key path objects except the last must be a map or a list")
	}

	return parentObj, nil
}
