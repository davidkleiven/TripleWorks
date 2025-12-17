package api

import (
	"net/http"

	"com.github/davidkleiven/tripleworks/pkg"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	pkg.Index(w)
}

func Setup(mux *http.ServeMux) {
	mux.HandleFunc("/", RootHandler)
}
