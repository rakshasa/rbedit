package objects

func AsMap(data interface{}) (map[string]interface{}, bool) {
	obj, ok := data.(map[string]interface{})
	return obj, ok
}
