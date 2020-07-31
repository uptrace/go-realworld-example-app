package org

import "github.com/uptrace/go-realworld-example-app/rwe"

func init() {
	rwe.API.GET("/users", listUsers)
	rwe.API.GET("/users/:user_id", showUser)
	rwe.API.POST("/users", createUser)
}
