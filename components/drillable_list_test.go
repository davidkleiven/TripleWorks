package components

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllFieldsPresent(t *testing.T) {
	items := []DrillableListItem{
		{Content: map[string]any{"mrid": "000-000", "name": "my component", "voltage": 22}},
	}

	component := MakeList(items)
	var buf bytes.Buffer
	component.Render(context.Background(), &buf)
	asString := buf.String()
	for k := range items[0].Content {
		require.Contains(t, asString, k)
	}

	// Table head and one row
	require.Equal(t, 2, strings.Count(asString, "</tr>"))
}

func TestMaxFieldsAreReturned(t *testing.T) {
	content := make(map[string]any)
	for i := range 50 {
		content[fmt.Sprintf("f%d", i)] = i
	}

	items := []DrillableListItem{{Content: content}}
	fields := orderedFields(items)
	require.Equal(t, 20, len(fields))
}

func TestEmptyStringOnNonExsitingField(t *testing.T) {
	result := asStringOrEmpty(map[string]any{}, "whatever-key")
	require.Equal(t, "", result)
}
