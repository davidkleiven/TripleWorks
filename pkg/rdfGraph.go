package pkg

import (
	"strings"

	"gonum.org/v1/gonum/graph/formats/rdf"
)

type StatmentMatcher func(stmt *rdf.Statement) bool

func SubjectEndswith(suffix string) StatmentMatcher {
	return func(stmt *rdf.Statement) bool {
		return strings.HasSuffix(stmt.Subject.Value, suffix)
	}
}

func FindFirstMatch(g *rdf.Graph, matcher StatmentMatcher) *rdf.Statement {
	it := g.AllStatements()

	for it.Next() {
		stmt := it.Statement()
		if matcher(stmt) {
			return stmt
		}
	}
	return nil
}
