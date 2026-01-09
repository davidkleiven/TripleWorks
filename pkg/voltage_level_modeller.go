package pkg

import (
	"context"
	"fmt"
	"strings"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type VoltageLevelModel struct {
	ReportingGroup    models.ReportingGroup
	BusNameMarkers    []models.BusNameMarker
	ConnectivityNodes []models.ConnectivityNode
	Switches          []models.Switch
	Terminals         []models.Terminal
}

func (v *VoltageLevelModel) AssignCommitId(commitId int) {
	v.ReportingGroup.CommitId = commitId
	for i := range v.BusNameMarkers {
		v.BusNameMarkers[i].CommitId = commitId
	}

	for i := range v.ConnectivityNodes {
		v.ConnectivityNodes[i].CommitId = commitId
	}

	for i := range v.Switches {
		v.Switches[i].CommitId = commitId
	}

	for i := range v.Terminals {
		v.Terminals[i].CommitId = commitId
	}
}

func (v *VoltageLevelModel) Entities(modelId int) []models.Entity {
	entities := make([]models.Entity, 0, len(v.BusNameMarkers)+len(v.ConnectivityNodes)+len(v.Switches)+len(v.Terminals)+1)

	entities = append(entities, models.Entity{Mrid: v.ReportingGroup.Mrid, ModelEntity: models.ModelEntity{ModelId: modelId}})

	for _, bnm := range v.BusNameMarkers {
		entity := models.Entity{Mrid: bnm.Mrid, ModelEntity: models.ModelEntity{ModelId: modelId}}
		entities = append(entities, entity)
	}

	for _, conNode := range v.ConnectivityNodes {
		entity := models.Entity{Mrid: conNode.Mrid, ModelEntity: models.ModelEntity{ModelId: modelId}}
		entities = append(entities, entity)
	}

	for _, breaker := range v.Switches {
		entity := models.Entity{Mrid: breaker.Mrid, ModelEntity: models.ModelEntity{ModelId: modelId}}
		entities = append(entities, entity)
	}

	for _, terminal := range v.Terminals {
		entity := models.Entity{Mrid: terminal.Mrid, ModelEntity: models.ModelEntity{ModelId: modelId}}
		entities = append(entities, entity)
	}
	return entities
}

func (v *VoltageLevelModel) Write(ctx context.Context, db *bun.DB, modelId int, commitMsg string) error {
	entities := v.Entities(modelId)
	commit := models.Commit{
		Message: commitMsg,
		Author:  "VoltageLevelModeller",
	}

	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := ReturnOnFirstError(
			func() error {
				_, err := tx.NewInsert().Model(&commit).Exec(ctx)
				return err
			},
			func() error {
				v.AssignCommitId(int(commit.Id))
				return nil
			},
			func() error {
				_, err := tx.NewInsert().Model(&entities).Exec(ctx)
				return err
			},
			func() error {
				_, err := tx.NewInsert().Model(&v.BusNameMarkers).Exec(ctx)
				return err
			},
			func() error {
				_, err := tx.NewInsert().Model(&v.ConnectivityNodes).Exec(ctx)
				return err
			},
			func() error {
				_, err := tx.NewInsert().Model(&v.Switches).Exec(ctx)
				return err
			},
			func() error {
				_, err := tx.NewInsert().Model(&v.Terminals).Exec(ctx)
				return err
			},
			func() error {
				_, err := tx.NewInsert().Model(&v.ReportingGroup).Exec(ctx)
				return err
			},
		)
		return err
	})
}

type VoltageLevelEquipment struct {
	VoltageLevel        models.VoltageLevel
	Lines               []models.ACLineSegment
	Generators          []models.SynchronousMachine
	ConformLoads        []models.ConformLoad
	LineTerminalNumbers map[uuid.UUID]int
}

func NewVoltageLevelEquipment(opts ...func(v *VoltageLevelEquipment)) *VoltageLevelEquipment {
	v := VoltageLevelEquipment{
		Lines:               []models.ACLineSegment{},
		Generators:          []models.SynchronousMachine{},
		ConformLoads:        []models.ConformLoad{},
		LineTerminalNumbers: make(map[uuid.UUID]int),
	}

	for _, opt := range opts {
		opt(&v)
	}
	return &v
}

