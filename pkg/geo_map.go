package pkg

import (
	"io"
	"iter"
	"slices"
)

type SubstationMapData struct {
	Lat  float64
	Lng  float64
	Mrid string
	Name string
}

type LineMapData struct {
	LatFrom float64
	LatTo   float64
	LngFrom float64
	LngTo   float64
	Mrid    string
	Name    string
	Voltage int
}

func RenderMap(w io.Writer, substations iter.Seq[SubstationMapData], lines iter.Seq[LineMapData]) {
	tmpl := Map()
	data := struct {
		Substations []SubstationMapData
		Lines       []LineMapData
	}{
		Substations: slices.Collect(substations),
		Lines:       slices.Collect(lines),
	}
	PanicOnErr(tmpl.Execute(w, data))
}
