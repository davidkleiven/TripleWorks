package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"com.github/davidkleiven/tripleworks/migrations"
	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/repository"
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
	timeout := 10 * time.Minute
	mustPerformMigrations(db, timeout)

	entityHandler := NewEntityStore(db, config.Timeout)
	inVoltageLevel := InVoltageLevelEndpoint{voltageLevelRepo: repository.NewBunVoltageLevelReadRepository(db), timeout: timeout}

	acLineRepo := repository.BunReadRepository[models.ACLineSegment]{Db: db, UseLatestView: true}
	substationRepo := repository.BunReadRepository[models.Substation]{Db: db, UseLatestView: true}
	equipmentInVoltageLevel := EquipmentInVoltageLevelEndpoint{
		VoltageLevel:    &repository.BunReadRepository[models.VoltageLevel]{Db: db, UseLatestView: true},
		ConNodes:        &repository.BunConnectivityNodeReadRepository{BunReadRepository: repository.BunReadRepository[models.ConnectivityNode]{Db: db, UseLatestView: true}},
		Terminals:       &repository.BunTerminalReadRepository{BunReadRepository: repository.BunReadRepository[models.Terminal]{Db: db, UseLatestView: true}},
		Generators:      &repository.BunReadRepository[models.SynchronousMachine]{Db: db, UseLatestView: true},
		Lines:           &acLineRepo,
		Switches:        &repository.BunReadRepository[models.Switch]{Db: db, UseLatestView: true},
		ConformLoads:    &repository.BunReadRepository[models.ConformLoad]{Db: db, UseLatestView: true},
		NonConformLoads: &repository.BunReadRepository[models.NonConformLoad]{Db: db, UseLatestView: true},
		Transformers:    &repository.BunReadRepository[models.PowerTransformer]{Db: db, UseLatestView: true},
		timeout:         timeout,
	}

	commit := CommitEndpoint{Db: &repository.BunInserter{Db: db}, timeout: timeout}
	validate := NewBunValidationEndpoint(db, timeout)
	xiidmEndpoint := XiidmExport{BusBreakerRepo: &repository.BunBusBreakerRepo{Db: db}, Timeout: timeout}
	userIdentifier := NoopMiddleware
	if config.WithTailscaleUserIdentification {
		tailscaleMiddleware := UserIdentificationMiddleware{
			Identifier: &TailscaleUserIdentifier{},
		}
		userIdentifier = tailscaleMiddleware.Apply
	}

	modelsEndpoint := ModelsEndpoint{
		Repo:    &repository.BunReadRepository[models.Model]{Db: db},
		Timeout: timeout,
	}

	substationWorkbench := SubstationConnectorWorkbench{
		LineRepo: &acLineRepo,
		Timeout:  timeout,
	}

	querySub := SubstationListQueryHandler{
		SubstationRepo: &substationRepo,
		Timeout:        timeout,
	}

	substationConnector := SubstationConnector{
		LineRepo:         &acLineRepo,
		SubstationRepo:   &substationRepo,
		TerminalRepo:     &repository.BunReadRepository[models.Terminal]{Db: db, UseLatestView: true},
		VoltageLevelRepo: &repository.BunReadRepository[models.VoltageLevel]{Db: db, UseLatestView: true},
		Inserter:         &repository.BunInserter{Db: db},
		Timeout:          timeout,
	}

	ptdfRecalc := RecalcPtdf{
		Bucket:            config.PtdfBucket,
		Doer:              &http.Client{},
		Model:             &repository.BunBusBreakerRepo{Db: db},
		PtdfEndpoint:      config.LoadflowServiceEndpoint + "/ptdf",
		PtdfWriterFactory: config.PtdfWriterFactory(),
		Timeout:           timeout,
	}

	actionForm := ActionFormEndpoint{Timeout: timeout}

	var ptdfs []pkg.PtdfRecord
	if config.PtdfProvider == "random" {
		slog.Info("Initializing random ptdfs")
		ptdfs = pkg.MustCreateRandomPtdf(&acLineRepo, &substationRepo)
	}
	flow := FlowEndpoint{
		Ptdf:        pkg.NewPtdfMatrix(ptdfs),
		MaxNumFlows: 100,
		Timeout:     timeout,
	}

	mux.HandleFunc("/", RootHandler)
	mux.HandleFunc("/cim-types", CimTypes)
	mux.HandleFunc("/entity-form", EntityForm)
	mux.HandleFunc("GET /entity-form/{mrid}", entityHandler.EditComponentForm)
	mux.HandleFunc("GET /resource/{mrid}", entityHandler.Resource)
	mux.Handle("GET /voltage-levels/{mrid}/items", &equipmentInVoltageLevel)
	mux.HandleFunc("GET /substations/{mrid}/voltage-levels", inVoltageLevel.VoltageLevelsInSubstation)
	mux.HandleFunc("/entities", entityHandler.GetEntityForKind)
	mux.HandleFunc("/enum", entityHandler.GetEnumOptions)
	mux.HandleFunc("/entity-list", entityHandler.EntityList)
	mux.Handle("POST /commit", userIdentifier(&commit))
	mux.HandleFunc("DELETE /commit/{id}", entityHandler.DeleteCommit)
	mux.HandleFunc("POST /autofill", AutofillHandler)
	mux.HandleFunc("GET /substations/{mrid}/diagram", entityHandler.SubstationDiagram)
	mux.HandleFunc("/export", entityHandler.Export)
	mux.Handle("/xiidm", &xiidmEndpoint)
	mux.HandleFunc("/upload/{kind}", entityHandler.SimpleUpload)
	mux.HandleFunc("GET /commits", entityHandler.Commits)
	mux.HandleFunc("/map", entityHandler.Map)
	mux.Handle("POST /connect-dangling", userIdentifier(http.HandlerFunc(entityHandler.ConnectDanglingLines)))
	mux.Handle("PATCH /resource", userIdentifier(http.HandlerFunc(entityHandler.ApplyJsonPatch)))
	mux.HandleFunc("/connection/{mrid}", entityHandler.Connection)
	mux.Handle("PUT /validate", validate)
	mux.Handle("/models", &modelsEndpoint)

	// Substation connection workkbench
	mux.Handle("POST /connect/{mrid}", &substationConnector)
	mux.Handle("GET /substation-connector/{mrid}", &substationWorkbench)
	mux.Handle("/substation-list", &querySub)
	mux.HandleFunc("/substation-selection", SetSelectedSubstation)
	mux.Handle("POST /ptdf/recalculate", &ptdfRecalc)
	mux.Handle("POST /production", &actionForm)
	mux.Handle("POST /flow", &flow)
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
