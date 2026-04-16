package pkg

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"com.github/davidkleiven/tripleworks/models"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/google/uuid"
)

func IdentifiedLine(mrid uuid.UUID, name string) *models.ACLineSegment {
	var line models.ACLineSegment
	line.Mrid = mrid
	line.Name = name
	return &line
}

func IdentifiedSubstation(mrid uuid.UUID, name string, opts ...func(s *models.Substation)) *models.Substation {
	var substation models.Substation
	substation.Mrid = mrid
	substation.Name = name
	for _, opt := range opts {
		opt(&substation)
	}
	return &substation
}

func WithSubstation(sub uuid.UUID) func(v *models.VoltageLevel) {
	return func(v *models.VoltageLevel) {
		v.SubstationMrid = sub
	}
}

func IdentifiedVoltageLevel(mrid uuid.UUID, name string, opts ...func(v *models.VoltageLevel)) *models.VoltageLevel {
	var vl models.VoltageLevel
	vl.Mrid = mrid
	vl.Name = name
	for _, opt := range opts {
		opt(&vl)
	}
	return &vl
}

func IdentifiedTerminal(mrid uuid.UUID, name string, seqNo int, opts ...func(t *models.Terminal)) *models.Terminal {
	var t models.Terminal
	t.Mrid = mrid
	t.Name = name
	t.SequenceNumber = seqNo
	for _, opt := range opts {
		opt(&t)
	}
	return &t
}

func WithConductingEquipment(mrid uuid.UUID) func(t *models.Terminal) {
	return func(t *models.Terminal) {
		t.ConductingEquipmentMrid = mrid
	}
}

func WithConnectivityNode(mrid uuid.UUID) func(t *models.Terminal) {
	return func(t *models.Terminal) {
		t.ConnectivityNodeMrid = mrid
	}
}

func WithBaseVoltage(bvMrid uuid.UUID) func(v *models.VoltageLevel) {
	return func(v *models.VoltageLevel) {
		v.BaseVoltageMrid = bvMrid
	}
}

func WithLocation(locMrid uuid.UUID) func(s *models.Substation) {
	return func(s *models.Substation) {
		s.LocationMrid = locMrid
	}
}

func IdentifiedBaseVoltage(mrid uuid.UUID, name string, nominalVoltage float64) *models.BaseVoltage {
	var bv models.BaseVoltage
	bv.Mrid = mrid
	bv.Name = name
	bv.NominalVoltage = nominalVoltage
	return &bv
}

func IdentifiedConnectivityNode(mrid uuid.UUID, name string, containerMrid uuid.UUID) *models.ConnectivityNode {
	var cn models.ConnectivityNode
	cn.Mrid = mrid
	cn.Name = name
	cn.ConnectivityNodeContainerMrid = containerMrid
	return &cn
}

func IdentifiedLocation(mrid uuid.UUID, name string, coordSystemMrid uuid.UUID) *models.Location {
	var loc models.Location
	loc.Mrid = mrid
	loc.Name = name
	loc.CoordinateSystemMrid = coordSystemMrid
	return &loc
}

func IdentifiedPositionPoint(x, y, z float64, seqNo int, locationMrid uuid.UUID) *models.PositionPoint {
	var pp models.PositionPoint
	pp.XPosition = x
	pp.YPosition = y
	pp.ZPosition = z
	pp.SequenceNumber = seqNo
	pp.LocationMrid = locationMrid
	return &pp
}

type E2EData struct {
	Model  models.Model
	Commit models.Commit
	Data   []any
}

