package org

import (
	"github.com/gin-gonic/gin"
	"github.com/uptrace/go-realworld-example-app/rwe"
)

func createArticle(c *gin.Context) {
	user := c.MustGet("user").(*User)

	article := new(Article)
	if err := c.BindJSON(article); err != nil {
		return
	}

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
