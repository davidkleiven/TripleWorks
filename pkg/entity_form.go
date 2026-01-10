package pkg

import (
	"cmp"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"reflect"
	"slices"
	"strings"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
)

func FormTypes() map[string]any {
	return map[string]any{
		"ACDCConverter":                  &models.ACDCConverter{},
		"ACDCConverterDCTerminal":        &models.ACDCConverterDCTerminal{},
		"ACDCTerminal":                   &models.ACDCTerminal{},
		"ACLineSegment":                  &models.ACLineSegment{},
		"ActivePower":                    &models.ActivePower{},
		"ActivePowerPerCurrentFlow":      &models.ActivePowerPerCurrentFlow{},
		"ActivePowerPerFrequency":        &models.ActivePowerPerFrequency{},
		"AngleDegrees":                   &models.AngleDegrees{},
		"AngleRadians":                   &models.AngleRadians{},
		"ApparentPower":                  &models.ApparentPower{},
		"AsynchronousMachine":            &models.AsynchronousMachine{},
		"BaseVoltage":                    &models.BaseVoltage{},
		"BasicIntervalSchedule":          &models.BasicIntervalSchedule{},
		"Breaker":                        &models.Breaker{},
		"BusNameMarker":                  &models.BusNameMarker{},
		"BusbarSection":                  &models.BusbarSection{},
		"Capacitance":                    &models.Capacitance{},
		"CapacitancePerLength":           &models.CapacitancePerLength{},
		"Conductance":                    &models.Conductance{},
		"ConductingEquipment":            &models.ConductingEquipment{},
		"Conductor":                      &models.Conductor{},
		"ConformLoad":                    &models.ConformLoad{},
		"ConformLoadGroup":               &models.ConformLoadGroup{},
		"ConformLoadSchedule":            &models.ConformLoadSchedule{},
		"ConnectivityNode":               &models.ConnectivityNode{},
		"ConnectivityNodeContainer":      &models.ConnectivityNodeContainer{},
		"Connector":                      &models.Connector{},
		"ControlArea":                    &models.ControlArea{},
		"ControlAreaGeneratingUnit":      &models.ControlAreaGeneratingUnit{},
		"CsConverter":                    &models.CsConverter{},
		"CurrentFlow":                    &models.CurrentFlow{},
		"CurrentLimit":                   &models.CurrentLimit{},
		"Curve":                          &models.Curve{},
		"CurveData":                      &models.CurveData{},
		"DCBaseTerminal":                 &models.DCBaseTerminal{},
		"DCBreaker":                      &models.DCBreaker{},
		"DCBusbar":                       &models.DCBusbar{},
		"DCChopper":                      &models.DCChopper{},
		"DCConductingEquipment":          &models.DCConductingEquipment{},
		"DCConverterUnit":                &models.DCConverterUnit{},
		"DCDisconnector":                 &models.DCDisconnector{},
		"DCEquipmentContainer":           &models.DCEquipmentContainer{},
		"DCGround":                       &models.DCGround{},
		"DCLine":                         &models.DCLine{},
		"DCLineSegment":                  &models.DCLineSegment{},
		"DCNode":                         &models.DCNode{},
		"DCSeriesDevice":                 &models.DCSeriesDevice{},
		"DCShunt":                        &models.DCShunt{},
		"DCSwitch":                       &models.DCSwitch{},
		"DCTerminal":                     &models.DCTerminal{},
		"Disconnector":                   &models.Disconnector{},
		"EnergyConsumer":                 &models.EnergyConsumer{},
		"EnergySchedulingType":           &models.EnergySchedulingType{},
		"EnergySource":                   &models.EnergySource{},
		"Equipment":                      &models.Equipment{},
		"EquipmentContainer":             &models.EquipmentContainer{},
		"EquipmentVersion":               &models.EquipmentVersion{},
		"EquivalentBranch":               &models.EquivalentBranch{},
		"EquivalentEquipment":            &models.EquivalentEquipment{},
		"EquivalentInjection":            &models.EquivalentInjection{},
		"EquivalentNetwork":              &models.EquivalentNetwork{},
		"EquivalentShunt":                &models.EquivalentShunt{},
		"ExternalNetworkInjection":       &models.ExternalNetworkInjection{},
		"FossilFuel":                     &models.FossilFuel{},
		"Frequency":                      &models.Frequency{},
		"GeneratingUnit":                 &models.GeneratingUnit{},
		"GeographicalRegion":             &models.GeographicalRegion{},
		"HydroGeneratingUnit":            &models.HydroGeneratingUnit{},
		"HydroPowerPlant":                &models.HydroPowerPlant{},
		"HydroPump":                      &models.HydroPump{},
		"IdentifiedObject":               &models.IdentifiedObject{},
		"Inductance":                     &models.Inductance{},
		"InductancePerLength":            &models.InductancePerLength{},
		"InitialReactiveCapabilityCurve": &models.ReactiveCapabilityCurve{},
		"Junction":                       &models.Junction{},
		"Length":                         &models.Length{},
		"Line":                           &models.Line{},
		"LinearShuntCompensator":         &models.LinearShuntCompensator{},
		"LoadBreakSwitch":                &models.LoadBreakSwitch{},
		"LoadGroup":                      &models.LoadGroup{},
		"LoadResponseCharacteristic":     &models.LoadResponseCharacteristic{},
		"Money":                          &models.Money{},
		"NonConformLoad":                 &models.NonConformLoad{},
		"NonConformLoadGroup":            &models.NonConformLoadGroup{},
		"NonConformLoadSchedule":         &models.NonConformLoadSchedule{},
		"NonlinearShuntCompensator":      &models.NonlinearShuntCompensator{},
		"NonlinearShuntCompensatorPoint": &models.NonlinearShuntCompensatorPoint{},
		"NuclearGeneratingUnit":          &models.NuclearGeneratingUnit{},
		"OperationalLimit":               &models.OperationalLimit{},
		"OperationalLimitSet":            &models.OperationalLimitSet{},
		"OperationalLimitType":           &models.OperationalLimitType{},
		"PerLengthDCLineParameter":       &models.PerLengthDCLineParameter{},
		"PhaseTapChanger":                &models.PhaseTapChanger{},
		"PhaseTapChangerAsymmetrical":    &models.PhaseTapChangerAsymmetrical{},
		"PhaseTapChangerLinear":          &models.PhaseTapChangerLinear{},
		"PhaseTapChangerNonLinear":       &models.PhaseTapChangerNonLinear{},
		"PhaseTapChangerSymmetrical":     &models.PhaseTapChangerSymmetrical{},
		"PhaseTapChangerTable":           &models.PhaseTapChangerTable{},
		"PhaseTapChangerTablePoint":      &models.PhaseTapChangerTablePoint{},
		"PhaseTapChangerTabular":         &models.PhaseTapChangerTabular{},
		"PowerSystemResource":            &models.PowerSystemResource{},
		"PowerTransformer":               &models.PowerTransformer{},
		"PowerTransformerEnd":            &models.PowerTransformerEnd{},
		"ProtectedSwitch":                &models.ProtectedSwitch{},
		"RatioTapChanger":                &models.RatioTapChanger{},
		"RatioTapChangerTable":           &models.RatioTapChangerTable{},
		"RatioTapChangerTablePoint":      &models.RatioTapChangerTablePoint{},
		"Reactance":                      &models.Reactance{},
		"ReactiveCapabilityCurve":        &models.ReactiveCapabilityCurve{},
		"ReactivePower":                  &models.ReactivePower{},
		"RegularIntervalSchedule":        &models.RegularIntervalSchedule{},
		"RegulatingCondEq":               &models.RegulatingCondEq{},
		"RegulatingControl":              &models.RegulatingControl{},
		"ReportingGroup":                 &models.ReportingGroup{},
		"Resistance":                     &models.Resistance{},
		"ResistancePerLength":            &models.ResistancePerLength{},
		"RotatingMachine":                &models.RotatingMachine{},
		"RotationSpeed":                  &models.RotationSpeed{},
		"SeasonDayTypeSchedule":          &models.SeasonDayTypeSchedule{},
		"Seconds":                        &models.Seconds{},
		"SeriesCompensator":              &models.SeriesCompensator{},
		"ShuntCompensator":               &models.ShuntCompensator{},
		"Simple":                         &models.Simple_Float{},
		"SolarGeneratingUnit":            &models.SolarGeneratingUnit{},
		"StaticVarCompensator":           &models.StaticVarCompensator{},
		"SubGeographicalRegion":          &models.SubGeographicalRegion{},
		"Substation":                     &models.Substation{},
		"Susceptance":                    &models.Susceptance{},
		"Switch":                         &models.Switch{},
		"SynchronousMachine":             &models.SynchronousMachine{},
		"TapChanger":                     &models.TapChanger{},
		"TapChangerControl":              &models.TapChangerControl{},
		"TapChangerTablePoint":           &models.TapChangerTablePoint{},
		"Temperature":                    &models.Temperature{},
		"Terminal":                       &models.Terminal{},
		"ThermalGeneratingUnit":          &models.ThermalGeneratingUnit{},
		"TieFlow":                        &models.TieFlow{},
		"TransformerEnd":                 &models.TransformerEnd{},
		"Voltage":                        &models.Voltage{},
		"VoltageLevel":                   &models.VoltageLevel{},
		"VoltageLimit":                   &models.VoltageLimit{},
		"VoltagePerReactivePower":        &models.VoltagePerReactivePower{},
		"VsCapabilityCurve":              &models.VsCapabilityCurve{},
		"VsConverter":                    &models.VsConverter{},
		"WindGeneratingUnit":             &models.WindGeneratingUnit{},
	}
}

