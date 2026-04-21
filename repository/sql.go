package repository

import (
	"embed"
	"path/filepath"
)

//go:embed sql/*
var sqlFS embed.FS

func MustGetQuery(name string) string {
	query, err := sqlFS.ReadFile(filepath.Join("sql", name))
	if err != nil {
		panic(err)
	}
	return string(query)
}
