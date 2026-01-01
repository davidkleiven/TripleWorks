package models

import (
	"time"

	"github.com/google/uuid"
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

type ModelEntity struct {
	ModelId int    `bun:"model_id"`
	Model   *Model `bun:"rel:belongs-to,join:model_id=id"`
}

type BaseEntity struct {
	Id       int     `bun:"id,pk,autoincrement"`
	CommitId int     `bun:"commit_id"`
	Commit   *Commit `bun:"rel:belongs-to,join:commit_id=id"`
	Deleted  bool    `bun:"deleted" json:"deleted"`
}

func (b *BaseEntity) SetCommitId(commitId int) {
	b.CommitId = commitId
}

func (b BaseEntity) GetCommitId() int {
	return b.CommitId
}

type MridGetter interface {
	GetMrid() uuid.UUID
}

type VersionedIdentifiedObject interface {
	MridGetter
	GetCommitId() int
}

type MridNameGetter interface {
	MridGetter
	GetName() string
}

type CommitIdSetter interface {
	SetCommitId(commitId int)
}
