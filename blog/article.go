package blog

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10/orm"
	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"
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

type FavoriteArticle struct {
	tableName struct{} `pg:"alias:fa"`

	UserID    uint64
	ArticleID uint64
}

func articleTagsSubquery(q *orm.Query) (*orm.Query, error) {
	subq := rwe.PGMain().Model((*ArticleTag)(nil)).
		ColumnExpr("array_agg(t.tag)::text[]").
		Where("t.article_id = a.id")

	return q.ColumnExpr("(?) AS tag_list", subq), nil
}

func articleFavoritedSubquery(c *gin.Context) func(*orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		user, ok := c.Get("user")
		if !ok {
			return q.ColumnExpr("false AS favorited"), nil
		}

		subq := rwe.PGMain().Model((*FavoriteArticle)(nil)).
			Where("fa.article_id = a.id").
			Where("fa.user_id = ?", user.(*org.User).ID)

		return q.ColumnExpr("EXISTS(?) AS favorited", subq), nil
	}
}

func articleFavoritesCountSubquery(q *orm.Query) (*orm.Query, error) {
	subq := rwe.PGMain().Model((*FavoriteArticle)(nil)).
		ColumnExpr("count(*)").
		Where("fa.article_id = a.id")

	return q.ColumnExpr("(?) AS favorites_count", subq), nil
}

func SelectArticle(c *gin.Context, slug string) (*Article, error) {
	article := new(Article)
	if err := rwe.PGMain().
		ModelContext(c, article).
		ColumnExpr("?TableColumns").
		Relation("Author").
		Apply(articleTagsSubquery).
		Apply(articleFavoritedSubquery(c)).
		Apply(articleFavoritesCountSubquery).
		Where("slug = ?", slug).
		Select(); err != nil {
		return nil, err
	}

	return article, nil
}

func SelectArticles(c *gin.Context, f *ArticleFilter) ([]*Article, error) {
	articles := make([]*Article, 0)
	err := rwe.PGMain().ModelContext(c, &articles).
		ColumnExpr("?TableColumns").
		Apply(articleTagsSubquery).
		Apply(articleFavoritedSubquery(c)).
		Apply(articleFavoritesCountSubquery).
		Apply(f.Filters).
		Relation("Author").
		Limit(f.Pager.GetLimit()).
		Offset(f.Pager.GetOffset()).
		Select()
	if err != nil {
		return nil, err
	}

	return articles, nil
}
