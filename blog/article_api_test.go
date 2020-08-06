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

func assertArticle(article map[string]interface{}) {
	Expect(article).To(MatchAllKeys(Keys{
		"description":    Equal("Ever wonder how?"),
		"body":           Equal("You have to believe"),
		"author":         Equal(map[string]interface{}{"following": false, "username": "hello", "bio": "", "image": ""}),
		"tagList":        ConsistOf([]interface{}{"reactjs", "angularjs", "dragons"}),
		"favoritesCount": Equal(float64(0)),
		"favorited":      Equal(false),
		"slug":           HavePrefix("how-to-train-your-dragon-"),
		"title":          Equal("How to train your dragon"),
		"updatedAt":      Equal("0001-01-01T00:00:00Z"),
	}))
}

var _ = Describe("createArticle", func() {
	var resp map[string]interface{}
	var slug string

	checkTime := func(article map[string]interface{}, key string, expectedTime time.Time) {
		tm, err := time.Parse(time.RFC3339, article[key].(string))
		Expect(err).NotTo(HaveOccurred())

		ExpectWithOffset(1, tm).To(BeTemporally("~", expectedTime, time.Second))
		delete(article, key)
	}

	BeforeEach(func() {
		rwe.PGMain().Exec("TRUNCATE users;")
		rwe.PGMain().Exec("TRUNCATE articles;")
		rwe.PGMain().Exec("TRUNCATE article_tags;")

		user := &org.User{
			Username:     "hello",
			Email:        "foo@bar.com",
			PasswordHash: "hash",
		}

		_, err := rwe.PGMain().Model(user).Insert()
		Expect(err).NotTo(HaveOccurred())

		data := `{"slug": "how-to-train-your-dragon", "title": "How to train your dragon", "description": "Ever wonder how?", "body": "You have to believe", "tagList": ["reactjs", "angularjs", "dragons"]}`

		req := NewReqWithToken("POST", "/api/articles", data, user.ID)

		ProcessReq(req, 200, &resp)

		slug = resp["article"].(map[string]interface{})["slug"].(string)
	})

	It("creates new article", func() {
		article := resp["article"].(map[string]interface{})
		checkTime(article, "createdAt", time.Now())

		assertArticle(article)

		Expect(article).To(MatchAllKeys(Keys{
			"description":    Equal("Ever wonder how?"),
			"body":           Equal("You have to believe"),
			"author":         Equal(map[string]interface{}{"following": false, "username": "hello", "bio": "", "image": ""}),
			"tagList":        Equal([]interface{}{"reactjs", "angularjs", "dragons"}),
			"favoritesCount": Equal(float64(0)),
			"favorited":      Equal(false),
			"slug":           HavePrefix("how-to-train-your-dragon-"),
			"title":          Equal("How to train your dragon"),
			"updatedAt":      Equal("0001-01-01T00:00:00Z"),
		}))
	})

	Describe("showArticle", func() {
		BeforeEach(func() {
			url := fmt.Sprintf("/api/articles/%s", slug)
			req := NewReq("GET", url, "")
			ProcessReq(req, 200, &resp)
		})

		It("returns article", func() {
			article := resp["article"].(map[string]interface{})
			checkTime(article, "createdAt", time.Now())

			assertArticle(article)
		})
	})

	Describe("listArticles", func() {
		BeforeEach(func() {
			req := NewReq("GET", "/api/articles", "")
			ProcessReq(req, 200, &resp)
		})

		It("returns articles", func() {
			articles := resp["articles"].([]interface{})
			article := articles[0].(map[string]interface{})
			checkTime(article, "createdAt", time.Now())

			assertArticle(article)
		})
	})
})
