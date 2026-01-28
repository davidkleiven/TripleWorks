package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type VoltageLevelModel struct {
	BusNameMarkers    []models.BusNameMarker
	ConnectivityNodes []models.ConnectivityNode
	ReportingGroup    models.ReportingGroup
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
				for i := range entities {
					entities[i].CommitId = int(commit.Id)
				}
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

// CreateFullyConnectedVoltageLevel connects all equipment to each other via a switch. If a connection
// already exists, a new one is not created.
// The state of connector is modified during the run to take into account emerging connections
func CreateFullyConnectedVoltageLevel(equipment *VoltageLevelEquipment, connector *EquipmentConnector) *VoltageLevelModel {
	name := equipment.VoltageLevel.Name
	var (
		reportingGroup           models.ReportingGroup
		numLinesAlreadyConnected int
		numGenAlreadyConnected   int
		numLoadAlreadyConnected  int
	)
	reportingGroup.Mrid = uuid.New()
	reportingGroup.Name = fmt.Sprintf("Reporting group %s", name)
	reportingGroup.ShortName = fmt.Sprintf("RG %s", name)
	reportingGroup.Description = fmt.Sprintf("Reporting group for voltage level %s", name)

	connections := []ConnectionResult{}

	// Connect line to line
	for i, line1 := range equipment.Lines {
		for _, line2 := range equipment.Lines[i+1:] {

			// Equipment within a voltage level is connected if there exists a path of equipment (E)
			// and connectivity nodes (CN) such that
			//
			// E1 ---- CN ---- Switch ---- CN ---- E2
			//
			// Which results in four edges (e.g. strictly smaller than 5)
			if connector.IsConnected(line1.Mrid, line2.Mrid, 5) {
				numLinesAlreadyConnected++
				continue
			}

			seq1, ok1 := equipment.LineTerminalNumbers[line1.Mrid]
			seq2, ok2 := equipment.LineTerminalNumbers[line2.Mrid]
			if ok1 {
				seq2 = seq1%2 + 1
			} else if ok2 {
				seq1 = seq2%2 + 1
			} else {
				seq1, seq2 = 1, 2
			}

			params := ConnectParams{
				Mrid1:              line1.Mrid,
				Mrid2:              line2.Mrid,
				CreateSeqNo1:       seq1,
				CreateSeqNo2:       seq2,
				ReportingGroupMrid: reportingGroup.Mrid,
				VoltageLevel:       equipment.VoltageLevel,
			}

			result := connector.Connect(&params)
			connections = append(connections, *result)

			// Update connector with new terminals
			connector.AddTerminals(result.Terminals...)
		}
	}

	// Connect generator to lines
	for _, gen := range equipment.Generators {
		for _, line := range equipment.Lines {
			if connector.IsConnected(gen.Mrid, line.Mrid, 5) {
				numGenAlreadyConnected++
				continue
			}
			params := ConnectParams{
				Mrid1:              gen.Mrid,
				Mrid2:              line.Mrid,
				CreateSeqNo1:       1,
				ReportingGroupMrid: reportingGroup.Mrid,
				VoltageLevel:       equipment.VoltageLevel,
			}

			result := connector.Connect(&params)
			connections = append(connections, *result)

			// Update connector with new terminals
			connector.AddTerminals(result.Terminals...)
		}
	}

	// Connect loads to lines
	for _, load := range equipment.ConformLoads {
		for _, line := range equipment.Lines {
			if connector.IsConnected(load.Mrid, line.Mrid, 5) {
				numLoadAlreadyConnected++
				continue
			}
			params := ConnectParams{
				Mrid1:              load.Mrid,
				Mrid2:              line.Mrid,
				CreateSeqNo1:       1,
				ReportingGroupMrid: reportingGroup.Mrid,
				VoltageLevel:       equipment.VoltageLevel,
			}

			result := connector.Connect(&params)
			connections = append(connections, *result)

			// Update connector with new terminals
			connector.AddTerminals(result.Terminals...)
		}
	}

	slog.Info("Fraction already connected",
		"lines", fmt.Sprintf("%d/%d", numLinesAlreadyConnected, len(equipment.Lines)),
		"generators", fmt.Sprintf("%d/%d", numGenAlreadyConnected, len(equipment.Generators)),
		"load", fmt.Sprintf("%d/%d", numLoadAlreadyConnected, len(equipment.ConformLoads)),
	)

	voltageLevelModel := VoltageLevelModel{
		BusNameMarkers:    []models.BusNameMarker{},
		ConnectivityNodes: []models.ConnectivityNode{},
		ReportingGroup:    reportingGroup,
		Switches:          []models.Switch{},
		Terminals:         []models.Terminal{},
	}

	for _, con := range connections {
		voltageLevelModel.BusNameMarkers = append(voltageLevelModel.BusNameMarkers, con.BusNameMarkers...)
		voltageLevelModel.ConnectivityNodes = append(voltageLevelModel.ConnectivityNodes, con.ConnectivityNodes...)
		voltageLevelModel.Switches = append(voltageLevelModel.Switches, con.Switch)
		voltageLevelModel.Terminals = append(voltageLevelModel.Terminals, con.Terminals...)
	}
	return &voltageLevelModel
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

func CreateTerminal(cnMrid uuid.UUID, conductingEquipmentMrid uuid.UUID, bnm models.BusNameMarker, seqNo int) models.Terminal {
	var terminal models.Terminal

	terminal.Mrid = uuid.New()
	terminal.Name = fmt.Sprintf("Terminal %s", bnm.Name)
	terminal.ShortName = fmt.Sprintf("T %s", bnm.Name)
	terminal.Description = fmt.Sprintf("Terminal for %s", bnm.Name)
	terminal.SequenceNumber = seqNo
	terminal.BusNameMarkerMrid = bnm.Mrid
	terminal.ConnectivityNodeMrid = cnMrid
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
