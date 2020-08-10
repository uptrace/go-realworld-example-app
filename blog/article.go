package blog

import (
	"time"

	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"

	"github.com/gin-gonic/gin"
)

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

func SelectArticle(c *gin.Context, slug string) (*Article, error) {
	article := new(Article)
	if err := rwe.PGMain().ModelContext(c, article).
		Where("slug = ?", slug).
		Select(); err != nil {
		return nil, err
	}
	return article, nil
}
