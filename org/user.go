package org

import "github.com/uptrace/go-realworld-example-app/rwe"

type User struct {
	tableName struct{} `pg:",alias:u"`

	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	// Image        *string `json:"image"`
	Password     string `json:"password" pg:"-"`
	PasswordHash string `json:"-"`
}

func SelectUser(id uint64) (*User, error) {
	user := new(User)
	if err := rwe.PGMain().
		Model(user).
		Where("id = ?", id).
		Select(); err != nil {
		return nil, err
	}

	return user, nil
}
