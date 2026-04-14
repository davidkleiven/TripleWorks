package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"com.github/davidkleiven/tripleworks/components"
)

type ActionFormEndpoint struct {
	Timeout time.Duration
}

func (a *ActionFormEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		slog.Error("Could not parse form", "error", err)
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	actionType := r.FormValue("action")
	mrid := r.FormValue("mrid")

	formItems := make([]components.NamedInjection, 0)

	for k, v := range r.Form {
		if len(v) > 1 && k != "mrid" && k != "action" && k != "name" {
			name, value := v[0], v[1]
			formItems = append(formItems, components.NamedInjection{
				Mrid:  k,
				Name:  name,
				Value: value,
			})
		}
	}

	if actionType == "delete" {
		filtered := make([]components.NamedInjection, 0)
		for _, item := range formItems {
			if item.Mrid != mrid {
				filtered = append(filtered, item)
			}
		}
		formItems = filtered
	} else {
		name := r.FormValue("name")
		formItems = append(formItems, components.NamedInjection{
			Mrid:  mrid,
			Name:  name,
			Value: "1.0",
		})
	}

	comp := components.ActionFormItems(formItems)
	ctx, cancel := context.WithTimeout(r.Context(), a.Timeout)
	defer cancel()
	w.Header().Set("HX-Trigger-After-Swap", "action-form-changed")
	comp.Render(ctx, w)
}
