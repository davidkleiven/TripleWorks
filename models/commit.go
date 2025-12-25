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
