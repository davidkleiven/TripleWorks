package pkg

import (
	"context"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"reflect"
	"strings"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func Export(w io.Writer, items iter.Seq[models.MridGetter]) {
	for item := range items {
		ExportItem(w, item)
	}
}

func ExportItem(w io.Writer, item models.MridGetter) {
	itemType := reflect.TypeOf(item)
	itemValue := reflect.ValueOf(item)
	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
		itemValue = itemValue.Elem()
	}

	mrid := item.GetMrid()
	subject := fmt.Sprintf("<urn:uuid:%s>", mrid)
	fmt.Fprintf(w, "%s <%stype> <%s%s> .\n", subject, Rdf, Cim16, StructName(item))

	fieldWithoutIri := make(map[string]struct{})
	fields := FlattenStruct(item)
	uuidType := reflect.TypeOf(uuid.UUID{})
	for name, field := range fields {
		if field.Iri == "" {
			fieldWithoutIri[name] = struct{}{}
			continue
		}
		iri := field.Iri
		if strings.HasPrefix(iri, "cim:") {
			iri = strings.ReplaceAll(field.Iri, "cim:", Cim16)
		}
		if reflect.TypeOf(field.Value) == uuidType && name != "Mrid" {
			fmt.Fprintf(w, "%s <%s> <%s%s> .\n", subject, iri, "urn:uuid:", field.Value)
		} else {
			fmt.Fprintf(w, "%s <%s> \"%v\"%s .\n", subject, iri, field.Value, typeSpecifier(field.Value))
		}
	}
	slog.Info("Fields missing iris", "fields", fieldWithoutIri)
}

func typeSpecifier(value any) string {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "^^<http://www.w3.org/2001/XMLSchema#int>"
	case float64:
		return "^^<http://www.w3.org/2001/XMLSchema#float>"
	}
	return ""
}

func LatestOfAllItems(ctx context.Context, db *bun.DB, modelId int) ([]models.VersionedObject, error) {
	items := []models.VersionedObject{}
	for name, getter := range Finders {
		result, err := getter(ctx, db, modelId)
		if err != nil {
			return items, fmt.Errorf("Could not get data for %s: %w", name, err)
		}
		items = append(items, result...)
	}
	return items, nil
}
