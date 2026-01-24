package models

type SimpleUpload struct {
	BaseEntity
	Data []byte `bun:"data"`
}
