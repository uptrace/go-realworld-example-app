package org

import (
	"github.com/uptrace/go-realworld-example-app/rwe"
)

func init() {
	rwe.API.GET("/users", listUsers)
	rwe.API.POST("/users", createUser)
	rwe.API.POST("/users/login", loginUser)

	rwe.API.Use(AuthMiddleware())
	rwe.API.POST("/users/current", currentUser)
}
