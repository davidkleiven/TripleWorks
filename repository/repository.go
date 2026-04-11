package repository

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"reflect"
	"slices"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

type Lister[T any] interface {
	List(ctx context.Context) ([]T, error)
}

type ReadRepository[T any] interface {
	Lister[T]
	GetByMrid(ctx context.Context, mrid string) (T, error)
	ListByMrids(ctx context.Context, mrids iter.Seq[string]) ([]T, error)
}

type InMemReadRepository[T models.VersionedIdentifiedObject] struct {
	Items []T
}

func (i *InMemReadRepository[T]) GetByMrid(ctx context.Context, mrid string) (T, error) {
	var (
		newest T
		found  bool
	)
	for _, candidate := range i.Items {
		if candidate.GetMrid().String() == mrid && candidate.GetCommitId() >= newest.GetCommitId() {
			found = true
			newest = candidate
		}
	}

	if !found {
		return newest, fmt.Errorf("No resource with mrid %s", mrid)
	}
	return newest, nil
}

func (i *InMemReadRepository[T]) List(ctx context.Context) ([]T, error) {
	data := make([]T, 0, len(i.Items))
	for _, item := range i.Items {
		data = append(data, item)
	}
	return data, nil
}

func (i *InMemReadRepository[T]) ListByMrids(ctx context.Context, mrids iter.Seq[string]) ([]T, error) {
	mridSet := make(map[string]struct{})
	for mrid := range mrids {
		mridSet[mrid] = struct{}{}
	}

	var result []T
	for _, item := range i.Items {
		_, ok := mridSet[item.GetMrid().String()]
		if ok {
			result = append(result, item)
		}
	}
	return result, nil
}

type BunReadRepository[T any] struct {
	Db            *bun.DB
	UseLatestView bool
}

func (brp *BunReadRepository[T]) GetByMrid(ctx context.Context, mrid string) (T, error) {
	var result T
	latestView := brp.TableName()
	err := brp.Db.NewSelect().Table(latestView).Where("mrid = ?", mrid).OrderBy("commit_id", bun.OrderDesc).Limit(1).Scan(ctx, &result)
	return result, err
}

func (brp *BunReadRepository[T]) List(ctx context.Context) ([]T, error) {
	var result []T
	latestView := brp.TableName()
	err := brp.Db.NewSelect().Table(latestView).Scan(ctx, &result)
	return result, err
}

func (brp *BunReadRepository[T]) ListByMrids(ctx context.Context, mrids iter.Seq[string]) ([]T, error) {
	var result []T
	latestView := brp.TableName()
	err := brp.Db.NewSelect().Table(latestView).Where("mrid IN (?)", bun.List(slices.Collect(mrids))).Scan(ctx, &result)
	return result, err
}

func (brp *BunReadRepository[T]) TableName() string {
	var item T
	tableName := brp.Db.Table(reflect.TypeOf(item)).Name
	if brp.UseLatestView {
		return fmt.Sprintf("v_%s_latest", tableName)
	}
	return tableName
}

type FailingReadRepo[T any] struct{}

func (f *FailingReadRepo[T]) GetByMrid(ctx context.Context, mrid string) (T, error) {
	var item T
	return item, fmt.Errorf("Failed to read mrid: %s", mrid)
}

func (f *FailingReadRepo[T]) List(ctx context.Context) ([]T, error) {
	var result []T
	return result, errors.New("failed to list items")
}

func (f *FailingReadRepo[T]) ListByMrids(ctx context.Context, mrids iter.Seq[string]) ([]T, error) {
	var items []T
	return items, fmt.Errorf("Failed to list items: %v", mrids)
}

type InMemLister[T any] struct {
	Items []T
	Err   error
}

func (i *InMemLister[T]) List(ctx context.Context) ([]T, error) {
	return i.Items, i.Err
}
