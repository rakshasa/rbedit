package objects

// func PrintObjectAsPlain(obj interface{}, options ...printOpFunction) {
// 	fmt.Printf("%s\n", SprintObjectAsPlain(obj))
// }

// func SprintObjectAsPlain(obj interface{}, options ...printOpFunction) string {
// 	if _, ok := AsMap(obj); ok {
// 		return "<map-object>"
// 	} else if _, ok := AsList(obj); ok {
// 		return "<list-object>"
// 	} else {
// 		return "<unknown-object>"
// 	}
// }

// func PrintMapKeysAsPlain(v map[string]interface{}, options ...printOpFunction) {
// 	fmt.Printf("%s\n", strings.Join(MapKeysAsPlainList(v), "\n"))
// }

// func PrintMapObjectKeysAsPlain(obj interface{}, options ...printOpFunction) error {
// 	v, ok := AsMap(obj)
// 	if !ok {
// 		return fmt.Errorf("object is not a map")
// 	}

// 	PrintMapKeysAsPlain(v)
// 	return nil
// }
