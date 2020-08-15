package blog

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"
)

func listArticleComments(c *gin.Context) {
	article, err := SelectArticle(c, c.Param("slug"))
	if err != nil {
		c.Error(err)
		return
	}

	var userID uint64
	user, ok := c.Get("user")
	if ok {
		userID = user.(*org.User).ID
	}

	articles := make([]*Comment, 0)
	if err := rwe.PGMain().ModelContext(c, &articles).
		ColumnExpr("c.*").
		Relation("Author").
		Apply(authorFollowingColumn(userID)).
		Where("article_id = ?", article.ID).
		Select(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"comments": articles})
}

func showComment(c *gin.Context) {
	article, err := SelectArticle(c, c.Param("slug"))
	if err != nil {
		c.Error(err)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	var userID uint64
	user, ok := c.Get("user")
	if ok {
		userID = user.(*org.User).ID
	}

	comment := new(Comment)
	if err := rwe.PGMain().ModelContext(c, comment).
		ColumnExpr("c.*").
		Relation("Author").
		Apply(authorFollowingColumn(userID)).
		Where("c.id = ?", id).
		Where("article_id = ?", article.ID).
		Select(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"comment": comment})
}

func createComment(c *gin.Context) {
	user := c.MustGet("user").(*org.User)

	article, err := SelectArticle(c, c.Param("slug"))
	if err != nil {
		c.Error(err)
		return
	}

	var in struct {
		Comment *Comment `json:"comment"`
	}

	if err := c.BindJSON(&in); err != nil {
		return
	}

	if in.Comment == nil {
		c.Error(errors.New(`JSON field "comment" is required`))
		return
	}

	comment := in.Comment

	comment.AuthorID = user.ID
	comment.ArticleID = article.ID
	comment.CreatedAt = rwe.Clock.Now().UTC()
	comment.UpdatedAt = rwe.Clock.Now().UTC()

	if _, err := rwe.PGMain().
		ModelContext(c, comment).
		Insert(); err != nil {
		c.Error(err)
		return
	}

	comment.Author = org.NewProfile(user)
	c.JSON(200, gin.H{"comment": comment})
}

func deleteComment(c *gin.Context) {
	user := c.MustGet("user").(*org.User)

	article, err := SelectArticle(c, c.Param("slug"))
	if err != nil {
		c.Error(err)
		return
	}

	if _, err := rwe.PGMain().
		ModelContext(c, (*Comment)(nil)).
		Where("author_id = ?", user.ID).
		Where("article_id = ?", article.ID).
		Delete(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, nil)
}
