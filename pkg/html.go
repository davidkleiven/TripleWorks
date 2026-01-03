package pkg

import (
	"embed"
	"io"
)

//go:embed html/*
var htmlPages embed.FS

func Index(w io.Writer) {
	reader := Must(htmlPages.Open("html/index.html"))
	io.Copy(w, reader)
}
