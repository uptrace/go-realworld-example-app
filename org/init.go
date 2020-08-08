package org

import (
	"github.com/uptrace/go-realworld-example-app/rwe"
)

func init() {
	g := rwe.API.Group("")

	g.Use(UserMiddleware)

	g.POST("/users", createUser)
	g.POST("/users/login", loginUser)

	g.Use(MustUserMiddleware)

	g.GET("/user", currentUser)
	g.PUT("/users", updateUser)
}
