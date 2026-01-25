package migrations

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
)

func init() {
	migrations.MustRegister(createSequenceNumberConstraint, revertCreateSequenceNumberConstraint)
}

var sequenceNumberTables = []string{"terminals", "acdc_terminals", "dc_base_terminals"}

func createSequenceNumberConstraint(ctx context.Context, db *bun.DB) error {
	tmpl := template.Must(
		template.New("query").Parse(`
		ALTER TABLE {{.Table}} DROP CONSTRAINT IF EXISTS sequence_number_positive;
		ALTER TABLE {{.Table}}
		ADD CONSTRAINT sequence_number_positive
		CHECK (sequence_number > 0);
		`),
	)
	switch db.Dialect().Name() {
	case dialect.PG:
		for i, table := range sequenceNumberTables {
			data := struct{ Table string }{Table: table}
			var queryBuf bytes.Buffer
			tmpl.Execute(&queryBuf, data)
			query := queryBuf.String()
			_, err := db.Exec(query)
			if err != nil {
				return fmt.Errorf("Failed to run query %d (\n%s\n): %w", i, query, err)
			}
		}
	}
	return nil
}

func revertCreateSequenceNumberConstraint(ctx context.Context, db *bun.DB) error {
	if db.Dialect().Name() != dialect.PG {
		return nil
	}

	for i, table := range sequenceNumberTables {
		query := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT IF EXISTS sequence_number_positive", table)
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("Failed to run query %d (\n%s\n): %w", i, query, err)
		}
	}
	return nil
}
