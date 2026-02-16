package repository

import (
	"context"

	"com.github/davidkleiven/tripleworks/models"
)

type ConnectivityNodeReadRepository interface {
	ReadRepository[models.ConnectivityNode]
	InContainer(ctx context.Context, mrid string) ([]models.ConnectivityNode, error)
}

type InMemConnectivityNodeReadRepository struct {
	InMemReadRepository[models.ConnectivityNode]
}

func (imc *InMemConnectivityNodeReadRepository) InContainer(ctx context.Context, mrid string) ([]models.ConnectivityNode, error) {
	var result []models.ConnectivityNode
	for _, item := range imc.Items {
		if item.ConnectivityNodeContainerMrid.String() == mrid {
			result = append(result, item)
		}
	}
	return result, nil
}

type BunConnectivityNodeReadRepository struct {
	BunReadRepository[models.ConnectivityNode]
}

func (bcr *BunConnectivityNodeReadRepository) InContainer(ctx context.Context, mrid string) ([]models.ConnectivityNode, error) {
	var result []models.ConnectivityNode
	err := bcr.Db.NewSelect().Model(&result).Where("connectivity_node_container_mrid = ?", mrid).Scan(ctx)
	return result, err
}
