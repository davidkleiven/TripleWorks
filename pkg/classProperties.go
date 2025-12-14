package pkg

import (
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
