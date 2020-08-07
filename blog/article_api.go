package blog

import (
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"
)

const charsBytes = "01234567890"

func randBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charsBytes[rand.Intn(len(charsBytes))]
	}
	return string(b)
}

func newSlug(title string) string {
	return slug.Make(randBytes(6) + " " + title)
}

func listArticles(c *gin.Context) {
	f, err := decodeArticleFilter(c)
	if err != nil {
		c.Error(err)
		return
	}

	articles := make([]*Article, 0)
	err := rwe.PGMain().ModelContext(c, &articles).
		ColumnExpr("?TableColumns").
		Apply(f.query).
		Limit(f.Pager.GetLimit()).
		Offset(f.Pager.GetOffset()).
		Select()
	if err != nil {
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
		Where("slug = ?", slug).
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

	article.Slug = newSlug(article.Title)
	article.AuthorID = user.ID
	article.CreatedAt = rwe.Clock.Now()

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

func updateArticle(c *gin.Context) {
	user := c.MustGet("user").(*org.User)

	article, err := SelectArticle(c, c.Param("slug"))
	if err != nil {
		c.Error(err)
		return
	}

	newArticle := new(Article)
	if err := c.BindJSON(newArticle); err != nil {
		return
	}

	newArticle.Slug = slug.Make(randBytes(6) + " " + article.Title)
	newArticle.AuthorID = user.ID

	if _, err := rwe.PGMain().
		ModelContext(c, newArticle).
		Set("title = ?", newArticle.Title).
		Set("slug = ?", newSlug(newArticle.Title)).
		Set("description = ?", newArticle.Description).
		Set("body = ?", newArticle.Body).
		Set("updated_at = ?", rwe.Clock.Now()).
		Where("id = ?", article.ID).
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

	newArticle.Author = &Author{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: false,
	}
	c.JSON(200, gin.H{"article": newArticle})
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

	c.JSON(200, gin.H{"profile": article})
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
		Where("user_id = ?", article.ID).
		Where("article_id = ?", user.ID).
		Delete(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"profile": article})
}
