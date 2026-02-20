package repository

import "embed"

//go:embed sql/*
var sqlFS embed.FS
