package repository

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type ReadRepository[T any] interface {
	GetByMrid(ctx context.Context, mrid string) (T, error)
	List(ctx context.Context) ([]T, error)
}

type InMemReadRepository[T any] struct {
	items map[string]T
}

func (i *InMemReadRepository[T]) GetByMrid(ctx context.Context, mrid string) (T, error) {
	result, ok := i.items[mrid]
	if !ok {
		return result, fmt.Errorf("No object at %s", mrid)
	}
	return result, nil
}

func (i *InMemReadRepository[T]) List(ctx context.Context) ([]T, error) {
	data := make([]T, 0, len(i.items))
	for _, item := range i.items {
		data = append(data, item)
	}
	return data, nil
}

type BunReadRepository[T any] struct {
	db *bun.DB
}

func (brp *BunReadRepository[T]) GetByMrid(ctx context.Context, mrid string) (T, error) {
	var result T
	err := brp.db.NewSelect().Model(&result).Where("mrid = ?", mrid).Scan(ctx)
	return result, err
}

func (brp *BunReadRepository[T]) List(ctx context.Context) ([]T, error) {
	var result []T
	err := brp.db.NewSelect().Model(&result).Scan(ctx)
	return result, err
}