func EntityOptions(w io.Writer) {
	formTypes := FormTypes()
	items := make([]string, len(formTypes))
	num := 0
	for k := range formTypes {
		items[num] = k
		num++
	}
	slices.Sort(items)

	for _, k := range items {
		fmt.Fprintf(w, "<option>%s</option>\n", k)
	}
}

func FormInputFieldsForType(itemName string) (any, error) {
	fTypes := FormTypes()
	item, ok := fTypes[itemName]
	if !ok {
		return item, fmt.Errorf("Unknown type '%s'", itemName)
	}
	return item, nil
}

func FormInputFields(w io.Writer, item any) {
	checksum := MustGetHash(item)
	fmt.Fprintf(w, "<div id=\"entity-editor\" class=\"card-content\" checksum=\"%s\">\n", checksum)
	fields := FlattenStruct(item)
	fieldNames := make([]string, 0, len(fields))
	for name := range fields {
		fieldNames = append(fieldNames, name)
	}

	priorities := map[string]int{
		"Mrid":        0,
		"Name":        1,
		"ShortName":   2,
		"Description": 3,
	}
	slices.SortFunc(fieldNames, func(n1, n2 string) int {
		priority1, ok1 := priorities[n1]
		priority2, ok2 := priorities[n2]
		if ok1 && !ok2 {
			return -1
		} else if ok2 && !ok1 {
			return 1
		} else if ok1 && ok2 {
			return cmp.Compare(priority1, priority2)
		}
		return cmp.Compare(n1, n2)
	})

	for _, name := range fieldNames {
		if name == "Id" || name == "Deleted" {
			continue
		}
		formField := fields[name]
		value := formField.Value
		if formField.IsBunRelation {
			continue
		}

		fmt.Fprintf(w, "<div class=\"field\" id=\"%s-field\">\n", name)
		fmt.Fprintf(w, "<label id=\"%s-label\" class=\"label\" json-tag=\"%s\">%s</label>\n", name, formField.JsonTag, name)
		fmt.Fprintf(w, "<div class=\"control\">\n")

		borderStyle := getBorderStyle(name, value)
		if strings.Contains(name, "Mrid") {
			value = randomOrCurrentUuid(value)
		}

		writeCfg := writeInputConfig{
			name:        name,
			value:       value,
			borderStyle: borderStyle,
			queryKind:   mustGetQueryKind(name, fields),
		}
		writeInputItem(w, &writeCfg)
		fmt.Fprintf(w, "</div>\n</div>\n")
	}

	// Add delete marker
	deleteRow := `
	<div id="deleted-field" class="field has-background-danger-light">
		<label id="deleted-label" class="label" json-tag="delete">Mark for deletion</label>
		<input type="checkbox" id="deleted-value">
	</div>
	`
	fmt.Fprint(w, deleteRow)
	fmt.Fprintf(w, "</div>\n")
}

