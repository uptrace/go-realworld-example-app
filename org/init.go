package org

import (
	"github.com/gin-gonic/gin"

	"github.com/uptrace/go-realworld-example-app/rwe"
)

func init() {
	rwe.API.GET("/users", m(listUsers))
	rwe.API.GET("/users/:user_id", m(showUser))
	rwe.API.POST("/users", m(createUser))
}

type handler func(c *gin.Context) error

func m(fn handler) func(c *gin.Context) {
	return func(c *gin.Context) {
		if err := fn(c); err != nil {
			c.Error(err)
		}
	}
}