func WithLines(lines []models.ACLineSegment) func(v *VoltageLevelEquipment) {
	return func(v *VoltageLevelEquipment) {
		v.Lines = lines
	}
}

func WithGenerators(gens []models.SynchronousMachine) func(v *VoltageLevelEquipment) {
	return func(v *VoltageLevelEquipment) {
		v.Generators = gens
	}
}

func WithConformLoads(loads []models.ConformLoad) func(v *VoltageLevelEquipment) {
	return func(v *VoltageLevelEquipment) {
		v.ConformLoads = loads
	}
}

func CreateFullyConnectedVoltageLevel(equipment *VoltageLevelEquipment) *VoltageLevelModel {
	numLines := len(equipment.Lines)
	numGens := len(equipment.Generators)
	numLoads := len(equipment.ConformLoads)
	name := equipment.VoltageLevel.Name

	var reportingGroup models.ReportingGroup
	reportingGroup.Mrid = uuid.New()
	reportingGroup.Name = fmt.Sprintf("Reporting group %s", name)
	reportingGroup.ShortName = fmt.Sprintf("RG %s", name)
	reportingGroup.Description = fmt.Sprintf("Reporting group for voltage level %s", name)

	numConnectivityNodes := numLines + numGens + numLoads

	// ConnectivityNodes of lines goes first. Very important
	conNodes := make([]models.ConnectivityNode, 0, numConnectivityNodes)
	for i := range numConnectivityNodes {
		conNodes = append(conNodes, CreateConnectivityNode(fmt.Sprintf("%d %s", i, name)))
	}

	equipmentBusNameMarkers := make([]models.BusNameMarker, 0, numConnectivityNodes)
	equipmentTerminals := make([]models.Terminal, 0, numConnectivityNodes)
	for i := range numConnectivityNodes {
		var conductingEquipmentMrid uuid.UUID
		sequenceNumber := 1
		if i < numLines {
			conductingEquipmentMrid = equipment.Lines[i].Mrid
			num, ok := equipment.LineTerminalNumbers[conductingEquipmentMrid]
			if ok {
				sequenceNumber = num
			}
		} else if i < numLines+numGens {
			conductingEquipmentMrid = equipment.Generators[i-numLines].Mrid
		} else {
			conductingEquipmentMrid = equipment.ConformLoads[i-numLines-numGens].Mrid
		}
		equipmentBusNameMarkers = append(equipmentBusNameMarkers, CreateBusNameMarker(fmt.Sprintf("%d %s", i, name), reportingGroup.Mrid))
		equipmentTerminals = append(equipmentTerminals, CreateTerminal(conNodes[i], conductingEquipmentMrid, equipmentBusNameMarkers[i], sequenceNumber))
	}

	numSwitches := numLines*(numLines+1)/2 + numLines*(numGens+numLoads)
	switches := make([]models.Switch, 0, numSwitches)
	for i := range numSwitches {
		breaker := CreateSwitch(fmt.Sprintf("%s %d", name, i), &equipment.VoltageLevel)
		switches = append(switches, breaker)
	}

	switchTerminals := make([]models.Terminal, 0, 2*len(switches))
	switchTerminalBnms := make([]models.BusNameMarker, 0, len(switchTerminals))

	// Add a switch between any equipment and all lines (also line to line)
	switchNo := 0
	for lineNo := range numLines {
		for srcEquipment := lineNo; srcEquipment < numConnectivityNodes; srcEquipment++ {
			linkName := fmt.Sprintf("%d-%d", srcEquipment, lineNo)
			cn1 := conNodes[srcEquipment]
			if srcEquipment == lineNo {
				// Create a connectivity node disconnecting the line
				//          /
				// -------*      ----- substation
				//
				// * Is the new connectivity node
				cn1 = CreateConnectivityNode(linkName)
				conNodes = append(conNodes, cn1)
			}
			cn2 := conNodes[lineNo]
			AssertDifferent(cn1, cn2)

			switchMrid := switches[switchNo].Mrid
			switchNo++

			bnm1 := CreateBusNameMarker(linkName, reportingGroup.Mrid)
			terminal1 := CreateTerminal(cn1, switchMrid, bnm1, 1)
			bnm2 := CreateBusNameMarker(linkName, reportingGroup.Mrid)
			terminal2 := CreateTerminal(cn2, switchMrid, bnm2, 2)

			switchTerminalBnms = append(switchTerminalBnms, bnm1)
			switchTerminalBnms = append(switchTerminalBnms, bnm2)
			switchTerminals = append(switchTerminals, terminal1)
			switchTerminals = append(switchTerminals, terminal2)
		}
	}
	return &VoltageLevelModel{
		ConnectivityNodes: conNodes,
		ReportingGroup:    reportingGroup,
		BusNameMarkers:    append(equipmentBusNameMarkers, switchTerminalBnms...),
		Switches:          switches,
		Terminals:         append(equipmentTerminals, switchTerminals...),
	}

}

