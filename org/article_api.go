package org

import (
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/uptrace/go-realworld-example-app/rwe"
)

const charsBytes = "01234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charsBytes[rand.Intn(len(charsBytes))]
	}
	return string(b)
}

func createArticle(c *gin.Context) {
	user := c.MustGet("user").(*User)

	article := new(Article)
	if err := c.BindJSON(article); err != nil {
		return
	}

	article.Slug = slug.Make(article.Title + " " + randStringBytes(6))

	article.AuthorID = user.ID
	if _, err := rwe.PGMain().
		ModelContext(c, article).
		Insert(); err != nil {
		c.Error(err)
		return
	}

	tags := make([]ArticleTag, 0, len(article.TagList))
	for _, t := range article.TagList {
		tags = append(tags, ArticleTag{
			ArticleID: article.ID,
			Tag:       t,
		})
	}

	if _, err := rwe.PGMain().
		ModelContext(c, &tags).
		Insert(); err != nil {
		c.Error(err)
		return
	}

	article.Author = &Author{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: false,
	}
	c.JSON(200, gin.H{"article": article})
}

func showArticle(c *gin.Context) {
	slug := c.Param("slug")

	article := new(Article)
	if err := rwe.PGMain().
		ModelContext(c, article).
		ColumnExpr("a.*").
		ColumnExpr("array_agg(tag) as tag_list").
		Join("JOIN article_tags AS at ON at.article_id = a.id").
		Where("slug = ?", slug).
		GroupExpr("a.id").
		Select(); err != nil {
		c.Error(err)
		return
	}

	user, err := SelectUser(c, article.AuthorID)
	if err != nil {
		c.Error(err)
		return
	}

	article.Author = &Author{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: false,
	}
	c.JSON(200, gin.H{"article": article})
}

func listArticles(c *gin.Context) {
	articleFilter := &ArticleFilter{
		Tag:       c.Query("tag"),
		Author:    c.Query("author"),
		Favorited: c.Query("favorited"),
	}

	articles, err := SelectArticles(c, articleFilter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"articles": articles})
}