func getBorderStyle(name string, value any) string {
	var borderStyle string
	zero := uuid.UUID{}
	if strings.Contains(name, "Mrid") && value.(uuid.UUID) == zero {
		if name != "Mrid" {
			borderStyle = "is-danger"
		}
	}
	return borderStyle
}

func mustGetQueryKind(name string, fieldMap map[string]formField) string {
	var queryKind string
	if strings.HasSuffix(name, "Mrid") && name != "Mrid" {
		queryKind = strings.TrimSuffix(name, "Mrid")
	} else if strings.HasSuffix(name, "Id") {
		kindField := strings.TrimSuffix(name, "Id")
		val, ok := fieldMap[kindField]
		if !ok {
			panic(fmt.Sprintf("Kind field should exist for all id fields: %s (kindField=%s)", name, kindField))
		}
		queryKind = StructName(val.Value)
	}
	return queryKind
}

type writeInputConfig struct {
	name        string
	value       any
	borderStyle string
	queryKind   string
}

type inputData struct {
	Id          string
	Type        string
	Value       string
	Checked     bool
	Class       string
	Name        string
	SelectId    string
	SelectHxGet string
}

func writeInputItem(w io.Writer, config *writeInputConfig) {
	fieldType := reflect.TypeOf(config.value).Name()
	isEnum := strings.HasSuffix(config.name, "Id")
	isExternalId := strings.HasSuffix(config.name, "Mrid") && config.name != "Mrid"
	endpoint := ""
	fmt.Fprintf(w, "<div class=\"is-flex\">\n")

	if isEnum {
		endpoint = "enum"
	} else if isExternalId {
		endpoint = "entities"
	}

	data := inputData{
		Id:          fmt.Sprintf("%s-value", config.name),
		Type:        "text",
		Name:        config.name,
		Class:       strings.Trim(fmt.Sprintf("input %s", config.borderStyle), " "),
		Value:       fmt.Sprintf("%v", config.value),
		SelectId:    fmt.Sprintf("%s-select", config.name),
		SelectHxGet: fmt.Sprintf("/%s?kind=%s&choice=%v", endpoint, config.queryKind, config.value),
	}

	inputTemplate := `<input id="{{.Id}}" type="{{.Type}}" value="{{.Value}}" name="{{.Name}}" class="{{.Class}}" {{if .Checked}}checked{{end}}>`
	switch fieldType {
	case "bool":
		data.Checked = config.value.(bool)
		data.Type = "checkbox"
		data.Value = "true"
		data.Class = ""
	case "int", "float64":
		data.Type = "number"
	}

	templ := template.Must(template.New("inputField").Parse(inputTemplate))
	PanicOnErr(templ.Execute(w, data))

	if isExternalId || isEnum {
		selectHtml := `
<div class="select">
	<select id="{{.SelectId}}" value-field="{{.Id}}" hx-get="{{.SelectHxGet}}" hx-target="this" hx-swap="innerHTML" hx-trigger="load"></select>
</div>
		`
		selectTempl := template.Must(template.New("selectField").Parse(selectHtml))
		PanicOnErr(selectTempl.Execute(w, data))
	} else {
		maybeWriteAutofillCheckbox(w, config.name, config.value, data.Id)
	}
	fmt.Fprintf(w, "</div>\n")
}

