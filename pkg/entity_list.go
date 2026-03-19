package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"com.github/davidkleiven/tripleworks/components"
)

func CreateList[T any](w io.Writer, items []T) {
	var (
		listItems []components.DrillableListItem
		errs      []error
	)
	for i, item := range items {
		name := StructName(item)
		jsonBytes, err := json.Marshal(item)
		if err != nil {
			errs = append(errs, fmt.Errorf("Failed to marshal item no %d: %w", i, err))
			continue
		}

		var mapRep map[string]any
		if err := json.Unmarshal(jsonBytes, &mapRep); err != nil {
			errs = append(errs, fmt.Errorf("Failed to unmarshal item no %d: %w", i, err))
			continue
		}

		listItems = append(listItems, components.DrillableListItem{Type: name, Content: mapRep})
	}

	listRep := components.MakeList(listItems)
	err := listRep.Render(context.Background(), w)
	errs = append(errs, err)

	if err := errors.Join(errs...); err != nil {
		slog.Warn("Failed to interpret some items", "error", err)
	}
}
