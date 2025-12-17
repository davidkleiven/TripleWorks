package pkg

import (
	"embed"
	"io"
)

//go:embed html/*
var html embed.FS

func Index(w io.Writer) {
	reader := Must(html.Open("html/index.html"))
	io.Copy(w, reader)
}
