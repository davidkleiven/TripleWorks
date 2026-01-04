package pkg

import (
	"embed"
	"fmt"
	"net/http"
	"strings"
)

type AutofillTargetField struct {
	Id       string `json:"id"`
	Label    string `json:"label"`
	Checksum string `json:"checksum"`
	Value    any    `json:"value"`
}

type FormState struct {
	Kind   string  `json:"kind"`
	Length float64 `json:"Length"`
	Name   string  `json:"Name"`
}

type AutofillInput struct {
	Fields []AutofillTargetField `json:"fields"`
	State  FormState             `json:"state"`
}

type FloatAutofiller func(a *FormState) float64
type StringAutofiller func(a *FormState) string

var floatAutofillers = map[string]FloatAutofiller{
	"R": func(a *FormState) float64 {
		rPerKm := 0.05 // Ohm/km
		return a.Length * rPerKm
	},
	"X": func(a *FormState) float64 {
		xPerKm := 0.4 // Ohm/km
		return a.Length * xPerKm
	},
	"Gch": func(a *FormState) float64 {
		gchPerKm := 0.05 // microF/km
		return gchPerKm * a.Length
	},
	"Bch": func(a *FormState) float64 {
		bchPerKm := 4.0e-5 // S/km
		return bchPerKm * a.Length
	},
}

var stringAutofillers = map[string]StringAutofiller{
	"ShortName": func(f *FormState) string {
		splitted := splitOnAny(f.Name, ",- ")
		result := ""
		for _, part := range splitted {
			if len(part) > 3 {
				result += part[:3]
			} else {
				result += part
			}
		}
		return result
	},
	"Description": func(f *FormState) string {
		return f.Name
	},
}

func splitOnAny(s string, separators string) []string {

	// Split using FieldsFunc
	return strings.FieldsFunc(s, func(r rune) bool {
		return strings.ContainsRune(separators, r)
	})
}

func GetAutofillValue(name string, state *FormState) (any, error) {
	filler, ok := floatAutofillers[name]
	if ok {
		return filler(state), nil
	}

	stringFiller, ok := stringAutofillers[name]
	if ok {
		return stringFiller(state), nil
	}
	return nil, fmt.Errorf("Could not find autofiller for '%s'", name)
}

//go:embed js/*
var jsFiles embed.FS

func JsServer() http.Handler {
	return http.FileServerFS(jsFiles)
}