func MakeE2eData() *E2EData {
	model := models.Model{Name: "e2e-model", Id: 1}

	subCMrid := uuid.MustParse("06a747c0-c16b-4a16-8e67-d9b6177d76cb")
	subDMrid := uuid.MustParse("d657c174-0e1e-47cf-bf19-29fb85955295")
	vlcMrid := uuid.MustParse("ab58a50d-c587-4004-813a-06d52008b6cb")
	vldMrid := uuid.MustParse("d18749a1-1043-4606-9435-3f1bbec690ba")

	lineMrid := uuid.MustParse("4e832836-ef53-458e-9711-903982551fcf")
	bvCMrid := uuid.MustParse("5c9ae4b1-3c5d-4b7a-8f2e-1d6c3a9b4e8f")
	bvDMrid := uuid.MustParse("6d7bf5c2-4d6e-5c8b-9a3f-2e7d4b0c5f9a")
	locCMrid := uuid.MustParse("7e8cf6d3-5e7f-6d9c-1b4a-3f8e5c1d6a0b")
	locDMrid := uuid.MustParse("8f9da7e4-6f8a-7e1d-2c5b-4a9f6d2e7b1c")
	coordMrid := uuid.MustParse("9a1eb8f5-7b9c-8f2e-3d6c-5b0c7e3f8c2d")
	cnCMrid := uuid.MustParse("1b2fc9c6-8c1d-9c3f-4e7d-6c1d8f4c9d3e")
	cnDMrid := uuid.MustParse("2c3ad0d7-9d2e-1d4c-5f8e-7d2e9c5d0e4f")
	term1Mrid := uuid.MustParse("3d4be1e8-1e4f-2e5d-6f9f-8e4f0d6e1f5f")
	term2Mrid := uuid.MustParse("4e5cf2f9-2f5a-3f6e-7f0a-8f5a1e7f2a6f")

	concreteKinds := []models.MridGetter{
		// e2e test where Sub A is connected to Sub B
		IdentifiedLine(uuid.MustParse("ce8e57c7-8f6c-42c3-8b8e-e06aa39f0da3"), "Unconnected line"),
		IdentifiedSubstation(uuid.MustParse("fed4f58f-199c-43c7-95f1-b353f55ae12c"), "Substation A"),
		IdentifiedSubstation(uuid.MustParse("8fbd0382-e14c-491b-b4d1-7b2b13be27fb"), "Substation B"),

		// Data for map test
		IdentifiedBaseVoltage(bvCMrid, "BaseVoltage 138kV", 138000),
		IdentifiedBaseVoltage(bvDMrid, "BaseVoltage 138kV", 138000),
		IdentifiedLocation(locCMrid, "Location Sub C", coordMrid),
		IdentifiedLocation(locDMrid, "Location Sub D", coordMrid),
		IdentifiedSubstation(subCMrid, "Substation C", WithLocation(locCMrid)),
		IdentifiedSubstation(subDMrid, "Substation D", WithLocation(locDMrid)),
		IdentifiedVoltageLevel(vlcMrid, "Vl sub C", WithSubstation(subCMrid), WithBaseVoltage(bvCMrid)),
		IdentifiedVoltageLevel(vldMrid, "Vl sub D", WithSubstation(subDMrid), WithBaseVoltage(bvDMrid)),
		IdentifiedConnectivityNode(cnCMrid, "Node C", vlcMrid),
		IdentifiedConnectivityNode(cnDMrid, "Node D", vldMrid),
		IdentifiedLine(lineMrid, "Connected line"),
		IdentifiedTerminal(term1Mrid, "Terminal 1", 1, WithConductingEquipment(lineMrid), WithConnectivityNode(cnCMrid)),
		IdentifiedTerminal(term2Mrid, "Terminal 2", 2, WithConductingEquipment(lineMrid), WithConnectivityNode(cnDMrid)),
		IdentifiedLocation(coordMrid, "WGS84", uuid.Nil),
	}

	positionPoints := []any{
		IdentifiedPositionPoint(10.3951, 63.4305, 0, 1, locCMrid),
		IdentifiedPositionPoint(10.7509, 59.9139, 0, 2, locDMrid),
	}

	var entities []*models.Entity
	for _, item := range concreteKinds {
		entities = append(entities, &models.Entity{
			Mrid:        item.GetMrid(),
			ModelEntity: models.ModelEntity{ModelId: 1},
			EntityType:  StructName(item),
		})
	}

	commit := models.Commit{
		Id:        1,
		Message:   "Populate database with e2e data",
		Author:    "TripleWorks",
		CreatedAt: time.Now(),
	}

	var data []any
	for _, e := range entities {
		data = append(data, e)
	}

	for _, c := range concreteKinds {
		data = append(data, c)
	}

	for _, pp := range positionPoints {
		data = append(data, pp)
	}

	return &E2EData{
		Model:  model,
		Commit: commit,
		Data:   data,
	}
}

func InsertE2eData(data *E2EData, inserter repository.Inserter) {
	slog.Info("Inserting E2E data")
	err1 := inserter.Insert(context.Background(), &data.Model)

	var numInserted int

	onInsert := func(item any) error {
		numInserted++
		return nil
	}

	err := InsertAllInserter(context.Background(), inserter, data.Commit, slices.Values(data.Data), onInsert)
	slog.Info("Inserted e2e data", "numRecords", numInserted, "errModelInsert", err1, "errEntityInsert", err)
}
