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

func (p *PropertyList) WriteBunModel(w io.Writer, term rdf.Term, types Types, classPrefix string) {
	props, ok := p.properties[term]
	if !ok {
		props = []rdf.Term{}
	}
	fmt.Fprintf(w, "type %s%s struct {\n", classPrefix, bunName(term.Value))
	for _, prop := range props {
		name := bunName(prop.Value)
		dbName := camelToSnake(name)
		iri := iriValue(prop.Value)
		dataType := types.Get(iri)

		if strings.Contains(dataType, "#") {
			// Class type
			splitted := strings.Split(dataType, "#")
			className := iriValue(splitted[len(splitted)-1])
			className = capitalizeFirst(className)
			fmt.Fprintf(w, "\t%sId int `bun:\"%s_id\" json:\"%s_id\"`\n", name, dbName, dbName)
			fmt.Fprintf(w, "\t%s *%s `bun:\"rel:belongs-to,join:%s_id=id\" json:\"%s,omitempty\"`\n", name, className, dbName, dbName)
		} else {
			fmt.Fprintf(w, "\t%s %s `bun:\"%s\" json:\"%s\" iri:\"%s\"`\n", name, dataType, dbName, dbName, iri)
		}
	}

	super, ok := p.superclass[term]
	if ok {
		name := bunName(super.Value)
		dbName := camelToSnake(name)

		// Add Id
		fmt.Fprintf(w, "\t%sId int `bun:\"%s_id\" json:\"%s_id\"`\n", name, dbName, dbName)

		// Add relationship
		fmt.Fprintf(w, "\t%s *%s `bun:\"rel:belongs-to,join:%s_id=id\" json:\"%s,omitempty\" rdfs:\"subClassOf\"`\n", name, name, dbName, dbName)
	}
	fmt.Fprintf(w, "}\n")
}

func (p *PropertyList) WriteAllBunModels(w io.Writer, types Types, classPrefix string) {
	fmt.Fprint(w, "package tripleworks_rdf\n")
	fmt.Fprintf(w, "import (\n\t\"time\"\n)\n")
	seen := make(map[rdf.Term]struct{})
	for term := range p.properties {
		seen[term] = struct{}{}
		p.WriteBunModel(w, term, types, classPrefix)
	}

	addTerm := func(term rdf.Term) {
		_, ok := seen[term]
		if !ok {
			seen[term] = struct{}{}
			p.WriteBunModel(w, term, types, classPrefix)
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
		splitted := strings.Split(value, "#")
		value = splitted[len(splitted)-1]
	}

	if strings.Contains(value, ".") {
		splitted := strings.Split(value, ".")
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
