package main

import (
	"log"
	"os"

	"com.github/davidkleiven/tripleworks/pkg"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: rdfsgen <rdfs file> <out go file> (<class prefix>)")
	}

	schema := os.Args[1]
	out := os.Args[2]
	prefix := ""
	if len(os.Args) >= 4 {
		prefix = os.Args[3]
	}

	schemaFile, err := os.Open(schema)
	if err != nil {
		log.Fatalf("Could not open file %s\n", schema)
	}
	defer schemaFile.Close()

	graph, err := pkg.LoadObjects(schemaFile)
	if err != nil {
		log.Fatalf("Could not load objects: %s\n", err)
	}

	outfile, err := os.Create(out)
	if err != nil {
		log.Fatalf("Could not open outfile %s: %s\n", out, err)
	}
	defer outfile.Close()

	rdfsGraph := pkg.RdfsGraph{Graph: graph}
	properties := rdfsGraph.Properties()
	types := rdfsGraph.GolangTypes()
	properties.WriteAllBunModels(outfile, *pkg.NewTypes(types), prefix)
	log.Printf("Golang data models written to %s\n", out)
}
