package repository

import (
	"context"
	"errors"

	"com.github/davidkleiven/tripleworks/models"
	"github.com/uptrace/bun"
)

var ErrUserNotFound = errors.New("user not found")

type BunUserRepository struct {
	Db *bun.DB
}

func (r *BunUserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.Db.NewSelect().Model(&user).Where("email = ?", email).Limit(1).Scan(ctx)
	if err != nil {
		return user, ErrUserNotFound
	}
	return user, nil
}
