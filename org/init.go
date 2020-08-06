package org

import (
	"github.com/uptrace/go-realworld-example-app/rwe"
)

func init() {
	rwe.API.POST("/users", createUser)
	rwe.API.POST("/users/login", loginUser)

	rwe.API.GET("/articles", listArticles)
	rwe.API.GET("/articles/:slug", showArticle)

	rwe.API.Use(AuthMiddleware)
	rwe.API.GET("/user", currentUser)
	rwe.API.PUT("/users", updateUser)

	rwe.API.POST("/articles", createArticle)
}
