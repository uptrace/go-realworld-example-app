package blog

import (
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/urlstruct"
	"github.com/gosimple/slug"
	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"
)

const charsBytes = "01234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type ArticleFilter struct {
	Author    string
	Tag       string
	Favorited string
	urlstruct.Pager
}

func (f *ArticleFilter) Filters(q *orm.Query) (*orm.Query, error) {
	if f.Author != "" {
		q = q.Where("author__username = ?", f.Author)
	}

	if f.Tag != "" {
		q = q.
			Join("JOIN article_tags AS t ON t.article_id = a.id").
			Where("t.tag = ?", f.Tag)
	}

	return q, nil
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charsBytes[rand.Intn(len(charsBytes))]
	}
	return string(b)
}

func createArticle(c *gin.Context) {
	user := c.MustGet("user").(*org.User)

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
	article := new(Article)
	if err := rwe.PGMain().
		ModelContext(c, article).
		ColumnExpr("?TableColumns").
		Relation("Author").
		Apply(articleTagsSubquery).
		Where("slug = ?", c.Param("slug")).
		Select(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"article": article})
}

func articleTagsSubquery(q *orm.Query) (*orm.Query, error) {
	subq := rwe.PGMain().Model((*ArticleTag)(nil)).
		ColumnExpr("array_agg(t.tag)::text[]").
		Where("t.article_id = a.id")

	return q.ColumnExpr("(?) AS tag_list", subq), nil
}

func listArticles(c *gin.Context) {
	f := &ArticleFilter{
		Tag:       c.Query("tag"),
		Author:    c.Query("author"),
		Favorited: c.Query("favorited"),
	}

	articles := make([]*Article, 0)
	err := rwe.PGMain().ModelContext(c, &articles).
		ColumnExpr("?TableColumns").
		Apply(articleTagsSubquery).
		Apply(f.Filters).
		Relation("Author").
		Limit(f.Pager.GetLimit()).
		Offset(f.Pager.GetOffset()).
		Select()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"articles": articles})
}
