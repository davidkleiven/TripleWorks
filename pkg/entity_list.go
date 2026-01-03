package pkg

import (
	"fmt"
	"html"
	"io"
	"reflect"
	"sort"

	"com.github/davidkleiven/tripleworks/models"
)

// CreateList writes an HTML table of the given items using Bulma CSS.
func CreateList(w io.Writer, items []models.MridNameGetter) {
	if len(items) == 0 {
		fmt.Fprintln(w, `<table class="table is-striped is-hoverable"><tr><td>No items</td></tr></table>`)
		return
	}

	// Use the first item's type to determine columns
	firstType := reflect.TypeOf(items[0])
	if firstType.Kind() == reflect.Ptr {
		firstType = firstType.Elem()
	}
	RequireStruct(firstType)

	// Collect field names
	fields := FlattenStruct(items[0])
	fieldNames := make([]string, 0, len(fields))
	for name, v := range fields {
		if !v.IsBunRelation && name != "Deleted" {
			fieldNames = append(fieldNames, name)
		}
	}

	// Optional: sort fields alphabetically
	sort.Strings(fieldNames)

	// Start table with Bulma classes
	fmt.Fprintln(w, `<table class="table is-striped is-hoverable is-fullwidth">`)

	// Table header
	fmt.Fprint(w, "<thead><tr>")
	for _, name := range fieldNames {
		fmt.Fprintf(w, "<th>%s</th>", html.EscapeString(name))
	}
	fmt.Fprintln(w, "</tr></thead>")

	// Table body
	fmt.Fprintln(w, "<tbody>")
	for _, item := range items {
		val := reflect.ValueOf(item)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		fmt.Fprint(w, "<tr>")
		fields := FlattenStruct(item)
		for _, name := range fieldNames {
			field := MustGet(fields, name)
			cell := fmt.Sprint(field.Value)
			// Escape HTML content
			fmt.Fprintf(w, "<td>%s</td>", html.EscapeString(cell))
		}
		fmt.Fprintln(w, "</tr>")
	}
	fmt.Fprintln(w, "</tbody>")
	fmt.Fprintln(w, "</table>")
}
