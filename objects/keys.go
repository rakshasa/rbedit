package objects

import (
	"fmt"

	"github.com/rakshasa/rbedit/types"
)

func changeObject(parentObj interface{}, parentKeys []string, changeFn func(parentObj interface{}) (interface{}, error)) (interface{}, error) {
	if len(parentKeys) == 0 {
		return changeFn(parentObj)
	}

	childKey := parentKeys[0]
	if len(childKey) == 0 {
		return nil, types.NewKeysLookupError("empty key path element", []string{childKey})
	}

	if m, ok := AsMap(parentObj); ok {
		childObj, ok := m[childKey]
		if !ok {
			return nil, types.NewKeysLookupError("key not found in map", []string{childKey})
		}

		childObj, err := changeObject(childObj, parentKeys[1:], changeFn)
		if err != nil {
			return nil, types.PrependKeyStringIfKeysError(err, childKey)
		}

		m[childKey] = childObj
		return m, nil

	} else if l, ok := AsList(parentObj); ok {
		idx, err := stringToIntInRange(childKey, 0, len(l))
		if err != nil {
			return nil, types.NewKeysLookupError(fmt.Sprintf("invalid key '%s', %v", childKey, err), []string{childKey})
		}

		childObj, err := changeObject(l[idx], parentKeys[1:], changeFn)
		if err != nil {
			return nil, types.PrependKeyStringIfKeysError(err, childKey)
		}

		l[idx] = childObj
		return l, nil

	} else {
		return nil, types.NewKeysLookupError("non-terminating path element is not a key or list", []string{childKey})
	}
}

func LookupKeyPath(parentObj interface{}, keys []string) (interface{}, error) {
	if len(keys) == 0 {
		return parentObj, nil
	}

	childKey := keys[0]
	if len(childKey) == 0 {
		return nil, types.NewKeysLookupError("empty key path element", []string{childKey})
	}

	var childObj interface{}

	if m, ok := AsMap(parentObj); ok {
		if childObj, ok = m[childKey]; !ok {
			return nil, types.NewKeysLookupError("could not find key path object in map", []string{childKey})
		}

	} else if l, ok := AsList(parentObj); ok {
		idx, err := stringToIntInRange(childKey, 0, len(l))
		if err != nil {
			return nil, types.NewKeysLookupError(fmt.Sprintf("invalid list index, %v", err), []string{childKey})
		}

		childObj = l[idx]

	} else {
		return nil, types.NewKeysLookupError("not a key or list", []string{childKey})
	}

	result, err := LookupKeyPath(childObj, keys[1:])
	return result, types.PrependKeyStringIfKeysError(err, childKey)
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
				return nil, types.NewKeysLookupError(fmt.Sprintf("invalid list index, %v", err), []string{lastKey})
			}

			return append(l[:idx], l[idx+1:]...), nil

		} else {
			return nil, fmt.Errorf("non-terminating path element is not a key or list")
		}
	})
}

// Returns the root object with the modified key path object.
func SetObject(rootObj, setObj interface{}, keys []string) (interface{}, error) {
	if setObj == nil {
		return nil, fmt.Errorf("null object")
	}
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