func FlattenStruct(v any) map[string]formField {
	result := make(map[string]formField)
	flatten(reflect.ValueOf(v), result)
	return result
}

type formField struct {
	Value         any
	JsonTag       string
	IsBunRelation bool
}

func flatten(val reflect.Value, result map[string]formField) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	t := val.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		value := val.Field(i)
		if field.Name == "CommitId" {
			continue
		}

		tag := field.Tag.Get("bun")
		jsonTag := field.Tag.Get("json")

		if field.Anonymous && value.Kind() == reflect.Struct {
			flatten(value, result)
		} else {
			result[field.Name] = formField{
				Value:         value.Interface(),
				JsonTag:       jsonTag,
				IsBunRelation: strings.Contains(tag, "belongs-to"),
			}
		}
	}
}

func randomOrCurrentUuid(current any) uuid.UUID {
	asUuid, ok := current.(uuid.UUID)
	if !ok || (asUuid == uuid.UUID{}) {
		return Must(uuid.NewUUID())
	}
	return asUuid
}

func MustGetHash(item any) string {
	data := Must(json.Marshal(item))
	hash := md5.Sum(data)
	return string(hex.EncodeToString(hash[:]))
}

func isEmptyString(v any) bool {
	asString, ok := v.(string)
	return ok && asString == ""
}

func isNumberZero(v any) bool {
	asInt, ok := v.(int)
	asFloat, okFloat := v.(float64)
	return (ok && asInt == 0) || (okFloat && asFloat == 0.0)
}

func isBool(v any) bool {
	_, ok := v.(bool)
	return ok
}

func maybeWriteAutofillCheckbox(w io.Writer, name string, value any, valueId string) {
	if isBool(value) {
		return
	}
	skipAutoFill := map[string]struct{}{
		"Name":   {},
		"Length": {},
	}
	isMrid := strings.Contains(name, "Mrid")
	_, skip := skipAutoFill[name]
	if skip || isMrid {
		return
	}

	content := Must(io.ReadAll(Must(htmlPages.Open("html/autofill_checkbox.html"))))

	templ := Must(template.New("root").Parse(string(content)))
	data := NewAutofillCheckboxData(value, valueId, isEmptyString(value) || isNumberZero(value))
	templ.ExecuteTemplate(w, "autofill-checkbox", data)
}

type AutofillCheckboxData struct {
	Value    any
	TargetId string
	Checked  bool
	Checksum string
}

func NewAutofillCheckboxData(value any, targetId string, checked bool) AutofillCheckboxData {
	return AutofillCheckboxData{
		Value:    value,
		TargetId: targetId,
		Checked:  checked,
		Checksum: MustGetHash(value),
	}
}
