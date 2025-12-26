package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Commit struct {
	bun.BaseModel `bun:"table:commits"`
	Id            int64     `bun:"id,pk,autoincrement"`
	Branch        string    `bun:"branch,default:'main'"`
	Message       string    `bun:"message"`
	Author        string    `bun:"author"`
	CreatedAt     time.Time `bun:"created_at,default:current_timestamp"`
}

type Model struct {
	bun.BaseModel `bun:"table:models"`
	Id            int    `bun:"id,pk,autoincrement"`
	Name          string `bun:"name"`
}

type BaseEntity struct {
	Id       int64 `bun:"id,pk,autoincrement"`
	ModelId  int64 `bun:"model_id"`
	CommitId int64 `bun:"commit_id"`

	Commit *Commit `bun:"rel:belongs-to,join:commit_id=id"`
	Model  *Model  `bun:"rel:belongs-to,join:model_id=id"`
}
