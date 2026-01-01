package models

type RdfsEnum struct {
	Id      int    `bun:"id,pk"`
	Code    string `bun:"code"`
	Iri     string `bun:"iri"`
	Comment string `bun:"comment"`
}

func (r RdfsEnum) GetId() int {
	return r.Id
}

func (r RdfsEnum) GetCode() string {
	return r.Code
}

type Enum interface {
	GetId() int
	GetCode() string
}

type ControlAreaTypeKind struct{ RdfsEnum }
type Currency struct{ RdfsEnum }
type CurveStyle struct{ RdfsEnum }
type DCConverterOperatingModeKind struct{ RdfsEnum }
type DCPolarityKind struct{ RdfsEnum }
type FuelType struct{ RdfsEnum }
type GeneratorControlSource struct{ RdfsEnum }
type HydroEnergyConversionKind struct{ RdfsEnum }
type HydroPlantStorageKind struct{ RdfsEnum }
type LimitTypeKind struct{ RdfsEnum }
type OperationalLimitDirectionKind struct{ RdfsEnum }
type PetersenCoilModeKind struct{ RdfsEnum }
type PhaseCode struct{ RdfsEnum }
type RegulatingControlModeKind struct{ RdfsEnum }
type SVCControlMode struct{ RdfsEnum }
type ShortCircuitRotorKind struct{ RdfsEnum }
type Source struct{ RdfsEnum }
type SynchronousMachineKind struct{ RdfsEnum }
type TransformerControlMode struct{ RdfsEnum }
type UnitMultiplier struct{ RdfsEnum }
type UnitSymbol struct{ RdfsEnum }
type Validity struct{ RdfsEnum }
type WindGenUnitKind struct{ RdfsEnum }
type WindingConnection struct{ RdfsEnum }
