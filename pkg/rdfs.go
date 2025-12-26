package pkg

import (
	"embed"
	"fmt"
	"io"
	"strings"

	"gonum.org/v1/gonum/graph/formats/rdf"
)

//go:embed resources/*
var resource embed.FS

const (
	Rdf     = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	Rdfs    = "http://www.w3.org/2000/01/rdf-schema#"
	RdfsExt = "http://iec.ch/TC57/1999/rdf-schema-extensions-19990926#"
	Cim16   = "http://iec.ch/TC57/2013/CIM-schema-cim16#"
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

func (r *RdfsGraph) UnusedAssociations() map[rdf.Term]struct{} {
	associationUsed := Must(rdf.NewIRITerm(RdfsExt + "AssociationUsed")).Value
	iter := r.Graph.AllStatements()
	skip := make(map[rdf.Term]struct{})
	for iter.Next() {
		statement := iter.Statement()
		if statement.Predicate.Value == associationUsed && strings.Contains(strings.ToLower(statement.Object.Value), "no") {
			skip[statement.Subject] = struct{}{}
		}
	}
	return skip
}

// Extract all subjects having 'term' as their domain
func (r *RdfsGraph) Properties() *PropertyList {
	properties := NewPropertyList()
	rdfDomain := Must(rdf.NewIRITerm(Rdfs + "domain")).Value
	subClassOf := Must(rdf.NewIRITerm(Rdfs + "subClassOf")).Value

	skip := r.UnusedAssociations()
	iter := r.Graph.AllStatements()
	for iter.Next() {
		statement := iter.Statement()
		_, shouldSkip := skip[statement.Subject]
		if shouldSkip {
			continue
		}
		switch statement.Predicate.Value {
		case rdfDomain:
			properties.AddProperty(statement.Object, statement.Subject)
		case subClassOf:
			properties.SetSuperClass(statement.Subject, statement.Object)
		}
	}
	return properties
}

// Extract the closes golang type for all properties
func (r *RdfsGraph) GolangTypes() map[string]string {
	iter := r.Graph.AllStatements()
	dataType := Must(rdf.NewIRITerm(RdfsExt + "dataType")).Value
	rdfDomain := Must(rdf.NewIRITerm(Rdfs + "domain")).Value
	rdfDtype := make(map[rdf.Term]rdf.Term)
	domains := make(map[rdf.Term][]rdf.Term)
	for iter.Next() {
		statement := iter.Statement()
		switch statement.Predicate.Value {
		case rdfDomain:
			// Voltage.value <domain> Voltage
			// Used to traverse multiple levels of voltages
			v, ok := domains[statement.Object]
			if !ok {
				v = []rdf.Term{statement.Subject}
			} else {
				v = append(v, statement.Subject)
			}
			domains[statement.Object] = v
		case dataType:
			// <subject> dataType Voltage
			rdfDtype[statement.Subject] = statement.Object
		}
	}

	golangTypes := make(map[string]string)
	for subj, dtype := range rdfDtype {
		currentType := dtype
		for {
			alternatives, ok := domains[currentType]
			if ok {
				for _, candidate := range alternatives {
					newType, ok := rdfDtype[candidate]
					if ok {
						currentType = newType
						break
					}
				}
			} else {
				break
			}
		}

		var goDtype string
		rdfType := strings.ToLower(currentType.Value)
		if strings.Contains(rdfType, "float") {
			goDtype = "float64"
		} else if strings.Contains(rdfType, "int") {
			goDtype = "int"
		} else if strings.Contains(rdfType, "bool") {
			goDtype = "bool"
		} else {
			goDtype = "string"
		}
		golangTypes[subj.Value] = goDtype
	}
	return golangTypes
}
