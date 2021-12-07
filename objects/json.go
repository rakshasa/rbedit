package objects

import (
	"encoding/json"
	"fmt"
	"strings"
)

func sanitizeJSONToBencodeObject(object interface{}) error {
	if _, ok := object.(string); ok {
		return nil
	}
	if _, ok := object.(int64); ok {
		return nil
	}
	if l, ok := object.([]interface{}); ok {
		for idx, o := range l {
			if n, ok := o.(json.Number); ok {
				v, err := n.Int64()
				if err != nil {
					return fmt.Errorf("not an int64: %v", o)
				}

				l[idx] = v
				continue
			}

			if err := sanitizeJSONToBencodeObject(o); err != nil {
				return err
			}
		}
		return nil
	}
	if m, ok := object.(map[string]interface{}); ok {
		for key, o := range m {
			if n, ok := o.(json.Number); ok {
				v, err := n.Int64()
				if err != nil {
					return fmt.Errorf("not an int64: %v", o)
				}

				m[key] = v
				continue
			}

			if err := sanitizeJSONToBencodeObject(o); err != nil {
				return err
			}
		}
		return nil
	}

	return fmt.Errorf("not a valid bencode type: %T", object)
}

func ConvertJSONToBencodeObject(data string) (interface{}, error) {
	decoder := json.NewDecoder(strings.NewReader("[" + data + "]"))
	decoder.UseNumber()

	var object []interface{}
	if err := decoder.Decode(&object); err != nil {
		return nil, err
	}
	if len(object) != 1 {
		return nil, fmt.Errorf("getJSONValueFromFlag received an unexpected list result")
	}

	if err := sanitizeJSONToBencodeObject(object[0]); err != nil {
		return nil, err
	}

	return object[0], nil
}
