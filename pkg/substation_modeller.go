package pkg

import (
	"context"
	"fmt"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ConnectableVoltageLevel struct {
	ConnectivityNodes []models.ConnectivityNode
	VoltageLevel      models.VoltageLevel
	BaseVoltage       models.BaseVoltage
}

func (c *ConnectableVoltageLevel) RequireConsistentVoltageMrid() {
	if c.BaseVoltage.Mrid != c.VoltageLevel.BaseVoltageMrid {
		panic(fmt.Sprintf("Inconsistent base voltage mrid: %s (base voltage) and %s (voltage level)", c.BaseVoltage.Mrid, c.VoltageLevel.BaseVoltageMrid))
	}
}

type SubstationModel struct {
	ReportingGroup    models.ReportingGroup
	BusNameMarkers    []models.BusNameMarker
	ConnectivityNodes []models.ConnectivityNode
	Switches          []models.Switch
	Terminals         []models.Terminal
	TransformerEnds   []models.PowerTransformerEnd
	Transformers      []models.PowerTransformer
}

func (s *SubstationModel) AssignCommitId(commitId int) {
	s.ReportingGroup.CommitId = commitId
	for i := range s.BusNameMarkers {
		s.BusNameMarkers[i].CommitId = commitId
	}

	for i := range s.ConnectivityNodes {
		s.ConnectivityNodes[i].CommitId = commitId
	}

	for i := range s.Switches {
		s.Switches[i].CommitId = commitId
	}

	for i := range s.Terminals {
		s.Terminals[i].CommitId = commitId
	}
	for i := range s.Transformers {
		s.Transformers[i].CommitId = commitId
	}
	for i := range s.TransformerEnds {
		s.TransformerEnds[i].CommitId = commitId
	}
}

func (s *SubstationModel) Entities(modelId int) []models.Entity {
	entities := make([]models.Entity, 0, len(s.BusNameMarkers)+len(s.ConnectivityNodes)+len(s.Switches)+len(s.Terminals)+len(s.Transformers)+len(s.TransformerEnds)+1)

	entities = append(entities, models.Entity{Mrid: s.ReportingGroup.Mrid, ModelEntity: models.ModelEntity{ModelId: modelId}})

	for _, bnm := range s.BusNameMarkers {
		entity := models.Entity{Mrid: bnm.Mrid, ModelEntity: models.ModelEntity{ModelId: modelId}}
		entities = append(entities, entity)
	}

	for _, conNode := range s.ConnectivityNodes {
		entity := models.Entity{Mrid: conNode.Mrid, ModelEntity: models.ModelEntity{ModelId: modelId}}
		entities = append(entities, entity)
	}

	for _, breaker := range s.Switches {
		entity := models.Entity{Mrid: breaker.Mrid, ModelEntity: models.ModelEntity{ModelId: modelId}}
		entities = append(entities, entity)
	}

	for _, terminal := range s.Terminals {
		entity := models.Entity{Mrid: terminal.Mrid, ModelEntity: models.ModelEntity{ModelId: modelId}}
		entities = append(entities, entity)
	}

	for _, transformers := range s.Transformers {
		entity := models.Entity{Mrid: transformers.Mrid, ModelEntity: models.ModelEntity{ModelId: modelId}}
		entities = append(entities, entity)
	}

	for _, windings := range s.TransformerEnds {
		entity := models.Entity{Mrid: windings.Mrid, ModelEntity: models.ModelEntity{ModelId: modelId}}
		entities = append(entities, entity)
	}
	return entities
}

func (s *SubstationModel) Write(ctx context.Context, db *bun.DB, modelId int, msg string) error {
	entities := s.Entities(modelId)
	commit := models.Commit{
		Message: msg,
		Author:  "SubstationModeller",
	}
	return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := ReturnOnFirstError(
			func() error {
				_, err := tx.NewInsert().Model(&commit).Exec(ctx)
				return err
			},
			func() error {
				s.AssignCommitId(int(commit.Id))
				return nil
			},
			func() error {
				_, err := tx.NewInsert().Model(&entities).Exec(ctx)
				return err
			},
			func() error {
				_, err := tx.NewInsert().Model(&s.BusNameMarkers).Exec(ctx)
				return err
			},
			func() error {
				_, err := tx.NewInsert().Model(&s.ConnectivityNodes).Exec(ctx)
				return err
			},
			func() error {
				_, err := tx.NewInsert().Model(&s.Switches).Exec(ctx)
				return err
			},
			func() error {
				_, err := tx.NewInsert().Model(&s.Transformers).Exec(ctx)
				return err
			},
			func() error {
				_, err := tx.NewInsert().Model(&s.TransformerEnds).Exec(ctx)
				return err
			},
			func() error {
				_, err := tx.NewInsert().Model(&s.ReportingGroup).Exec(ctx)
				return err
			},
			func() error {
				_, err := tx.NewInsert().Model(&s.Terminals).Exec(ctx)
				return err
			},
		)
		return err
	})
}

type SubstationData struct {
	Substation    models.Substation
	VoltageLevels []ConnectableVoltageLevel
}

