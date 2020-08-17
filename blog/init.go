package blog

import (
	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"
)

func init() {
	g := rwe.API.Group("")

	g.Use(org.UserMiddleware)

	g.GET("/tags/", listTags)
	g.GET("/articles", listArticles)
	g.GET("/articles/:slug", showArticle)
	g.GET("/articles/:slug/comments", listArticleComments)
	g.GET("/articles/:slug/comments/:id", showComment)

	g.Use(org.MustUserMiddleware)

	g.POST("/articles", createArticle)
	g.PUT("/articles/:slug", updateArticle)
	g.DELETE("/articles/:slug", deleteArticle)

	g.POST("/articles/:slug/favorite", favoriteArticle)
	g.DELETE("/articles/:slug/favorite", unfavoriteArticle)

	g.POST("/articles/:slug/comments", createComment)
	g.DELETE("/articles/:slug/comments/:id", deleteComment)
}
