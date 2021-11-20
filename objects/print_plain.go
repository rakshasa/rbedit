package objects

import "fmt"

func PrintMapKeysAsPlain(v map[string]interface{}) {
	for key, _ := range v {
		// TODO: Not being escaped:
		fmt.Printf("%s\n", key)
	}
}

func PrintMapObjectKeysAsPlain(obj interface{}) error {
	v, ok := AsMap(obj)
	if !ok {
		return fmt.Errorf("object is not a map")
	}

	for key, _ := range v {
		// TODO: Not being escaped:
		fmt.Printf("%s\n", key)
	}

	return nil
}
