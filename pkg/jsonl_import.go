package pkg

import (
	"fmt"
	"iter"
	"math"
	"regexp"
	"strconv"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
)

func mridFromName(prefix string, name string) uuid.UUID {
	namespace := uuid.NewSHA1(uuid.NameSpaceURL, []byte(prefix))
	return uuid.NewSHA1(namespace, []byte(name))
}

func locationMrid(name string) uuid.UUID {
	return mridFromName("Location", name)
}

func subGeoRegionMrid(name string) uuid.UUID {
	return mridFromName("SubGeoGraphicalRegion", name)
}

func CoordinateSystemMrid() uuid.UUID {
	return mridFromName("CoordinateSystem", "EPSG: 4326")
}

func substationMrid(name string) uuid.UUID {
	return mridFromName("Substation", name)
}

func baseVoltageMrid(voltage int) uuid.UUID {
	return mridFromName("BaseVoltage", strconv.Itoa(voltage))
}

func lineMrid(name string) uuid.UUID {
	return mridFromName("ACLineSegment", name)
}

func lineContainerMrid(name string) uuid.UUID {
	return mridFromName("Line", name)
}

func syncMachineMrid(name string) uuid.UUID {
	return mridFromName("SynchronousMachine", name)
}

func voltageLevelMrid(substation string, voltage int) uuid.UUID {
	return mridFromName("VoltageLevel", fmt.Sprintf("%s %d", substation, voltage))
}

func regulatingControlMrid(name string) uuid.UUID {
	return mridFromName("RegulatingControl", name)
}

func terminalMrid(name string) uuid.UUID {
	return mridFromName("Terminal", name)
}

func busNameMarkerMrid(name string) uuid.UUID {
	return mridFromName("Bus name marker", name)
}

func reportingGroupMrid(name string) uuid.UUID {
	return mridFromName("Reporting group", name)
}

func conNodeMrid(name string) uuid.UUID {
	return mridFromName("Connectivity node", name)
}

func reactCurveMrid(pf float64) uuid.UUID {
	return mridFromName("ReactiveCurve", fmt.Sprintf("%.2f", pf))
}

func genUnitMrid(name string) uuid.UUID {
	return mridFromName("GeneratingUnit", name)
}

func hydroPowerPlantMrid(name string) uuid.UUID {
	return mridFromName("HydroPowerPlant", name)
}

func geoRegionMrid(name string) uuid.UUID {
	return mridFromName("GeographicalRegion", name)
}

func loadMrid(name string) uuid.UUID {
	return mridFromName("ConformLoad", name)
}

func constantLoadMrid(p float64) uuid.UUID {
	return mridFromName("ConstantLoad", fmt.Sprintf("%.2f", p))
}

func loadGroupMrid(name string) uuid.UUID {
	return mridFromName("LoadGroup", name)
}

func loadAreaMrid(name string) uuid.UUID {
	return mridFromName("LoadArea", name)
}

func subLoadAreaMrid(name string) uuid.UUID {
	return mridFromName("SubLoadArea", name)
}

