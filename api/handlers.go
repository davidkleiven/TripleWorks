package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/pkg"
	"github.com/uptrace/bun"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	pkg.Index(w)
}

func CimTypes(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("resourceType")
	pkg.EntityOptions(w, target)
}

func EntityForm(w http.ResponseWriter, r *http.Request) {
	entityType := r.URL.Query().Get("type")
	item, err := pkg.FormInputFieldsForType(entityType)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to retrieve entity", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pkg.FormInputFields(w, item)
}

func Setup(mux *http.ServeMux, config *pkg.Config) {
	db := config.DatabaseConnection()
	mustPerformMigrations(db, 10*time.Minute)

	entityHandler := NewEntityStore(db, config.Timeout)
	mux.HandleFunc("/", RootHandler)
	mux.HandleFunc("/cim-types", CimTypes)
	mux.HandleFunc("/entity-form", EntityForm)
	mux.HandleFunc("GET /entity-form/{mrid}", entityHandler.EditComponentForm)
	mux.HandleFunc("GET /resource/{mrid}", entityHandler.Resource)
	mux.HandleFunc("/entities", entityHandler.GetEntityForKind)
	mux.HandleFunc("/enum", entityHandler.GetEnumOptions)
	mux.HandleFunc("/entity-list", entityHandler.EntityList)
	mux.HandleFunc("POST /commit", entityHandler.Commit)
	mux.HandleFunc("DELETE /commit/{id}", entityHandler.DeleteCommit)
	mux.HandleFunc("POST /autofill", AutofillHandler)
	mux.HandleFunc("GET /substations/{mrid}/diagram", entityHandler.SubstationDiagram)
	mux.HandleFunc("/export", entityHandler.Export)
	mux.HandleFunc("/upload/{kind}", entityHandler.SimpleUpload)
	mux.HandleFunc("GET /commits", entityHandler.Commits)

	mux.Handle("/js/", pkg.JsServer())
}

func mustPerformMigrations(db *bun.DB, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	executed, err := migrations.RunUp(ctx, db)

	if err != nil {
		panic(err)
	}
	slog.Info("Executed migrations", "num", len(executed.Migrations))
}
