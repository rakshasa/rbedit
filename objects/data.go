package objects

import (
	"fmt"
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
	if !ok || !VerifyAbsoluteURI(s) {
		return "", false
	}

	return s, true
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
			return nil, fmt.Errorf("could not find key path object in map")
		}
	} else if l, ok := AsList(parentObj); ok {
		idx, err := stringToIntInRange(childKey, 0, len(l))
		if err != nil {
			return nil, fmt.Errorf("invalid key '%s', %v", childKey, err)
		}

		if childObj = l[idx]; !ok {
			return nil, fmt.Errorf("not find key path object in list")
		}
	} else {
		return nil, fmt.Errorf("non-terminating path element is not a key or list")
	}

	return LookupKeyPath(childObj, keys[1:])
}

func changeObject(parentObj interface{}, parentKeys []string, changeFn func(parentObj interface{}) (interface{}, error)) (interface{}, error) {
	if len(parentKeys) == 0 {
		return changeFn(parentObj)
	}

	childKey := parentKeys[0]
	if len(childKey) == 0 {
		return nil, fmt.Errorf("empty key path element")
	}

	if m, ok := AsMap(parentObj); ok {
		childObj, ok := m[childKey]
		if !ok {
			return nil, fmt.Errorf("could not find key path object in map")
		}

		childObj, err := changeObject(childObj, parentKeys[1:], changeFn)
		if err != nil {
			return nil, err
		}

		m[childKey] = childObj
		return m, nil

	} else if l, ok := AsList(parentObj); ok {
		idx, err := stringToIntInRange(childKey, 0, len(l))
		if err != nil {
			return nil, fmt.Errorf("invalid key '%s', %v", childKey, err)
		}

		childObj, err := changeObject(l[idx], parentKeys[1:], changeFn)
		if err != nil {
			return nil, err
		}

		l[idx] = childObj
		return l, nil

	} else {
		return nil, fmt.Errorf("non-terminating path element is not a key or list")
	}
}

// Returns the root object with the modified key path object.
func RemoveObject(rootObj interface{}, keys []string) (interface{}, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("empty keys")
	}

	lastKey := keys[len(keys)-1]
	if len(lastKey) == 0 {
		return nil, fmt.Errorf("empty last key")
	}

	return changeObject(rootObj, keys[:len(keys)-1], func(parentObj interface{}) (interface{}, error) {
		if m, ok := AsMap(parentObj); ok {
			delete(m, lastKey)
			return m, nil

		} else if l, ok := AsList(parentObj); ok {
			idx, err := stringToIntInRange(lastKey, 0, len(l))
			if err != nil {
				return nil, fmt.Errorf("invalid key '%s', %v", lastKey, err)
			}

			return append(l[:idx], l[idx+1:]...), nil

		} else {
			return nil, fmt.Errorf("non-terminating path element is not a key or list")
		}
	})
}

// Returns the root object with the modified key path object.
func SetObject(rootObj, setObj interface{}, keys []string) (interface{}, error) {
	if len(keys) == 0 {
		return setObj, nil
	}

	lastKey := keys[len(keys)-1]
	if len(lastKey) == 0 {
		return nil, fmt.Errorf("empty last key")
	}

	return changeObject(rootObj, keys[:len(keys)-1], func(parentObj interface{}) (interface{}, error) {
		if m, ok := AsMap(parentObj); ok {
			m[lastKey] = setObj
			return m, nil

		} else if l, ok := AsList(parentObj); ok {
			idx, err := stringToIntInRange(lastKey, 0, len(l))
			if err != nil {
				return nil, fmt.Errorf("invalid key '%s', %v", lastKey, err)
			}

			l[idx] = setObj
			return l, nil

		} else {
			return nil, fmt.Errorf("non-terminating path element is not a key or list")
		}
	})
}
