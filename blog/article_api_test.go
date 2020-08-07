package blog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

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
	ctx := context.Background()

	cfg, err := xconfig.LoadConfig("test")
	if err != nil {
		panic(err)
	}

	ctx = rwe.Init(ctx, cfg)
}

func assertArticle(article map[string]interface{}, keys Keys) {
	checkTime(article, "createdAt", time.Now())

	matchers := Keys{
		"description":    Equal("Ever wonder how?"),
		"body":           Equal("You have to believe"),
		"author":         Equal(map[string]interface{}{"following": false, "username": "hello", "bio": "", "image": ""}),
		"tagList":        ConsistOf([]interface{}{"reactjs", "angularjs", "dragons"}),
		"favoritesCount": Equal(float64(0)),
		"favorited":      Equal(false),
		"slug":           HaveSuffix("-how-to-train-your-dragon"),
		"title":          Equal("How to train your dragon"),
		"updatedAt":      Equal("0001-01-01T00:00:00Z"),
	}

	for key, value := range keys {
		matchers[key] = value
	}

	Expect(article).To(MatchAllKeys(matchers))
}

func checkTime(article map[string]interface{}, key string, expectedTime time.Time) {
	tm, err := time.Parse(time.RFC3339, article[key].(string))
	Expect(err).NotTo(HaveOccurred())

	ExpectWithOffset(1, tm).To(BeTemporally("~", expectedTime, time.Second))
	delete(article, key)
}

var _ = FDescribe("createArticle", func() {
	var resp map[string]interface{}
	var slug string
	var user *org.User

	var favoritedKey = Keys{
		"favorited":      Equal(true),
		"favoritesCount": Equal(float64(1)),
	}

	BeforeEach(func() {
		rwe.PGMain().Exec("TRUNCATE users;")
		rwe.PGMain().Exec("TRUNCATE articles;")
		rwe.PGMain().Exec("TRUNCATE article_tags;")

		user = &org.User{
			Username:     "hello",
			Email:        "foo@bar.com",
			PasswordHash: "hash",
		}

		_, err := rwe.PGMain().Model(user).Insert()
		Expect(err).NotTo(HaveOccurred())

		data := `{"title": "How to train your dragon", "description": "Ever wonder how?", "body": "You have to believe", "tagList": ["reactjs", "angularjs", "dragons"]}`
		req := NewReqWithToken("POST", "/api/articles", data, user.ID)

		ProcessReq(req, 200, &resp)

		slug = resp["article"].(map[string]interface{})["slug"].(string)
	})

	It("creates new article", func() {
		article := resp["article"].(map[string]interface{})
		assertArticle(article, nil)
	})

	Describe("showArticle", func() {
		BeforeEach(func() {
			url := fmt.Sprintf("/api/articles/%s", slug)
			req := NewReq("GET", url, "")
			ProcessReq(req, 200, &resp)
		})

		It("returns article", func() {
			article := resp["article"].(map[string]interface{})
			assertArticle(article, nil)
		})
	})

	Describe("favoriteArticle", func() {
		BeforeEach(func() {
			url := fmt.Sprintf("/api/articles/%s/favorite", slug)
			req := NewReqWithToken("POST", url, "", user.ID)
			ProcessReq(req, 200, &resp)

			url = fmt.Sprintf("/api/articles/%s", slug)
			req = NewReqWithToken("GET", url, "", user.ID)
			ProcessReq(req, 200, &resp)
		})

		It("returns favorited article", func() {
			article := resp["article"].(map[string]interface{})
			assertArticle(article, favoritedKey)
		})

		Describe("unfavoriteArticle", func() {
			BeforeEach(func() {
				url := fmt.Sprintf("/api/articles/%s/favorite", slug)
				req := NewReqWithToken("DELETE", url, "", user.ID)
				ProcessReq(req, 200, &resp)

				url = fmt.Sprintf("/api/articles/%s", slug)
				req = NewReqWithToken("GET", url, "", user.ID)
				ProcessReq(req, 200, &resp)
			})

			It("returns article", func() {
				article := resp["article"].(map[string]interface{})
				assertArticle(article, nil)
			})
		})
	})

	Describe("listArticles", func() {
		BeforeEach(func() {
			url := fmt.Sprintf("/api/articles/%s/favorite", slug)
			req := NewReqWithToken("POST", url, "", user.ID)
			ProcessReq(req, 200, &resp)

			req = NewReqWithToken("GET", "/api/articles", "", user.ID)
			ProcessReq(req, 200, &resp)
		})

		It("returns articles", func() {
			articles := resp["articles"].([]interface{})

			Expect(articles).To(HaveLen(1))
			article := articles[0].(map[string]interface{})

			assertArticle(article, favoritedKey)
		})
	})

	Describe("updateArticle", func() {
		BeforeEach(func() {
			data := `{"title": "Ice age", "description": "20,000 years before", "body": "All kinds of animals begin immigrating to the south", "tagList": ["drama", "comedy"]}`

			url := fmt.Sprintf("/api/articles/%s", slug)
			req := NewReqWithToken("PUT", url, data, user.ID)
			ProcessReq(req, 200, &resp)
		})

		It("returns article", func() {
			article := resp["article"].(map[string]interface{})

			checkTime(article, "createdAt", time.Now())
			checkTime(article, "updatedAt", time.Now())
			Expect(article).To(MatchAllKeys(Keys{
				"description":    Equal("20,000 years before"),
				"body":           Equal("All kinds of animals begin immigrating to the south"),
				"author":         Equal(map[string]interface{}{"following": false, "username": "hello", "bio": "", "image": ""}),
				"tagList":        ConsistOf([]interface{}{"drama", "comedy"}),
				"favoritesCount": Equal(float64(0)),
				"favorited":      Equal(false),
				"slug":           HaveSuffix("-ice-age"),
				"title":          Equal("Ice age"),
			}))
		})
	})
})
