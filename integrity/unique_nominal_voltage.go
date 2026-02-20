package integrity

import (
	"encoding/json"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/pkg"
)

type UniqueNominalVoltage struct {
	BaseVoltages []models.BaseVoltage
}

func (u *UniqueNominalVoltage) Check() QualityResult {
	latest := pkg.OnlyActiveLatest(u.BaseVoltages)
	perNominalVoltage := pkg.GroupBy(latest, func(bv models.BaseVoltage) int { return int(bv.NominalVoltage) })
	for k, v := range perNominalVoltage {
		if len(v) == 1 {
			delete(perNominalVoltage, k)
		}
	}
	return &UniqueNominalVoltageResult{BaseVoltages: perNominalVoltage}
}

type UniqueNominalVoltageResult struct {
	BaseVoltages map[int][]models.BaseVoltage
}

func (u *UniqueNominalVoltageResult) Report(enc *json.Encoder) error {
	report := struct {
		Name         string                       `json:"name"`
		BaseVoltages map[int][]models.BaseVoltage `json:"base_voltages"`
	}{
		Name:         "One base voltage per nominal voltage",
		BaseVoltages: u.BaseVoltages,
	}
	return enc.Encode(report)
}
