package objects

import "fmt"

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

func LookupKeyPath(obj interface{}, keys []string) (interface{}, error) {
	if len(keys) == 0 {
		return obj, nil
	}
	if len(keys[0]) == 0 {
		return nil, fmt.Errorf("empty key path element")
	}

	m, ok := AsMap(obj)
	if !ok {
		return nil, fmt.Errorf("not a map")
	}

	child, ok := m[keys[0]]
	if !ok {
		return nil, fmt.Errorf("could not find child object")
	}

	return LookupKeyPath(child, keys[1:])
}

// Returns the root object with the modified key path object.
func SetKeyPath(parentObj, setObj interface{}, keys []string) (interface{}, error) {
	if len(keys) == 0 {
		return setObj, nil
	}

	childKey := keys[0]
	if len(childKey) == 0 {
		return nil, fmt.Errorf("empty key path element")
	}

	m, ok := AsMap(parentObj)
	if !ok {
		return nil, fmt.Errorf("not a map")
	}

	childObj, ok := m[childKey]
	if !ok {
		if len(keys) != 1 {
			return nil, fmt.Errorf("a key path object was not map")
		}
	}

	childObj, err := SetKeyPath(childObj, setObj, keys[1:])
	if err != nil {
		return nil, err
	}

	m[childKey] = childObj

	return m, nil
}
