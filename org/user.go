package org

import (
	"context"

	"github.com/uptrace/go-realworld-example-app/rwe"
)

type User struct {
	tableName struct{} `pg:",alias:u"`

	ID           uint64
	Username     string
	Email        string
	Bio          string
	Image        string `pg:"img"`
	Password     string `pg:"-"`
	PasswordHash string
}

func SelectUser(ctx context.Context, id uint64) (*User, error) {
	user := new(User)
	if err := rwe.PGMain().
		ModelContext(ctx, user).
		Where("id = ?", id).
		Select(); err != nil {
		return nil, err
	}

	return user, nil
}
