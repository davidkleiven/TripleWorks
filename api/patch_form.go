package api

import (
	"net/http"

	"com.github/davidkleiven/tripleworks/pkg"
)

func PatchForm(w http.ResponseWriter, r *http.Request) {
	pkg.PatchForm(w)
}
