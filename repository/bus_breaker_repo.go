package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type BusBreakerConnection struct {
	Mrid           uuid.UUID `bun:"mrid"`
	R              float64   `bun:"r"`
	X              float64   `bun:"x"`
	Name           string    `bun:"name"`
	NominalVoltage float64   `bun:"nominal_voltage"`
	SubstationMrid uuid.UUID `bun:"substation_mrid"`
}

type BusBreakerRepo interface {
	Fetch(ctx context.Context) ([]BusBreakerConnection, error)
}

type CachedBusbReakerrepo struct {
	Items []BusBreakerConnection
}

func (c *CachedBusbReakerrepo) Fetch(ctx context.Context) ([]BusBreakerConnection, error) {
	return c.Items, nil
}

type BunBusBreakerRepo struct {
	Db *bun.DB
}

func (b *BunBusBreakerRepo) Fetch(ctx context.Context) ([]BusBreakerConnection, error) {
	query, err := sqlFS.ReadFile("sql/bus_breaker.sql")
	var result []BusBreakerConnection
	if err != nil {
		return result, fmt.Errorf("Failed to open query: %w", err)
	}
	err = b.Db.NewRaw(string(query)).Scan(ctx, &result)
	return result, err
}
