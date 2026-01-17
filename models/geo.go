package models

import (
	"github.com/google/uuid"
)

type LocatedPowerSystemResource struct {
	LocationMrid uuid.UUID `bun:"location_mrid,type:uuid" json:"location_mrid"`
	Location     *Entity   `bun:"rel:belongs-to,join:location_mrid=mrid" json:"location,omitempty"`
}

type Location struct {
	IdentifiedObject
	CoordinateSystemMrid uuid.UUID `bun:"coordinate_system_mrid,type:uuid" json:"coordinate_system_mrid"`
	CoordinateSystem     *Entity   `bun:"rel:belongs-to,join:coordinate_system_mrid=mrid" json:"coordinate_system,omitempty"`
}
type PositionPoint struct {
	YPosition      float64   `bun:"yposition" json:"yposition" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PositionPoint.yPosition"`
	ZPosition      float64   `bun:"zposition" json:"zposition" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PositionPoint.zPosition"`
	SequenceNumber int       `bun:"sequence_number" json:"sequence_number" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PositionPoint.sequenceNumber"`
	XPosition      float64   `bun:"xposition" json:"xposition" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#PositionPoint.xPosition"`
	LocationMrid   uuid.UUID `bun:"location_mrid,type:uuid" json:"location_mrid"`
	Location       *Entity   `bun:"rel:belongs-to,join:location_mrid=mrid" json:"location,omitempty"`
}
type CoordinateSystem struct {
	IdentifiedObject
	CrsUrn string `bun:"crs_urn" json:"crs_urn" iri:"http://iec.ch/TC57/2013/CIM-schema-cim16#CoordinateSystem.crsUrn"`
}
