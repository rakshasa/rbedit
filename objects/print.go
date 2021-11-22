package objects

import (
	"fmt"
	"strings"
)

func PrintObject(obj interface{}, options ...printOpFunction) {
	opts := NewPrintOptions(options)
	printObject(obj, opts)
}

func PrintListObject(obj interface{}, options ...printOpFunction) error {
	opts := NewPrintOptions(options)
	return printMapObject(obj, opts)
}

func PrintMapObject(obj interface{}, options ...printOpFunction) error {
	opts := NewPrintOptions(options)
	return printMapObject(obj, opts)
}

func printObject(obj interface{}, opts *printOptions) {
	strs := sprintObjectAsList(obj, opts)
	fmt.Printf("%s\n", strings.Join(strs, "\n"))
}

func sprintObjectAsList(obj interface{}, opts *printOptions) []string {
	if d, ok := AsInteger(obj); ok {
		return sprintInteger(d, opts)
	} else if l, ok := AsList(obj); ok {
		return sprintList(l, opts)
	} else if m, ok := AsMap(obj); ok {
		return sprintMap(m, opts)
	} else if s, ok := AsString(obj); ok {
		return sprintString(s, opts)
	} else {
		return []string{"<unknown-object>"}
	}
}

func printListObject(obj interface{}, opts *printOptions) error {
	m, ok := AsList(obj)
	if !ok {
		return fmt.Errorf("object is not a list")
	}

	printList(m, opts)
	return nil
}

func printMapObject(obj interface{}, opts *printOptions) error {
	m, ok := AsMap(obj)
	if !ok {
		return fmt.Errorf("object is not a map")
	}

	printMap(m, opts)
	return nil
}

func sprintInteger(d int64, opts *printOptions) []string {
	return []string{fmt.Sprintf("%d", d)}
}

func printList(l []interface{}, opts *printOptions) {
	fmt.Printf("%s\n", sprintList(l, opts))
}

func sprintList(l []interface{}, opts *printOptions) []string {
	var strs []string

	indent := strings.Repeat(" ", opts.indent)

	for idx, obj := range l {
		objStrs := sprintObjectAsList(obj, opts)

		strs = append(strs, fmt.Sprintf("%d:", idx))
		for _, s := range objStrs {
			strs = append(strs, fmt.Sprintf("%s%s", indent, s))
		}
	}

	return strs
}

func printMap(m map[string]interface{}, opts *printOptions) {
	fmt.Printf("%s\n", sprintMap(m, opts))
}

func sprintMap(m map[string]interface{}, opts *printOptions) []string {
	indent := strings.Repeat(" ", opts.indent)

	keys := make([]string, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}

	var strs []string

	for _, key := range keys {
		if opts.keysOnly {
			strs = append(strs, fmt.Sprintf("%s", key))
			continue
		}

		objStrs := sprintObjectAsList(m[key], opts)

		strs = append(strs, fmt.Sprintf("%s:", key))
		for _, s := range objStrs {
			strs = append(strs, fmt.Sprintf("%s%s", indent, s))
		}
	}

	return strs
}

func sprintString(s string, opts *printOptions) []string {
	// TODO: Add prettify options.
	var str string

	for idx, c := range []byte(s) {
		if idx >= 256 {
			str += " ..."
			break
		}

		if c < 0x20 || c >= 0x7f {
			str += fmt.Sprintf("\\x%02x", int(c))
		} else {
			str += string(c)
		}
	}

	return []string{str}
}
