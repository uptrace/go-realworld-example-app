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
	. "github.com/uptrace/go-realworld-example-app/testhelper"
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

var _ = FDescribe("createArticle", func() {
	var data map[string]interface{}
	var slug string
	var user *org.User

	var helloArticleKeys, fooArticleKeys, favoritedArticleKeys Keys

	BeforeEach(func() {
		TruncateUsersTable()
		TruncateArticlesTable()

		helloArticleKeys = Keys{
			"title":          Equal("Hello world"),
			"slug":           HaveSuffix("-hello-world"),
			"description":    Equal("Hello world article description!"),
			"body":           Equal("Hello world article body."),
			"author":         Equal(map[string]interface{}{"following": false, "username": "CurrentUser", "bio": "", "image": ""}),
			"tagList":        ConsistOf([]interface{}{"greeting", "welcome", "salut"}),
			"favoritesCount": Equal(float64(0)),
			"favorited":      Equal(false),
			"createdAt":      Equal(rwe.Clock.Now().Format(time.RFC3339)),
			"updatedAt":      Equal("0001-01-01T00:00:00Z"),
		}

		favoritedArticleKeys = Extend(helloArticleKeys, Keys{
			"favorited":      Equal(true),
			"favoritesCount": Equal(float64(1)),
		})

		fooArticleKeys = Keys{
			"title":          Equal("Foo bar"),
			"slug":           HaveSuffix("-foo-bar"),
			"description":    Equal("Foo bar article description!"),
			"body":           Equal("Foo bar article body."),
			"author":         Equal(map[string]interface{}{"following": false, "username": "CurrentUser", "bio": "", "image": ""}),
			"tagList":        ConsistOf([]interface{}{"foobar", "variable"}),
			"favoritesCount": Equal(float64(0)),
			"favorited":      Equal(false),
			"createdAt":      Equal(rwe.Clock.Now().Format(time.RFC3339)),
			"updatedAt":      Equal("0001-01-01T00:00:00Z"),
		}

		user = &org.User{
			Username:     "CurrentUser",
			Email:        "hello@world.com",
			PasswordHash: "#1",
		}
		_, err := rwe.PGMain().Model(user).Insert()
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		json := `{"title": "Hello world", "description": "Hello world article description!", "body": "Hello world article body.", "tagList": ["greeting", "welcome", "salut"]}`
		resp := PostWithToken("/api/articles", json, user.ID)

		data = ParseJSON(resp, http.StatusOK)
		slug = data["article"].(map[string]interface{})["slug"].(string)
	})

	It("creates new article", func() {
		Expect(data["article"]).To(MatchAllKeys(helloArticleKeys))
	})

	Describe("showFeed", func() {
		BeforeEach(func() {
			followedUser := &org.User{
				Username:     "FollowedUser",
				Email:        "foo@bar.com",
				PasswordHash: "h2",
			}
			_, err := rwe.PGMain().Model(followedUser).Insert()
			Expect(err).NotTo(HaveOccurred())

			url := fmt.Sprintf("/api/profiles/%s/follow", followedUser.Username)
			resp := PostWithToken(url, "", user.ID)
			_ = ParseJSON(resp, 200)

			json := `{"title": "Foo bar", "description": "Foo bar article description!", "body": "Foo bar article body.", "tagList": ["foobar", "variable"]}`
			resp = PostWithToken("/api/articles", json, followedUser.ID)

			_ = ParseJSON(resp, http.StatusOK)

			resp = GetWithToken("/api/articles/feed", user.ID)
			data = ParseJSON(resp, http.StatusOK)
		})

		It("returns article", func() {
			articles := data["articles"].([]interface{})

			Expect(articles).To(HaveLen(1))
			followedAuthorKeys := Extend(fooArticleKeys, Keys{
				"author": Equal(map[string]interface{}{"following": true, "username": "FollowedUser", "bio": "", "image": ""}),
			})
			Expect(articles[0].(map[string]interface{})).To(MatchAllKeys(followedAuthorKeys))
		})
	})

	Describe("showArticle", func() {
		BeforeEach(func() {
			url := fmt.Sprintf("/api/articles/%s", slug)
			resp := Get(url)

			data = ParseJSON(resp, http.StatusOK)
		})

		It("returns article", func() {
			Expect(data["article"]).To(MatchAllKeys(helloArticleKeys))
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
				Expect(data["article"]).To(MatchAllKeys(helloArticleKeys))
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
			json := `{"title": "Foo bar", "description": "Foo bar article description!", "body": "Foo bar article body.", "tagList": ["foobar", "variable"]}`

			url := fmt.Sprintf("/api/articles/%s", slug)
			resp := PutWithToken(url, json, user.ID)
			data = ParseJSON(resp, 200)
		})

		It("returns article", func() {
			updatedArticleKeys := Extend(fooArticleKeys, Keys{
				"updatedAt": Equal(rwe.Clock.Now().Format(time.RFC3339)),
			})
			Expect(data["article"]).To(MatchAllKeys(updatedArticleKeys))
		})
	})

	Describe("deleteArticle", func() {
		BeforeEach(func() {
			url := fmt.Sprintf("/api/articles/%s", slug)
			resp := DeleteWithToken(url, user.ID)
			data = ParseJSON(resp, 200)
		})

		It("returns article", func() {
			Expect(data).To(BeNil())
		})
	})
})
