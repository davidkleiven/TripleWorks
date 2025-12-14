package pkg

import (
	"embed"
	"fmt"
	"io"

	"gonum.org/v1/gonum/graph/formats/rdf"
)

//go:embed resources/*
var resource embed.FS

const (
	Rdf   = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	Rdfs  = "http://www.w3.org/2000/01/rdf-schema#"
	Cim16 = "http://iec.ch/TC57/2013/CIM-schema-cim16#"
)

// Loads an rdf resource and groups statements by subject
func LoadObjects(r io.Reader) (*rdf.Graph, error) {
	graph := rdf.NewGraph()
	dec := rdf.NewDecoder(r)

	counter := 0
	for {
		stmt, err := dec.Unmarshal()
		if err == io.EOF {
			break
		} else if err != nil || stmt == nil {
			return graph, fmt.Errorf("Failed to read statement no. %d (%v): %w", counter, stmt, err)
		}
		counter++
		graph.AddStatement(stmt)

	}
	return graph, nil
}

type RdfsGraph struct {
	Graph *rdf.Graph
}

func (r *RdfsGraph) IsClass(term rdf.Term) bool {
	it := r.Graph.From(term.ID())
	classPred := "<" + Rdfs + "Class>"
	for it.Next() {
		node := it.Node()
		if otherTerm, ok := node.(rdf.Term); ok && otherTerm.Value == classPred {
			return true
		}
	}
	return false
}

// Extract all subjects having 'term' as their domain
func (r *RdfsGraph) Properties() *PropertyList {
	properties := NewPropertyList()
	rdfDomain := Must(rdf.NewIRITerm(Rdfs + "domain")).Value
	subClassOf := Must(rdf.NewIRITerm(Rdfs + "subClassOf")).Value
	iter := r.Graph.AllStatements()
	for iter.Next() {
		statement := iter.Statement()
		switch statement.Predicate.Value {
		case rdfDomain:
			properties.AddProperty(statement.Object, statement.Subject)
		case subClassOf:
			properties.SetSuperClass(statement.Subject, statement.Object)
		}
	}
	return properties
}
