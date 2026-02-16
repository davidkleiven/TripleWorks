package repository

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

type VoltageLevelReadRepository interface {
	ReadRepository[models.VoltageLevel]
	InSubstation(ctx context.Context, mrid string) ([]models.VoltageLevel, error)
}

type InMemVoltageLevelReadRepository struct {
	InMemReadRepository[models.VoltageLevel]
	InSubstationErr error
}

func (imv *InMemVoltageLevelReadRepository) InSubstation(ctx context.Context, mrid string) ([]models.VoltageLevel, error) {
	var result []models.VoltageLevel
	for _, vl := range imv.items {
		if vl.SubstationMrid.String() == mrid {
			result = append(result, vl)
		}
	}
	return result, imv.InSubstationErr
}

func NewBunVoltageLevelReadRepository(db *bun.DB) *BunVoltageLevelRepository {
	return &BunVoltageLevelRepository{
		BunReadRepository: BunReadRepository[models.VoltageLevel]{db: db},
	}
}

type BunVoltageLevelRepository struct {
	BunReadRepository[models.VoltageLevel]
}

func (vlr *BunVoltageLevelRepository) InSubstation(ctx context.Context, mrid string) ([]models.VoltageLevel, error) {
	var vls []models.VoltageLevel
	err := vlr.db.NewSelect().Model(&vls).Where("substation_mrid = ?", mrid).Scan(ctx)
	return vls, err
}
