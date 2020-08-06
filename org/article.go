package org

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/go-realworld-example-app/rwe"

	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/urlstruct"
)

type Article struct {
	tableName struct{} `pg:"articles,alias:a"`

	ID          uint64 `json:"-"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`

	Author   *Author `json:"author" pg:"-"`
	AuthorID uint64  `json:"-"`

	TagList []string `json:"tagList" pg:"-,array"`

	Favorited      bool `json:"favorited" pg:"-"`
	FavoritesCount int  `json:"favoritesCount" pg:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Author struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type ArticleTag struct {
	ArticleID uint64
	Tag       string
}

type ArticleFilter struct {
	Slug      string
	Author    string
	Tag       string
	Favorited string
	urlstruct.Pager
}

func SelectArticles(c *gin.Context, f *ArticleFilter) ([]*Article, error) {
	// authorQ := rwe.PGMain().ModelContext(c, (*User)(nil)).
	// 	ColumnExpr("username").
	// 	ColumnExpr("bio").
	// 	ColumnExpr("image").
	// 	// ColumnExpr("following").
	// 	Where("id = u.id")

	var articles []*Article
	err := rwe.PGMain().ModelContext(c, &articles).
		// ColumnExpr("(?) as author", authorQ).
		ColumnExpr("a.*").
		ColumnExpr("array_agg(tag) as tag_list").
		Join("JOIN article_tags AS at ON at.article_id = a.id").
		// Join("JOIN users AS u ON u.id = a.author_id").
		GroupExpr("a.id").
		Limit(f.Pager.GetLimit()).
		Offset(f.Pager.GetOffset()).
		Apply(f.Filters(c)).
		Select()
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (f *ArticleFilter) Filters(c *gin.Context) func(*orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if f.Slug != "" {
			q = q.Where("slug = ?", f.Slug)
			return q, nil
		}

		if f.Author != "" {
			q = q.Join("JOIN users as u ON u.id = a.author_id").
				Where("u.username = ?", f.Author)
		}

		if f.Tag != "" {
			q = q.Join("JOIN article_tags at u ON at.article_id = a.id").
				Where("at.tag = ?", f.Tag)
		}

		return q, nil
	}
}
