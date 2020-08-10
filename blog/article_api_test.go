package blog_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"
	. "github.com/uptrace/go-realworld-example-app/testbed"
	"github.com/uptrace/go-realworld-example-app/xconfig"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "blog")
}

func init() {
	rwe.Clock = clock.NewMock()
	ctx := context.Background()

	cfg, err := xconfig.LoadConfig("test")
	if err != nil {
		panic(err)
	}

	ctx = rwe.Init(ctx, cfg)
}

var _ = Describe("createArticle", func() {
	var data map[string]interface{}
	var slug string
	var user *org.User

	var articleKeys, favoritedArticleKeys Keys

	BeforeEach(func() {
		_, err := rwe.PGMain().Exec("TRUNCATE users, follow_users;")
		Expect(err).NotTo(HaveOccurred())

		_, err = rwe.PGMain().Exec("TRUNCATE articles, article_tags, favorite_articles;")
		Expect(err).NotTo(HaveOccurred())

		articleKeys = Keys{
			"description":    Equal("Ever wonder how?"),
			"body":           Equal("You have to believe"),
			"author":         Equal(map[string]interface{}{"following": false, "username": "hello", "bio": "", "image": ""}),
			"tagList":        ConsistOf([]interface{}{"reactjs", "angularjs", "dragons"}),
			"favoritesCount": Equal(float64(0)),
			"favorited":      Equal(false),
			"slug":           HaveSuffix("-how-to-train-your-dragon"),
			"title":          Equal("How to train your dragon"),
			"createdAt":      Equal(rwe.Clock.Now().Format(time.RFC3339)),
			"updatedAt":      Equal("0001-01-01T00:00:00Z"),
		}

		favoritedArticleKeys = extend(articleKeys, Keys{
			"favorited":      Equal(true),
			"favoritesCount": Equal(float64(1)),
		})

		user = &org.User{
			Username:     "hello",
			Email:        "foo@bar.com",
			PasswordHash: "hash",
		}
		_, err = rwe.PGMain().Model(user).Insert()
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		json := `{"title": "How to train your dragon", "description": "Ever wonder how?", "body": "You have to believe", "tagList": ["reactjs", "angularjs", "dragons"]}`
		resp := PostWithToken("/api/articles", json, user.ID)

		data = ParseJSON(resp, http.StatusOK)
		slug = data["article"].(map[string]interface{})["slug"].(string)
	})

	It("creates new article", func() {
		Expect(data["article"]).To(MatchAllKeys(articleKeys))
	})

	Describe("showArticle", func() {
		BeforeEach(func() {
			url := fmt.Sprintf("/api/articles/%s", slug)
			resp := Get(url)

			data = ParseJSON(resp, http.StatusOK)
		})

		It("returns article", func() {
			Expect(data["article"]).To(MatchAllKeys(articleKeys))
		})
	})

	Describe("favoriteArticle", func() {
		BeforeEach(func() {
			url := fmt.Sprintf("/api/articles/%s/favorite", slug)
			resp := PostWithToken(url, "", user.ID)
			_ = ParseJSON(resp, 200)

			url = fmt.Sprintf("/api/articles/%s", slug)
			resp = GetWithToken(url, user.ID)
			data = ParseJSON(resp, 200)
		})

		It("returns favorited article", func() {
			Expect(data["article"]).To(MatchAllKeys(favoritedArticleKeys))
		})

		Describe("unfavoriteArticle", func() {
			BeforeEach(func() {
				url := fmt.Sprintf("/api/articles/%s/favorite", slug)
				resp := DeleteWithToken(url, user.ID)
				_ = ParseJSON(resp, 200)

				url = fmt.Sprintf("/api/articles/%s", slug)
				resp = GetWithToken(url, user.ID)
				data = ParseJSON(resp, 200)
			})

			It("returns article", func() {
				Expect(data["article"]).To(MatchAllKeys(articleKeys))
			})
		})
	})

	Describe("listArticles", func() {
		BeforeEach(func() {
			url := fmt.Sprintf("/api/articles/%s/favorite", slug)
			resp := PostWithToken(url, "", user.ID)
			_ = ParseJSON(resp, 200)

			resp = GetWithToken("/api/articles", user.ID)
			data = ParseJSON(resp, 200)
		})

		It("returns articles", func() {
			articles := data["articles"].([]interface{})

			Expect(articles).To(HaveLen(1))
			article := articles[0].(map[string]interface{})
			Expect(article).To(MatchAllKeys(favoritedArticleKeys))
		})
	})

	Describe("updateArticle", func() {
		BeforeEach(func() {
			json := `{"title": "Ice age", "description": "20,000 years before", "body": "All kinds of animals begin immigrating to the south", "tagList": ["drama", "comedy"]}`

			url := fmt.Sprintf("/api/articles/%s", slug)
			resp := PutWithToken(url, json, user.ID)
			data = ParseJSON(resp, 200)
		})

		It("returns article", func() {
			Expect(data["article"]).To(MatchAllKeys(Keys{
				"description":    Equal("20,000 years before"),
				"body":           Equal("All kinds of animals begin immigrating to the south"),
				"author":         Equal(map[string]interface{}{"following": false, "username": "hello", "bio": "", "image": ""}),
				"tagList":        ConsistOf([]interface{}{"drama", "comedy"}),
				"favoritesCount": Equal(float64(0)),
				"favorited":      Equal(false),
				"slug":           HaveSuffix("-ice-age"),
				"title":          Equal("Ice age"),
				"createdAt":      Equal(rwe.Clock.Now().Format(time.RFC3339)),
				"updatedAt":      Equal(rwe.Clock.Now().Format(time.RFC3339)),
			}))
		})
	})
})

func extend(a, b Keys) Keys {
	res := make(Keys)
	for k, v := range a {
		res[k] = v
	}
	for k, v := range b {
		res[k] = v
	}
	return res
}