func CreateFullyConnectedSubstation(data SubstationData, connector *EquipmentConnector) *SubstationModel {
	for _, vl := range data.VoltageLevels {
		vl.RequireConsistentVoltageMrid()
	}

	var (
		repGroup            models.ReportingGroup
		numAlreadyConnected int
	)
	repGroup.Mrid = uuid.New()
	repGroup.Name = fmt.Sprintf("%s voltage level connections", data.Substation.Name)
	repGroup.ShortName = fmt.Sprintf("%s vlc", data.Substation.Name)
	repGroup.Description = "Group representing all items created to connect voltage levels"

	bnms := []models.BusNameMarker{}
	conNodes := []models.ConnectivityNode{}
	switches := []models.Switch{}
	terminals := []models.Terminal{}
	transformers := []models.PowerTransformer{}
	windings := []models.PowerTransformerEnd{}

	vls := data.VoltageLevels
	for i, vl1 := range vls {
		for _, vl2 := range vls[i+1:] {
			for ii, cn1 := range vl1.ConnectivityNodes {
				for jj, cn2 := range vl2.ConnectivityNodes {
					// Connectivity nodes on different voltage levels has a direct connection
					// if they are related by
					//
					// CN1 --- Switch --- cn --- Transformer --- CN2
					//
					// which results in four edges
					if connector.IsConnected(cn1.Mrid, cn2.Mrid, 5) {
						numAlreadyConnected++
						continue
					}

					name := fmt.Sprintf("%.0f kV to %.0f kV (CN %d-%d)", vl1.BaseVoltage.NominalVoltage, vl2.BaseVoltage.NominalVoltage, ii, jj)
					transformer := CreateTransformer(name, data.Substation.Mrid, vl1.VoltageLevel.BaseVoltageMrid)

					breaker := CreateSwitch(name, &vl1.VoltageLevel)
					cnSwitch := CreateConnectivityNode(name)

					// Terminal 1 connects to new node, which also winding terminal is connected to
					switchBnm1 := CreateBusNameMarker(name, repGroup.Mrid)
					switchT1 := CreateTerminal(cnSwitch.Mrid, breaker.Mrid, switchBnm1, 1)

					// Switch two connects to connectivity node at v1
					switchBnm2 := CreateBusNameMarker(name, repGroup.Mrid)
					switchT2 := CreateTerminal(cn1.Mrid, breaker.Mrid, switchBnm2, 2)

					// Transformer terminals
					bnmT1 := CreateBusNameMarker(name, repGroup.Mrid)
					terminal1 := CreateTerminal(cnSwitch.Mrid, transformer.Mrid, bnmT1, 1)

					windingParams1 := WindingParams{
						Name:                 name,
						EndNumber:            terminal1.SequenceNumber,
						TerminalMrid:         terminal1.Mrid,
						BaseVoltageMrid:      vl1.VoltageLevel.BaseVoltageMrid,
						PowerTransformerMrid: transformer.Mrid,
						ConnectionKind:       4, // Y
						RatedU:               vl1.BaseVoltage.NominalVoltage,
					}
					winding1 := CreateWinding(windingParams1)

					bnmT2 := CreateBusNameMarker(name, repGroup.Mrid)
					terminal2 := CreateTerminal(cn2.Mrid, transformer.Mrid, bnmT2, 2)

					windingParams2 := WindingParams{
						Name:                 name,
						EndNumber:            terminal2.SequenceNumber,
						TerminalMrid:         terminal2.Mrid,
						BaseVoltageMrid:      vl2.VoltageLevel.BaseVoltageMrid,
						PowerTransformerMrid: transformer.Mrid,
						ConnectionKind:       5, // Yn
						RatedU:               vl2.BaseVoltage.NominalVoltage,
					}
					winding2 := CreateWinding(windingParams2)

					conNodes = append(conNodes, cnSwitch)
					switches = append(switches, breaker)
					terminals = append(terminals, switchT1, switchT2, terminal1, terminal2)
					windings = append(windings, winding1, winding2)
					bnms = append(bnms, switchBnm1, switchBnm2, bnmT1, bnmT2)
					transformers = append(transformers, transformer)

					// Update connectors terminals
					connector.AddTerminals(switchT1, switchT2, terminal1, terminal2)
				}
			}
		}
	}
	return &SubstationModel{
		ReportingGroup:    repGroup,
		BusNameMarkers:    bnms,
		ConnectivityNodes: conNodes,
		Switches:          switches,
		Terminals:         terminals,
		TransformerEnds:   windings,
		Transformers:      transformers,
	}
}

func CreateTransformer(name string, substationMrid uuid.UUID, bvMrid uuid.UUID) models.PowerTransformer {
	var transformer models.PowerTransformer
	transformer.Mrid = uuid.New()
	transformer.BaseVoltageMrid = bvMrid // Weird that a transformer has a base voltage
	transformer.Name = fmt.Sprintf("Transformer %s", name)
	transformer.ShortName = fmt.Sprintf("Trns %s", name)
	transformer.Description = fmt.Sprintf("Transformer %s", name)
	transformer.EquipmentContainerMrid = substationMrid
	return transformer
}

type WindingParams struct {
	Name                 string
	EndNumber            int
	TerminalMrid         uuid.UUID
	BaseVoltageMrid      uuid.UUID
	PowerTransformerMrid uuid.UUID
	ConnectionKind       int
	RatedU               float64
}

func CreateWinding(params WindingParams) models.PowerTransformerEnd {
	var winding models.PowerTransformerEnd
	winding.Mrid = uuid.New()
	winding.Name = fmt.Sprintf("Winding %s", params.Name)
	winding.ShortName = fmt.Sprintf("Wnd %s", params.Name)
	winding.Description = winding.Name
	winding.EndNumber = params.EndNumber
	winding.TerminalMrid = params.TerminalMrid
	winding.BaseVoltageMrid = params.BaseVoltageMrid
	winding.PowerTransformerMrid = params.PowerTransformerMrid
	winding.ConnectionKindId = params.ConnectionKind
	winding.RatedU = params.RatedU
	winding.RatedS = 100.0

	z_pu := winding.RatedU * winding.RatedU / winding.RatedS
	winding.R = 0.002 * z_pu
	winding.X = 0.1 * z_pu
	winding.G = 0.0
	winding.B = 0.005 / z_pu
	return winding
}
