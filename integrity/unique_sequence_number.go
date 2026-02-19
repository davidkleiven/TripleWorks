package integrity

import (
	"cmp"
	"encoding/json"
	"slices"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
	"github.com/google/uuid"
)

type UniqueSequenceNumberPerConductingEquipment struct {
	Terminals []models.Terminal
}

func (u *UniqueSequenceNumberPerConductingEquipment) Check() QualityResult {
	latest := pkg.OnlyActiveLatest(u.Terminals)
	grouped := pkg.GroupBy(latest, func(trm models.Terminal) uuid.UUID { return trm.ConductingEquipmentMrid })
	for k, terminals := range grouped {
		var (
			minSeq int
			maxSeq int
		)
		for _, term := range terminals {
			if term.SequenceNumber < minSeq || minSeq == 0 {
				minSeq = term.SequenceNumber
			}
			if term.SequenceNumber > maxSeq || maxSeq == 0 {
				maxSeq = term.SequenceNumber
			}
		}

		if minSeq == 1 && maxSeq == len(terminals) {
			// Valid
			delete(grouped, k)
		}
	}

	for k, term := range grouped {
		slices.SortFunc(term, func(a, b models.Terminal) int { return cmp.Compare(a.SequenceNumber, b.SequenceNumber) })
		grouped[k] = term
	}
	return &InvalidSequenceNumbers{
		Name:      "InvalidSequenceNumbers",
		Terminals: grouped,
	}
}

type InvalidSequenceNumbers struct {
	Name      string                          `json:"name"`
	Terminals map[uuid.UUID][]models.Terminal `json:"terminals"`
}

func (i *InvalidSequenceNumbers) Report(encoder *json.Encoder) error {
	return encoder.Encode(i)
}