func CreateConnectivityNode(associatedResourceName string) models.ConnectivityNode {
	var cn models.ConnectivityNode
	cn.Mrid = uuid.New()
	cn.Name = fmt.Sprintf("Connectivity Node %s", associatedResourceName)
	cn.ShortName = fmt.Sprintf("CN %s", associatedResourceName)
	cn.Description = fmt.Sprintf("Connectivity node for %s", associatedResourceName)
	return cn
}

func CreateBusNameMarker(name string, repGroupMrid uuid.UUID) models.BusNameMarker {
	var bn models.BusNameMarker
	bn.Mrid = uuid.New()
	bn.Name = fmt.Sprintf("Bus Name Marker for %s", name)
	bn.ShortName = fmt.Sprintf("BMN %s", name)
	bn.Description = fmt.Sprintf("Bus Name Marker for %s", name)
	bn.ReportingGroupMrid = repGroupMrid
	return bn
}

func CreateTerminal(cn models.ConnectivityNode, conductingEquipmentMrid uuid.UUID, bnm models.BusNameMarker, seqNo int) models.Terminal {
	var terminal models.Terminal

	terminal.Mrid = uuid.New()
	terminal.Name = fmt.Sprintf("Terminal %s", cn.Name)
	terminal.ShortName = fmt.Sprintf("T %s", cn.Name)
	terminal.Description = fmt.Sprintf("Terminal for %s", cn.Name)
	terminal.SequenceNumber = seqNo
	terminal.BusNameMarkerMrid = bnm.Mrid
	terminal.ConnectivityNodeMrid = cn.Mrid
	terminal.ConductingEquipmentMrid = conductingEquipmentMrid
	terminal.PhasesId = 1
	return terminal
}

func CreateSwitch(name string, vl *models.VoltageLevel) models.Switch {
	var breaker models.Switch
	breaker.Mrid = uuid.New()
	breaker.Name = fmt.Sprintf("Switch %s", name)
	breaker.ShortName = fmt.Sprintf("Sw %s", name)
	breaker.Description = fmt.Sprintf("Switch %s", name)
	breaker.BaseVoltageMrid = vl.BaseVoltageMrid
	breaker.EquipmentContainerMrid = vl.Mrid
	return breaker
}

func LineMrids(lines []models.ACLineSegment) []uuid.UUID {
	result := make([]uuid.UUID, len(lines))
	for i, line := range lines {
		result[i] = line.Mrid
	}
	return result
}

func GetTargetTerminalSequenceNumber(ctx context.Context, db *bun.DB, lines []uuid.UUID) (map[uuid.UUID]int, error) {
	result := make(map[uuid.UUID]int)
	var terminals []models.Terminal
	err := db.NewSelect().Model(&terminals).Where("conducting_equipment_mrid IN (?)", bun.In(lines)).Scan(ctx)
	if err != nil {
		return result, fmt.Errorf("could not fetch existing terminals: %w", err)
	}
	terminals = OnlyActiveLatest(terminals)

	haveMultipleTerminals := []string{}
	for _, terminal := range terminals {
		_, exists := result[terminal.ConductingEquipmentMrid]
		if exists {
			haveMultipleTerminals = append(haveMultipleTerminals, terminal.ConductingEquipmentMrid.String())
			continue
		}

		targetSequenceNumber := terminal.SequenceNumber%2 + 1
		result[terminal.ConductingEquipmentMrid] = targetSequenceNumber
	}

	if len(haveMultipleTerminals) > 0 {
		return result, fmt.Errorf("The lines %s have already multiple terminals.", strings.Join(haveMultipleTerminals, ", "))
	}
	return result, nil
}
