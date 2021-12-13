package objects

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rakshasa/rbedit/types"
)

func SprintObject(obj interface{}, options ...printOpFunction) string {
	strs := sprintObjectAsList(obj, NewPrintOptions(options))
	return strings.Join(strs, "\n")
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

func sprintInteger(d int64, opts *printOptions) []string {
	return []string{fmt.Sprintf("%d", d)}
}

func sprintList(l []interface{}, opts *printOptions) []string {
	var strs []string

	indent := strings.Repeat(" ", opts.indent)

	for idx, obj := range l {
		objStrs := sprintObjectAsList(obj, opts)

		if !opts.valuesOnly {
			strs = append(strs, fmt.Sprintf("%d:", idx))
		}

		for _, s := range objStrs {
			strs = append(strs, fmt.Sprintf("%s%s", indent, s))
		}
	}

	return strs
}

func sprintMap(m map[string]interface{}, opts *printOptions) []string {
	indent := strings.Repeat(" ", opts.indent)

	keys := make([]string, 0, len(m))
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

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

func sprintString(str string, opts *printOptions) []string {
	if len(str) > 256 {
		return []string{types.EscapeURIString(str[:256]) + " ..."}
	}

	return []string{types.EscapeURIString(str)}
}
