package repository

import (
	"context"
	"iter"
	"slices"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

type TerminalReadRepository interface {
	ReadRepository[models.Terminal]
	WithConnectivityNode(ctx context.Context, mrids iter.Seq[string]) ([]models.Terminal, error)
}

type InMemTerminalReadRepository struct {
	InMemReadRepository[models.Terminal]
	WithConNodeErr error
}

func (imt *InMemTerminalReadRepository) WithConnectivityNode(ctx context.Context, mrids iter.Seq[string]) ([]models.Terminal, error) {
	mridSet := make(map[string]struct{})
	for mrid := range mrids {
		mridSet[mrid] = struct{}{}
	}

	var result []models.Terminal
	for _, item := range imt.Items {
		_, ok := mridSet[item.ConnectivityNodeMrid.String()]
		if ok {
			result = append(result, item)
		}
	}
	return result, imt.WithConNodeErr
}

type BunTerminalReadRepository struct {
	BunReadRepository[models.Terminal]
}

func (btr *BunTerminalReadRepository) WithConnectivityNode(ctx context.Context, mrids iter.Seq[string]) ([]models.Terminal, error) {
	mridsSlice := slices.Collect(mrids)
	var results []models.Terminal
	err := btr.Db.NewSelect().Model(&results).Where("connectivity_node_mrid IN (?)", bun.In(mridsSlice)).Scan(ctx)
	return results, err
}
