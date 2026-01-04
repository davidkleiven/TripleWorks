package migrations

import (
	"context"
	"fmt"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
)

func init() {
	migrations.MustRegister(addConNodeToTerminals, revertAddConNodeToTerminals)
}

func addConNodeToTerminals(ctx context.Context, db *bun.DB) error {
	// Extract ids of existing terminsals
	var mrids []string
	err := db.NewSelect().Table("terminals").Column("mrid").Scan(ctx, &mrids)
	if err != nil {
		return fmt.Errorf("Could not extract existing terminals: %w", err)
	}

	conNodes := make([]models.ConnectivityNode, len(mrids))
	defaults := make(map[string]string)
	for i, terminal_mrid := range mrids {
		mrid, err := uuid.NewUUID()
		if err != nil {
			return fmt.Errorf("Failed to create mrid: %w", err)
		}

		conNodes[i].Name = fmt.Sprintf("CN %d", i)
		conNodes[i].Mrid = mrid
		defaults[terminal_mrid] = mrid.String()
	}

	if len(conNodes) > 0 {
		_, err = db.NewInsert().Model(&conNodes).Exec(ctx)
		if err != nil {
			return fmt.Errorf("Failed to insert new connectivity nodes: %w", err)
		}
	}

	zero := uuid.UUID{}

	switch db.Dialect().Name() {
	case dialect.PG:
		query := fmt.Sprintf("ALTER TABLE terminals ADD COLUMN IF NOT EXISTS connectivity_node_mrid UUID DEFAULT '%s'", zero)
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("Failed to execute alter query: %w", err)
		}
	case dialect.SQLite:
		query := fmt.Sprintf("ALTER TABLE terminals ADD COLUMN connectivity_node_mrid TEXT DEFAULT %s", zero)
		_, err := db.Exec(query)
		if isSQLiteDuplicateColumn(err) {
			err = nil
		}
	}

	entity := models.Terminal{}
	for terminal_mrid, con_node_mrid := range defaults {
		_, err = db.NewUpdate().
			Model(&entity).
			Set("connectivity_node_mrid = ?", con_node_mrid).
			Where("mrid = ?", terminal_mrid).
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("Failed to insert con_node_mrid %s for terminal %s: %w", con_node_mrid, terminal_mrid, err)
		}
	}

	if db.Dialect().Name() == dialect.PG {
		query := "ALTER TABLE terminals ADD CONSTRAINT fk_connectivity_node_mrid FOREIGN KEY(connectivity_node_mrid) REFERENCES entities(mrid)"
		_, err = db.Exec(query)
	}
	return err
}

func revertAddConNodeToTerminals(ctx context.Context, db *bun.DB) error {
	if db.Dialect().Name() == dialect.SQLite {
		return nil
	}
	query := "ALTER TABLE terminals DROP CONSTRAINT fk_connectivity_node_mrid"
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("Failed to drop constraint: %w", err)
	}

	query = "ALTER TABLE terminals DROP COLUMN connectivity_node_mrid"
	_, err = db.ExecContext(ctx, query)
	return err
}
