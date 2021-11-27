package objects

import (
	"fmt"
	"strings"
)

func PrintObject(obj interface{}, options ...printOpFunction) {
	strs := sprintObjectAsList(obj, NewPrintOptions(options))
	fmt.Printf("%s\n", strings.Join(strs, "\n"))
}

func PrintListObject(obj interface{}, options ...printOpFunction) error {
	l, ok := AsList(obj)
	if !ok {
		return fmt.Errorf("object is not a list")
	}

	fmt.Printf("%s\n", sprintList(l, NewPrintOptions(options)))
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
