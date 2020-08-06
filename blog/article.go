package blog

import (
	"time"
)

type Article struct {
	tableName struct{} `pg:"articles,alias:a"`

	ID          uint64 `json:"-"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`

	Author   *Author `json:"author"`
	AuthorID uint64  `json:"-"`

	Tags    []ArticleTag `json:"-"`
	TagList []string     `json:"tagList" pg:"-,array"`

	Favorited      bool `json:"favorited" pg:"-"`
	FavoritesCount int  `json:"favoritesCount" pg:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Author struct {
	tableName struct{} `pg:"users,alias:u"`

	ID        uint64 `json:"-"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following" pg:"-"`
}

type ArticleTag struct {
	tableName struct{} `pg:"alias:t"`

	ArticleID uint64
	Tag       string
}
