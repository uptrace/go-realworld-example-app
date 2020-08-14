package blog

import (
	"context"
	"encoding/json"
	"time"

	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"
)

const TimeFormatStr = "2006-01-02T15:04:05.999Z"

type Article struct {
	tableName struct{} `pg:"articles,alias:a"`

	ID          uint64 `json:"-"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`

	Author   *org.Profile `json:"author"`
	AuthorID uint64       `json:"-"`

	Tags    []ArticleTag `json:"-"`
	TagList []string     `json:"tagList" pg:"-,array"`

	Favorited      bool `json:"favorited" pg:"-"`
	FavoritesCount int  `json:"favoritesCount" pg:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (a *Article) MarshalJSON() ([]byte, error) {
	type Alias Article

	if a.TagList == nil {
		a.TagList = make([]string, 0)
	}

	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		Alias:     (*Alias)(a),
		CreatedAt: a.CreatedAt.UTC().Format(TimeFormatStr),
		UpdatedAt: a.UpdatedAt.UTC().Format(TimeFormatStr),
	})
}

type ArticleTag struct {
	tableName struct{} `pg:"alias:t"`

	ArticleID uint64
	Tag       string
}

type FavoriteArticle struct {
	tableName struct{} `pg:"alias:fa"`

	UserID    uint64
	ArticleID uint64
}

func SelectArticle(c context.Context, slug string) (*Article, error) {
	article := new(Article)
	if err := rwe.PGMain().ModelContext(c, article).
		Where("slug = ?", slug).
		Select(); err != nil {
		return nil, err
	}
	return article, nil
}
