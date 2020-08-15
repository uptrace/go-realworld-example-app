package blog

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
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

	c.JSON(200, gin.H{
		"articles":      articles,
		"articlesCount": len(articles),
	})
}

func showArticle(c *gin.Context) {
	slug := c.Param("slug")

	if slug == "feed" {
		listArticlesFeed(c)
		return
	}

	article, err := selectArticleByFilter(c)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"article": article})
}

func selectArticleByFilter(c *gin.Context) (*Article, error) {
	f, err := decodeArticleFilter(c)
	if err != nil {
		return nil, err
	}

	article := new(Article)
	if err := rwe.PGMain().
		ModelContext(c, article).
		ColumnExpr("?TableColumns").
		Apply(f.query).
		Select(); err != nil {
		return nil, err
	}

	if article.TagList == nil {
		article.TagList = make([]string, 0)
	}

	return article, nil
}

func listArticlesFeed(c *gin.Context) {
	e, ok := c.Get("authErr")
	if ok {
		c.AbortWithError(http.StatusUnauthorized, e.(error))
		return
	}

	f, err := decodeArticleFilter(c)
	if err != nil {
		c.Error(err)
		return
	}

	articles := make([]*Article, 0)
	if err := rwe.PGMain().
		ModelContext(c, &articles).
		ColumnExpr("?TableColumns").
		Apply(f.query).
		Select(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"articles":      articles,
		"articlesCount": len(articles),
	})
}

func createArticle(c *gin.Context) {
	user := c.MustGet("user").(*org.User)

	var in struct {
		Article *Article `json:"article"`
	}

	if err := c.BindJSON(&in); err != nil {
		return
	}

	if in.Article == nil {
		c.Error(errors.New(`JSON field "article" is required`))
		return
	}

	article := in.Article

	article.Slug = makeSlug(article.Title)
	article.AuthorID = user.ID
	article.CreatedAt = rwe.Clock.Now().UTC()
	article.UpdatedAt = rwe.Clock.Now().UTC()

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

	var in struct {
		Article *Article `json:"article"`
	}

	if err := c.BindJSON(&in); err != nil {
		return
	}

	if in.Article == nil {
		c.Error(errors.New(`JSON field "article" is required`))
		return
	}

	article := in.Article

	if _, err := rwe.PGMain().
		ModelContext(c, article).
		Set("title = ?", article.Title).
		Set("description = ?", article.Description).
		Set("body = ?", article.Body).
		Set("updated_at = ?", rwe.Clock.Now().UTC()).
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

	if article.TagList == nil {
		article.TagList = make([]string, 0)
	}

	article.Author = org.NewProfile(user)
	c.JSON(200, gin.H{"article": article})
}

func deleteArticle(c *gin.Context) {
	user := c.MustGet("user").(*org.User)

	if _, err := rwe.PGMain().
		ModelContext(c, (*Article)(nil)).
		Where("author_id = ?", user.ID).
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

	article, err := selectArticleByFilter(c)
	if err != nil {
		c.Error(err)
		return
	}

	favoriteArticle := &FavoriteArticle{
		UserID:    user.ID,
		ArticleID: article.ID,
	}
	res, err := rwe.PGMain().
		ModelContext(c, favoriteArticle).
		Insert()
	if err != nil {
		c.Error(err)
		return
	}

	if res.RowsAffected() != 0 {
		article.Favorited = true
		article.FavoritesCount = article.FavoritesCount + 1
	}
	c.JSON(200, gin.H{"article": article})
}

func unfavoriteArticle(c *gin.Context) {
	user := c.MustGet("user").(*org.User)

	article, err := selectArticleByFilter(c)
	if err != nil {
		c.Error(err)
		return
	}

	res, err := rwe.PGMain().
		ModelContext(c, (*FavoriteArticle)(nil)).
		Where("user_id = ?", user.ID).
		Where("article_id = ?", article.ID).
		Delete()
	if err != nil {
		c.Error(err)
		return
	}

	if res.RowsAffected() != 0 {
		article.Favorited = false
		article.FavoritesCount = article.FavoritesCount - 1
	}
	c.JSON(200, gin.H{"article": article})
}

func listTags(c *gin.Context) {
	tags := make([]string, 0)
	if err := rwe.PGMain().ModelContext(c, (*ArticleTag)(nil)).
		// ColumnExpr("distinct(tag)").
		ColumnExpr("tag").
		GroupExpr("tag").
		OrderExpr("count(tag) DESC").
		ColumnExpr("tag").
		Select(&tags); err != nil && err != pg.ErrNoRows {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"tags": tags})
}