type SubstationLight struct {
	Name   string  `json:"name"`
	Region string  `json:"region"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

func (s *SubstationLight) CimItems(modelId int) iter.Seq[any] {
	var coordinateSystem models.CoordinateSystem
	coordinateSystem.Mrid = CoordinateSystemMrid()
	coordinateSystem.Name = "Longitude, latitude"
	coordinateSystem.CrsUrn = "EPSG:4326"

	var location models.Location
	location.Mrid = locationMrid(s.Name)
	location.Name = "Location " + s.Name
	location.ShortName = "Loc " + s.Name
	location.Description = "Geographical location of " + s.Name
	location.CoordinateSystemMrid = coordinateSystem.Mrid

	var point models.PositionPoint
	point.XPosition = s.X
	point.YPosition = s.Y
	point.SequenceNumber = 1
	point.LocationMrid = location.Mrid

	re := regexp.MustCompile("[0-9]")
	geoName := re.ReplaceAllString(s.Region, "")

	var geoRegion models.GeographicalRegion
	geoRegion.Mrid = geoRegionMrid(geoName)
	geoRegion.Name = geoName
	geoRegion.ShortName = geoName

	var subGeoRegion models.SubGeographicalRegion
	subGeoRegion.Mrid = subGeoRegionMrid(s.Region)
	subGeoRegion.Name = s.Region
	subGeoRegion.ShortName = s.Region
	subGeoRegion.GeographicalRegionMrid = geoRegion.Mrid

	var substation models.Substation
	substation.Mrid = substationMrid(s.Name)
	substation.Name = s.Name
	substation.ShortName = s.Name
	substation.Description = s.Name
	substation.LocationMrid = location.Mrid
	substation.SubGeographicalRegionMrid = subGeoRegion.Mrid

	model := models.ModelEntity{ModelId: modelId}
	entities := []models.Entity{
		{ModelEntity: model, Mrid: location.Mrid, EntityType: StructName(location)},
		{ModelEntity: model, Mrid: substation.Mrid, EntityType: StructName(substation)},
		{ModelEntity: model, Mrid: coordinateSystem.Mrid, EntityType: StructName(coordinateSystem)},
		{ModelEntity: model, Mrid: geoRegion.Mrid, EntityType: StructName(geoRegion)},
		{ModelEntity: model, Mrid: subGeoRegion.Mrid, EntityType: StructName(subGeoRegion)},
	}
	return func(yield func(v any) bool) {
		for _, entity := range entities {
			if !yield(&entity) {
				return
			}
		}
		yieldMany(yield, &location, &point, &substation, &coordinateSystem, &geoRegion, &subGeoRegion)
	}
}

func CreateBaseVoltage(voltage int) models.BaseVoltage {
	var bv models.BaseVoltage
	bv.Mrid = baseVoltageMrid(voltage)
	bv.Name = fmt.Sprintf("Base voltage %d kV", voltage)
	bv.ShortName = fmt.Sprintf("%d kV", voltage)
	bv.NominalVoltage = float64(voltage)
	return bv
}

type LineLight struct {
	FromSubstation string  `json:"from"`
	ToSubstation   string  `json:"to"`
	Length         float64 `json:"length"`
	Voltage        int     `json:"voltage"`
}

func (l *LineLight) CimItems(modelId int) iter.Seq[any] {
	name := fmt.Sprintf("%s-%s (%d kV)", l.FromSubstation, l.ToSubstation, l.Voltage)
	formState := FormState{
		Length: l.Length,
		Name:   name,
	}

	bv := CreateBaseVoltage(l.Voltage)

	var geo models.GeographicalRegion
	geo.Mrid = geoRegionMrid("Transmission lines")
	geo.Name = "Transmission lines"
	geo.Description = "All transmission lines belongs to this region"

	var subRegion models.SubGeographicalRegion
	subRegion.Name = fmt.Sprintf("Transmission lines %d kV", l.Voltage)
	subRegion.Mrid = subGeoRegionMrid(subRegion.Name)
	subRegion.Description = fmt.Sprintf("All transmission lines at %d kV", l.Voltage)
	subRegion.GeographicalRegionMrid = geo.Mrid

	var lineContainer models.Line
	lineContainer.Mrid = lineContainerMrid(name)
	lineContainer.Name = "Line: " + name
	lineContainer.ShortName = "Line: " + name
	lineContainer.RegionMrid = subRegion.Mrid

	var line models.ACLineSegment
	line.Mrid = lineMrid(name)
	line.Name = name
	line.ShortName = MustGet(stringAutofillers, "ShortName")(&formState)
	line.Description = "AC line " + name
	line.Length = l.Length
	line.R = MustGet(floatAutofillers, "R")(&formState)
	line.X = MustGet(floatAutofillers, "X")(&formState)
	line.Bch = MustGet(floatAutofillers, "Bch")(&formState)
	line.Gch = MustGet(floatAutofillers, "Gch")(&formState)
	line.BaseVoltageMrid = bv.Mrid
	line.EquipmentContainerMrid = lineContainer.Mrid

	model := models.ModelEntity{ModelId: modelId}
	entities := []models.Entity{
		{ModelEntity: model, Mrid: lineContainer.Mrid, EntityType: StructName(lineContainer)},
		{ModelEntity: model, Mrid: line.Mrid, EntityType: StructName(line)},
		{ModelEntity: model, Mrid: bv.Mrid, EntityType: StructName(bv)},
		{ModelEntity: model, Mrid: geo.Mrid, EntityType: StructName(geo)},
		{ModelEntity: model, Mrid: subRegion.Mrid, EntityType: StructName(subRegion)},
	}

	return func(yield func(v any) bool) {
		for _, entity := range entities {
			if !yield(&entity) {
				return
			}
		}
		yieldMany(yield, &lineContainer, &line, &bv, &geo, &subRegion)
	}
}

func CreateVoltageLevel(substation string, voltage int) models.VoltageLevel {
	var vl models.VoltageLevel
	vl.Mrid = voltageLevelMrid(substation, voltage)
	vl.Name = fmt.Sprintf("%s %d kV", substation, voltage)
	vl.ShortName = vl.Name
	vl.SubstationMrid = substationMrid(substation) // Must exist
	vl.BaseVoltageMrid = baseVoltageMrid(voltage)  // Must exist
	vl.LowVoltageLimit = 0.9 * float64(voltage)
	vl.HighVoltageLimit = 1.1 * float64(voltage)
	return vl
}

type GeneratorLight struct {
	Kind       string  `json:"kind"`
	Substation string  `json:"substation"`
	Num        int     `json:"num"`
	MaxP       float64 `json:"maxP"`
	MinP       float64 `json:"minP"`
	Voltage    int     `json:"voltage"`
}

func (g *GeneratorLight) CimItems(modelId int) iter.Seq[any] {
	name := fmt.Sprintf("%s G%d", g.Substation, g.Num)
	var repGroup models.ReportingGroup
	repGroup.Mrid = reportingGroupMrid(name)
	repGroup.Name = "RG " + name
	repGroup.ShortName = "RG " + name
	repGroup.Description = "Reporting group for " + name

	var bnm models.BusNameMarker
	bnm.Mrid = busNameMarkerMrid(name)
	bnm.Name = "Bus Name Marker " + name
	bnm.ShortName = bnm.Name
	bnm.Description = "Bus name marker for the terminal of " + name
	bnm.ReportingGroupMrid = repGroup.Mrid

	vl := CreateVoltageLevel(g.Substation, g.Voltage)

	var conNode models.ConnectivityNode
	conNode.Mrid = conNodeMrid(name)
	conNode.Name = "Connectivity node " + name
	conNode.ShortName = "CN " + name
	conNode.Description = "Connectivity node for generator " + name
	conNode.ConnectivityNodeContainerMrid = vl.Mrid

	var terminal models.Terminal
	terminal.Mrid = terminalMrid(name)
	terminal.Name = name
	terminal.ShortName = name
	terminal.Description = fmt.Sprintf("Terminal for generator %s and its regulating control", name)
	terminal.BusNameMarkerMrid = bnm.Mrid
	terminal.SequenceNumber = 1
	terminal.PhasesId = 1
	terminal.ConnectivityNodeMrid = conNode.Mrid

	var regControl models.RegulatingControl
	regControl.Mrid = regulatingControlMrid(name)
	regControl.Name = "Reg: " + name
	regControl.ShortName = regControl.Name
	regControl.Description = "Regulating control for " + name
	regControl.TerminalMrid = terminal.Mrid
	regControl.ModeId = 1

	var reactCurve models.ReactiveCapabilityCurve
	reactCurve.Mrid = reactCurveMrid(0.95)
	reactCurve.Name = "PF 0.95 linear curve"
	reactCurve.XUnitId = 10     // W
	reactCurve.Y1UnitId = 11    // VAr
	reactCurve.Y2UnitId = 11    // VAr
	reactCurve.CurveStyleId = 2 // Linear

	qMult := math.Tan(math.Acos(0.95))
	pts := []models.CurveData{
		{Xvalue: 0.0, Y1value: 0.0, Y2value: 0.0, CurveMrid: reactCurve.Mrid},
		{Xvalue: 1e9, Y1value: -qMult * 1e9, Y2value: qMult * 1e9, CurveMrid: reactCurve.Mrid},
	}

	bv := CreateBaseVoltage(g.Voltage)

	var machine models.SynchronousMachine
	machine.Mrid = syncMachineMrid(name)
	machine.Name = name
	machine.ShortName = name
	machine.Description = name
	machine.BaseVoltageMrid = bv.Mrid
	machine.RatedU = float64(g.Voltage)
	machine.RatedS = g.MaxP
	machine.EquipmentContainerMrid = vl.Mrid
	machine.RatedPowerFactor = 0.95
	machine.RegulatingControlMrid = regControl.Mrid
	terminal.ConductingEquipmentMrid = machine.Mrid
	machine.TypeId = 2

	machine.MinQ = -qMult * g.MaxP
	machine.MaxQ = qMult * g.MaxP
	machine.InitialReactiveCapabilityCurveMrid = reactCurve.Mrid

	var genUnit models.GeneratingUnit
	genUnit.Mrid = genUnitMrid(name)
	genUnit.MaxOperatingP = g.MaxP
	genUnit.MinOperatingP = g.MinP
	genUnit.GenControlSourceId = 1
	genUnit.NominalP = g.MaxP
	genUnit.RatedGrossMaxP = g.MaxP
	genUnit.RatedGrossMinP = g.MinP
	genUnit.LongPF = machine.RatedPowerFactor
	genUnit.ShortPF = machine.RatedPowerFactor
	genUnit.TotalEfficiency = 1.0
	genUnit.EquipmentContainerMrid = machine.EquipmentContainerMrid
	machine.GeneratingUnitMrid = genUnit.Mrid

	var plant models.HydroPowerPlant
	plant.HydroPlantStorageTypeId = 1
	hyd := models.HydroGeneratingUnit{GeneratingUnit: genUnit}
	wind := models.WindGeneratingUnit{GeneratingUnit: genUnit}
	thermal := models.ThermalGeneratingUnit{GeneratingUnit: genUnit}

	switch g.Kind {
	case "hydro":
		plant.Mrid = hydroPowerPlantMrid(g.Substation)
		plant.Name = "Plant " + g.Substation
		hyd.HydroPowerPlantMrid = plant.Mrid
		hyd.EnergyConversionCapabilityId = 1
	case "wind":
		wind.WindGenUnitTypeId = 1
	}

	model := models.ModelEntity{ModelId: modelId}
	entities := []models.Entity{
		{ModelEntity: model, Mrid: repGroup.Mrid, EntityType: StructName(repGroup)},
		{ModelEntity: model, Mrid: bnm.Mrid, EntityType: StructName(bnm)},
		{ModelEntity: model, Mrid: conNode.Mrid, EntityType: StructName(conNode)},
		{ModelEntity: model, Mrid: terminal.Mrid, EntityType: StructName(terminal)},
		{ModelEntity: model, Mrid: regControl.Mrid, EntityType: StructName(regControl)},
		{ModelEntity: model, Mrid: reactCurve.Mrid, EntityType: StructName(reactCurve)},
		{ModelEntity: model, Mrid: machine.Mrid, EntityType: StructName(machine)},
		{ModelEntity: model, Mrid: genUnit.Mrid, EntityType: StructName(genUnit)},
		{ModelEntity: model, Mrid: vl.Mrid, EntityType: StructName(vl)},
		{ModelEntity: model, Mrid: bv.Mrid, EntityType: StructName(bv)},
	}

	if g.Kind == "hydro" {
		entities = append(entities, models.Entity{ModelEntity: model, Mrid: plant.Mrid, EntityType: StructName(plant)})
	}

	return func(yield func(v any) bool) {
		for _, entity := range entities {
			if !yield(&entity) {
				return
			}
		}

		for _, pt := range pts {
			if !yield(&pt) {
				return
			}
		}

		yieldMany(yield, &repGroup, &bnm, &conNode, &terminal, &regControl, &reactCurve, &machine, &bv, &vl)
		if g.Kind == "hydro" && !yield(&plant) {
			return
		}

		switch g.Kind {
		case "hydro":
			yield(&hyd)
		case "wind":
			yield(&wind)
		case "thermal":
			yield(&thermal)
		}
	}
}

func yieldMany(yield func(v any) bool, values ...any) {
	for _, v := range values {
		if !yield(v) {
			return
		}
	}
}

type LoadLight struct {
	Substation string  `json:"substation"`
	Num        int     `json:"num"`
	NominalP   float64 `json:"nominalP"`
	Voltage    int     `json:"voltage"`
}

func (l *LoadLight) CimItems(modelId int) iter.Seq[any] {
	name := fmt.Sprintf("%s L%d", l.Substation, l.Num)

	bv := CreateBaseVoltage(l.Voltage)
	vl := CreateVoltageLevel(l.Substation, l.Voltage)

	var loadResponse models.LoadResponseCharacteristic
	loadResponse.Mrid = constantLoadMrid(l.NominalP)
	loadResponse.Name = fmt.Sprintf("Load response %.2f", l.NominalP)
	loadResponse.ShortName = loadResponse.Name

	var loadArea models.LoadArea
	loadArea.Mrid = loadAreaMrid("loads")
	loadArea.Name = "Load areas"

	subLoadName := fmt.Sprintf("Load %d kV", l.Voltage)
	var subLoadArea models.SubLoadArea
	subLoadArea.Mrid = subLoadAreaMrid(subLoadName)
	subLoadArea.Name = subLoadName
	subLoadArea.LoadAreaMrid = loadArea.Mrid

	var loadGroup models.LoadGroup
	loadGroup.Mrid = loadGroupMrid(name)
	loadGroup.Name = fmt.Sprintf("LG %s", name)
	loadGroup.SubLoadAreaMrid = subLoadArea.Mrid

	var conformLoad models.ConformLoad
	conformLoad.Name = name
	conformLoad.Mrid = loadMrid(name)
	conformLoad.ShortName = name
	conformLoad.Pfixed = l.NominalP
	conformLoad.BaseVoltageMrid = bv.Mrid
	conformLoad.EquipmentContainerMrid = vl.Mrid
	conformLoad.LoadResponseMrid = loadResponse.Mrid
	conformLoad.LoadGroupMrid = loadGroup.Mrid

	model := models.ModelEntity{ModelId: modelId}
	bvEntity := models.Entity{ModelEntity: model, Mrid: bv.Mrid, EntityType: StructName(bv)}
	vlEntity := models.Entity{ModelEntity: model, Mrid: vl.Mrid, EntityType: StructName(vl)}
	loadRespEntity := models.Entity{ModelEntity: model, Mrid: loadResponse.Mrid, EntityType: StructName(loadResponse)}
	conformLoadEntity := models.Entity{ModelEntity: model, Mrid: conformLoad.Mrid, EntityType: StructName(conformLoad)}
	loadGroupEntity := models.Entity{ModelEntity: model, Mrid: loadGroup.Mrid, EntityType: StructName(loadGroup)}
	loadAreaEntity := models.Entity{ModelEntity: model, Mrid: loadArea.Mrid, EntityType: StructName(loadArea)}
	subLoadAreaEntity := models.Entity{ModelEntity: model, Mrid: subLoadArea.Mrid, EntityType: StructName(subLoadArea)}

	return func(yield func(v any) bool) {
		yieldMany(
			yield, &bvEntity, &vlEntity, &loadRespEntity, &conformLoadEntity, &loadGroupEntity, &loadAreaEntity, &subLoadAreaEntity,
			&bv, &vl, &loadResponse, &conformLoad, &loadGroup, &loadArea, &subLoadArea)
	}
}
