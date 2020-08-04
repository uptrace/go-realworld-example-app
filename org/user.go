package org

import (
	"context"

	"github.com/uptrace/go-realworld-example-app/rwe"
)

type User struct {
	tableName struct{} `pg:",alias:u"`

	ID           uint64 `json:"-"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Bio          string `json:"bio"`
	Image        string `json:"image"`
	Password     string `pg:"-" json:"password,omitempty"`
	Token        string `pg:"-" json:"token,omitempty"`
	PasswordHash string `json:"-"`
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
