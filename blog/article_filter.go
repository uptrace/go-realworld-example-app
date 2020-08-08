package blog

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/urlstruct"
	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"
)

type ArticleFilter struct {
	UserID    uint64
	Author    string
	Tag       string
	Favorited string
	urlstruct.Pager
}

func decodeArticleFilter(c *gin.Context) (*ArticleFilter, error) {
	f := &ArticleFilter{
		Tag:       c.Query("tag"),
		Author:    c.Query("author"),
		Favorited: c.Query("favorited"),
	}

	user, ok := c.Get("user")
	if ok {
		f.UserID = user.(*org.User).ID
	}

	return f, nil
}

func (f *ArticleFilter) query(q *orm.Query) (*orm.Query, error) {
	q = q.Relation("Author")

	{
		// subq := q.
		subq := rwe.PGMain().Model((*ArticleTag)(nil)).
			ColumnExpr("array_agg(t.tag)::text[]").
			Where("t.article_id = a.id")

		q = q.ColumnExpr("(?) AS tag_list", subq)
	}

	if f.UserID == 0 {
		q = q.ColumnExpr("false AS favorited")
	} else {
		subq := rwe.PGMain().Model((*FavoriteArticle)(nil)).
			Where("fa.article_id = a.id").
			Where("fa.user_id = ?", f.UserID)

		q = q.ColumnExpr("EXISTS (?) AS favorited", subq)
	}

	{
		subq := rwe.PGMain().Model((*FavoriteArticle)(nil)).
			ColumnExpr("count(*)").
			Where("fa.article_id = a.id")

		q = q.ColumnExpr("(?) AS favorites_count", subq)
	}

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
