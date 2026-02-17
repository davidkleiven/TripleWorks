package repository

import (
	"context"

	"github.com/uptrace/bun"
)

type Inserter interface {
	Insert(ctx context.Context, item any) error
}

type InMemInserter struct {
	Items       []any
	InsertError error
}

func (i *InMemInserter) Insert(ctx context.Context, item any) error {
	i.Items = append(i.Items, item)
	return i.InsertError
}

type BunInserter struct {
	Db bun.IDB
}

func (b *BunInserter) Insert(ctx context.Context, item any) error {
	_, err := b.Db.NewInsert().Model(item).Exec(ctx)
	return err
}

func (b *BunInserter) InTx(ctx context.Context, fn InsertFunc) error {
	return b.Db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		txInserter := BunInserter{Db: tx}
		return fn(ctx, &txInserter)
	})
}

type InsertFunc func(context.Context, Inserter) error

type TxRunner interface {
	InTx(ctx context.Context, fn InsertFunc) error
}

func WithTx(fn InsertFunc) InsertFunc {
	return func(ctx context.Context, inserter Inserter) error {
		if txRunner, ok := inserter.(TxRunner); ok {
			return txRunner.InTx(ctx, fn)
		}
		return fn(ctx, inserter)
	}
}
