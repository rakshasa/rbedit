package objects

import (
	"fmt"
)

func PrintObject(obj interface{}, options ...printOpFunction) {
	fmt.Println(SprintObject(obj, options...))
}

func PrintList(listObj []interface{}, options ...printOpFunction) {
	opts := NewPrintOptions(options)

	for _, uriStr := range sprintList(listObj, opts) {
		fmt.Printf("%s\n", uriStr)
	}
}

func PrintListObject(obj interface{}, options ...printOpFunction) error {
	listObj, ok := AsList(obj)
	if !ok {
		return fmt.Errorf("object is not a list")
	}

	PrintList(listObj, options...)
	return nil
}

func PrintMapObject(obj interface{}, options ...printOpFunction) error {
	m, ok := AsMap(obj)
	if !ok {
		return fmt.Errorf("object is not a map")
	}

	fmt.Printf("%s\n", sprintMap(m, NewPrintOptions(options)))
	return nil
}
