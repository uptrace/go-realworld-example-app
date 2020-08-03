package org

import (
	"github.com/uptrace/go-realworld-example-app/rwe"
)

func init() {
	rwe.API.POST("/users", createUser)
	rwe.API.POST("/users/login", loginUser)

	rwe.API.Use(AuthMiddleware)
	rwe.API.POST("/user", currentUser)
}
