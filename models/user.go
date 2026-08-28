package models

type User struct {
	Id    int    `bun:"id,pk,autoincrement"`
	Email string `bun:"email,unique,notnull"`
}
