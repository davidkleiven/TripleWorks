package pkg

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
)

//go:embed html/*
var htmlPages embed.FS

func Index(w io.Writer) {
	reader := Must(htmlPages.Open("html/index.html"))
	io.Copy(w, reader)
}

func toJSON(v any) template.JS {
	b, _ := json.Marshal(v)
	return template.JS(b)
}

func Map() *template.Template {
	return template.Must(template.New("map.html").Funcs(template.FuncMap{"toJSON": toJSON}).ParseFS(htmlPages, "html/map.html"))
}
