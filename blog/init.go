package blog

import (
	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"
)

func init() {
	g := rwe.API.Group("")

	g.Use(org.UserMiddleware)

	g.GET("/articles", listArticles)
	g.GET("/articles/:slug", showArticle)

	g.Use(org.AuthMiddleware)

	g.POST("/articles", createArticle)
	g.PUT("/articles/:slug", updateArticle)

	g.POST("/articles/:slug/favorite", favoriteArticle)
	g.DELETE("/articles/:slug/favorite", unfavoriteArticle)
}
