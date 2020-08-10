package blog

import (
	"math/rand"
	"strconv"

	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func makeSlug(title string) string {
	return slug.Make(strconv.Itoa(rand.Int()) + " " + title)
}

func listArticles(c *gin.Context) {
	f, err := decodeArticleFilter(c)
	if err != nil {
		c.Error(err)
		return
	}

	articles := make([]*Article, 0)
	if err := rwe.PGMain().ModelContext(c, &articles).
		ColumnExpr("?TableColumns").
		Apply(f.query).
		Limit(f.Pager.GetLimit()).
		Offset(f.Pager.GetOffset()).
		Select(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"articles": articles})
}

func showArticle(c *gin.Context) {
	f, err := decodeArticleFilter(c)
	if err != nil {
		c.Error(err)
		return
	}

	article := new(Article)
	if err := rwe.PGMain().
		ModelContext(c, article).
		ColumnExpr("?TableColumns").
		Apply(f.query).
		Where("slug = ?", c.Param("slug")).
		Select(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"article": article})
}

func createArticle(c *gin.Context) {
	user := c.MustGet("user").(*org.User)

	article := new(Article)
	if err := c.BindJSON(article); err != nil {
		return
	}

	article.Slug = makeSlug(article.Title)
	article.AuthorID = user.ID
	article.CreatedAt = rwe.Clock.Now()

	if _, err := rwe.PGMain().
		ModelContext(c, article).
		Insert(); err != nil {
		c.Error(err)
		return
	}

	if err := createTags(c, article); err != nil {
		c.Error(err)
		return
	}

	article.Author = org.NewProfile(user)
	c.JSON(200, gin.H{"article": article})
}

func updateArticle(c *gin.Context) {
	user := c.MustGet("user").(*org.User)

	article := new(Article)
	if err := c.BindJSON(article); err != nil {
		return
	}

	if _, err := rwe.PGMain().
		ModelContext(c, article).
		Set("title = ?", article.Title).
		Set("slug = ?", makeSlug(article.Title)).
		Set("description = ?", article.Description).
		Set("body = ?", article.Body).
		Set("updated_at = ?", rwe.Clock.Now()).
		Where("slug = ?", c.Param("slug")).
		Returning("*").
		Update(); err != nil {
		c.Error(err)
		return
	}

	if _, err := rwe.PGMain().ModelContext(c, (*ArticleTag)(nil)).
		Where("article_id = ?", article.ID).
		Delete(); err != nil {
		c.Error(err)
		return
	}

	if err := createTags(c, article); err != nil {
		c.Error(err)
		return
	}

	article.Author = org.NewProfile(user)
	c.JSON(200, gin.H{"article": article})
}

func deleteArticle(c *gin.Context) {
	if _, err := rwe.PGMain().
		ModelContext(c, (*Article)(nil)).
		Where("slug = ?", c.Param("slug")).
		Delete(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, nil)
}

func createTags(c *gin.Context, article *Article) error {
	if len(article.TagList) == 0 {
		return nil
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
		return err
	}

	return nil
}

func favoriteArticle(c *gin.Context) {
	user := c.MustGet("user").(*org.User)
	article, err := SelectArticle(c, c.Param("slug"))
	if err != nil {
		c.Error(err)
		return
	}

	favoriteArticle := &FavoriteArticle{
		UserID:    user.ID,
		ArticleID: article.ID,
	}
	if _, err := rwe.PGMain().
		ModelContext(c, favoriteArticle).
		Insert(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"article": article})
}

func unfavoriteArticle(c *gin.Context) {
	user := c.MustGet("user").(*org.User)
	article, err := SelectArticle(c, c.Param("slug"))
	if err != nil {
		c.Error(err)
		return
	}

	if _, err := rwe.PGMain().
		ModelContext(c, (*FavoriteArticle)(nil)).
		Where("user_id = ?", user.ID).
		Where("article_id = ?", article.ID).
		Delete(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"article": article})
}
