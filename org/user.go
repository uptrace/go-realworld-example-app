package org

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
