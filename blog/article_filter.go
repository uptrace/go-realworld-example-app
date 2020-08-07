package blog

import (
	"github.com/go-pg/pg/v10/orm"
	"github.com/go-pg/urlstruct"
)

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
