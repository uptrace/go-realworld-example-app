package blog

import (
	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"
)

func init() {
	rwe.API.GET("/articles", listArticles)
	rwe.API.GET("/articles/:slug", showArticle)

	rwe.API.Use(org.AuthMiddleware)

	rwe.API.POST("/articles", createArticle)
}
