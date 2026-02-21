package pkg

import (
	"fmt"
	"io"
	"iter"
	"regexp"
	"strings"

	"gonum.org/v1/gonum/graph/formats/rdf"
)

// PropertyList provides methods for efficiently getting all properties
// of a rdf.Term and the properties if its super classes
type PropertyList struct {
	properties map[rdf.Term][]rdf.Term
	superclass map[rdf.Term]rdf.Term
}

func NewPropertyList() *PropertyList {
	return &PropertyList{
		properties: make(map[rdf.Term][]rdf.Term),
		superclass: make(map[rdf.Term]rdf.Term),
	}
}

func (p *PropertyList) AddProperty(src rdf.Term, target rdf.Term) {
	items, ok := p.properties[src]
	if ok {
		items = append(items, target)
	} else {
		items = []rdf.Term{target}
	}
	p.properties[src] = items
}

func (p *PropertyList) SetSuperClass(src rdf.Term, target rdf.Term) {
	p.superclass[src] = target
}

func (p *PropertyList) GetProperties(term rdf.Term) []rdf.Term {
	directProps, ok := p.properties[term]
	if !ok {
		return []rdf.Term{}
	}

	superclass, ok := p.superclass[term]
	if !ok {
		return directProps
	}
	return append(directProps, p.GetProperties(superclass)...)
}

func (p *PropertyList) Classes() iter.Seq[rdf.Term] {
	return Keys(p.properties)
}

func (p *PropertyList) WriteBunModel(w io.Writer, term rdf.Term, params WriteBunModelParams) {
	props, ok := p.properties[term]
	if !ok {
		props = []rdf.Term{}
	}
	fmt.Fprintf(w, "type %s%s struct {\n", params.ClassPrefix, bunName(term.Value))

	super, ok := p.superclass[term]
	if ok {
		name := bunName(super.Value)
		// Embed super class
		fmt.Fprintf(w, "\t%s\n", name)
	}

	conflictHandler := newNameConflictHandler()
	for _, prop := range props {
		name := bunName(prop.Value)
		dbName := conflictHandler.Get(camelToSnake(name))
		iri := iriValue(prop.Value)
		dataType := params.Types.Get(iri)

		if strings.HasSuffix(dataType, "_enumeration") {
			dataType = strings.TrimSuffix(dataType, "_enumeration")
			splitted := MustSlice(strings.Split(dataType, "#"))
			dataType = splitted[len(splitted)-1]
			fmt.Fprintf(w, "\t%sId int `bun:\"%s_id\" json:\"%s_id\"`\n", name, dbName, dbName)
			fmt.Fprintf(w, "\t%s *%s `bun:\"rel:belongs-to,join:%s_id=id\" json:\"%s,omitempty\"`\n", name, dataType, dbName, dbName)
		} else if strings.Contains(dataType, "#") {
			// Class type
			fmt.Fprintf(w, "\t%sMrid uuid.UUID `bun:\"%s_mrid,type:uuid\" json:\"%s_mrid\"`\n", name, dbName, dbName)
			fmt.Fprintf(w, "\t%s *%s `bun:\"rel:belongs-to,join:%s_mrid=mrid\" json:\"%s,omitempty\"`\n", name, params.UuidType, dbName, dbName)
		} else {
			fmt.Fprintf(w, "\t%s %s `bun:\"%s\" json:\"%s\" iri:\"%s\"`\n", name, dataType, dbName, dbName, iri)
		}
	}
	fmt.Fprintf(w, "}\n")
}

func (p *PropertyList) WriteAllBunModels(w io.Writer, params WriteBunModelParams) {
	fmt.Fprintf(w, "package %s\n", params.Package)
	fmt.Fprintf(w, "import (\n")
	fmt.Fprintf(w, "\t\"time\"\n")
	fmt.Fprintf(w, "\t\"github.com/google/uuid\"\n")
	fmt.Fprintf(w, ")\n")
	seen := make(map[rdf.Term]struct{})

	if params.UuidType != "" {
		fmt.Fprintf(w, "type %s struct {\n", params.UuidType)
		fmt.Fprintf(w, "\tMrid uuid.UUID `bun:\"mrid,type:uuid\"`\n")
		fmt.Fprintf(w, "}\n")
	}

	for term := range p.properties {
		seen[term] = struct{}{}
		p.WriteBunModel(w, term, params)
	}

	addTerm := func(term rdf.Term) {
		_, ok := seen[term]
		if !ok {
			seen[term] = struct{}{}
			p.WriteBunModel(w, term, params)
		}
	}

	// Add super classes
	for baseclass, superclass := range p.superclass {
		addTerm(baseclass)
		addTerm(superclass)
	}
}

func bunName(value string) string {
	if strings.Contains(value, "#") {
		splitted := MustSlice(strings.Split(value, "#"))
		value = splitted[len(splitted)-1]
	}

	if strings.Contains(value, ".") {
		splitted := MustSlice(strings.Split(value, "."))
		value = splitted[len(splitted)-1]
	}
	value = strings.ReplaceAll(value, ">", "")
	return capitalizeFirst(value)
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func camelToSnake(s string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}

type Types struct {
	types map[string]string
}

func NewTypes(mapping map[string]string) *Types {
	return &Types{
		types: mapping,
	}
}

func (t *Types) Get(item string) string {
	value, ok := t.types[item]
	if !ok {
		return "string"
	}
	return value
}

type WriteBunModelParams struct {
	Types       Types
	ClassPrefix string
	UuidType    string
	Package     string
}

type nameConflictHandler struct {
	mapping map[string]string
}

func (n *nameConflictHandler) Get(name string) string {
	alternative, ok := n.mapping[name]
	if ok {
		return alternative
	}
	return name
}

func newNameConflictHandler() *nameConflictHandler {
	return &nameConflictHandler{
		mapping: map[string]string{
			"xmin": "x_min",
			"xmax": "x_max",
		},
	}
}
