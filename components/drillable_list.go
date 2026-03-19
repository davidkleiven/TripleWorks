package components

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/a-h/templ"
)

type DrillableListItem struct {
	Type    string
	Content map[string]any
}

func orderedFields(items []DrillableListItem) []string {
	count := make(map[string]int)
	for _, item := range items {
		for field := range item.Content {
			count[field] += 1
		}
	}

	_, hasMrid := count["mrid"]
	_, hasName := count["name"]
	delete(count, "mrid")
	delete(count, "name")
	delete(count, "Commit")

	var result []string
	if hasMrid {
		result = append(result, "mrid")
	}
	if hasName {
		result = append(result, "name")
	}

	remaining := make([]string, 0, len(count))
	for k := range count {
		remaining = append(remaining, k)
	}

	slices.SortStableFunc(remaining, func(a, b string) int {
		ca, cb := count[a], count[b]
		if num := cmp.Compare(ca, cb); num != 0 {
			return -num
		}
		return cmp.Compare(a, b)
	})
	result = append(result, remaining...)
	if len(result) > 20 {
		result = result[:20]
	}
	return result
}

func MakeList(items []DrillableListItem) templ.Component {
	fields := orderedFields(items)
	return DrillableList(items, fields)
}

func asStringOrEmpty(content map[string]any, field string) string {
	if raw, ok := content[field]; ok {
		return fmt.Sprintf("%v", raw)
	}
	return ""
}
