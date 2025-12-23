package models

import "time"

type Commit struct {
	Id        int64     `bun:"id,pk,autoincrement"`
	Message   string    `bun:"message"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp"`
}
