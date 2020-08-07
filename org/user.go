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
	PasswordHash string `json:"-"`

	Token string `pg:"-" json:"token,omitempty"`
}

type Profile struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
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
