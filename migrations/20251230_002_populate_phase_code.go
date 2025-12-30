package migrations

import (
	"context"
	"fmt"
	"strings"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

func init() {
	migrations.MustRegister(populatePhaseCodeEnum, revertPopulatePhaseCodeEnum)
}

func populatePhaseCodeEnum(ctx context.Context, db *bun.DB) error {
	enumData := []models.RdfsEnum{
		// PhaseCode
		{Id: 1, Code: "A", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.A", Comment: "Phase A."},
		{Id: 2, Code: "AB", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.AB", Comment: "Phases A and B."},
		{Id: 3, Code: "ABC", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.ABC", Comment: "Phases A, B, and C."},
		{Id: 4, Code: "ABCN", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.ABCN", Comment: "Phases A, B, C, and N."},
		{Id: 5, Code: "ABN", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.ABN", Comment: "Phases A, B, and neutral."},
		{Id: 6, Code: "AC", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.AC", Comment: "Phases A and C."},
		{Id: 7, Code: "ACN", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.ACN", Comment: "Phases A, C and neutral."},
		{Id: 8, Code: "AN", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.AN", Comment: "Phases A and neutral."},
		{Id: 9, Code: "B", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.B", Comment: "Phase B."},
		{Id: 10, Code: "BC", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.BC", Comment: "Phases B and C."},
		{Id: 11, Code: "BCN", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.BCN", Comment: "Phases B, C, and neutral."},
		{Id: 12, Code: "BN", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.BN", Comment: "Phases B and neutral."},
		{Id: 13, Code: "C", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.C", Comment: "Phase C."},
		{Id: 14, Code: "CN", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.CN", Comment: "Phases C and neutral."},
		{Id: 15, Code: "N", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.N", Comment: "Neutral phase."},
		{Id: 16, Code: "s1", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.s1", Comment: "Secondary phase 1."},
		{Id: 17, Code: "s12", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.s12", Comment: "Secondary phase 1 and 2."},
		{Id: 18, Code: "s12N", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.s12N", Comment: "Secondary phases 1, 2, and neutral."},
		{Id: 19, Code: "s1N", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.s1N", Comment: "Secondary phase 1 and neutral."},
		{Id: 20, Code: "s2", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.s2", Comment: "Secondary phase 2."},
		{Id: 21, Code: "s2N", Iri: "http://iec.ch/TC57/2013/CIM-schema-cim16#PhaseCode.s2N", Comment: "Secondary phase 2 and neutral."},
	}

	// Clean up comments
	for i := range enumData {
		comment := strings.TrimSuffix(enumData[i].Comment, "^^<http://www.w3.org/1999/02/22-rdf-syntax-ns#XMLLiteral>")
		comment = strings.Trim(comment, `"`)
		enumData[i].Comment = comment
	}

	phases := []models.PhaseCode{
		{RdfsEnum: enumData[0]},
		{RdfsEnum: enumData[1]},
		{RdfsEnum: enumData[2]},
		{RdfsEnum: enumData[3]},
		{RdfsEnum: enumData[4]},
		{RdfsEnum: enumData[5]},
		{RdfsEnum: enumData[6]},
		{RdfsEnum: enumData[7]},
		{RdfsEnum: enumData[8]},
		{RdfsEnum: enumData[9]},
		{RdfsEnum: enumData[10]},
		{RdfsEnum: enumData[11]},
		{RdfsEnum: enumData[12]},
		{RdfsEnum: enumData[13]},
		{RdfsEnum: enumData[14]},
		{RdfsEnum: enumData[15]},
		{RdfsEnum: enumData[16]},
		{RdfsEnum: enumData[17]},
		{RdfsEnum: enumData[18]},
		{RdfsEnum: enumData[19]},
		{RdfsEnum: enumData[20]},
	}

	for _, phase := range phases {
		_, err := db.NewInsert().Model(&phase).Exec(ctx)
		if err != nil {
			return fmt.Errorf("Failed to insert PhaseCode enum: %w", err)
		}
	}

	return nil
}

func revertPopulatePhaseCodeEnum(ctx context.Context, db *bun.DB) error {
	_, err := db.NewDelete().
		Table("phase_codes").
		Where("1=1").
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("Failed to clear table phase_codes: %w", err)
	}

	return nil
}
