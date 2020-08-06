package org_test

import (
	"fmt"
	"time"

	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

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

var _ = FDescribe("createArticle", func() {
	var resp map[string]interface{}
	var slug string

	var checkTime = func(article map[string]interface{}, key string, expectedTime time.Time) {
		tm, err := time.Parse(time.RFC3339, article[key].(string))
		Expect(err).NotTo(HaveOccurred())

		Expect(tm).To(BeTemporally("~", expectedTime, time.Second))
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

		req := newReqWithToken("POST", "/api/articles", data, user.ID)

		processReq(req, 200, &resp)

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
			req := newReq("GET", url, "")
			processReq(req, 200, &resp)
		})

		It("returns article", func() {
			article := resp["article"].(map[string]interface{})
			checkTime(article, "createdAt", time.Now())

			assertArticle(article)
		})
	})

	Describe("listArticles", func() {
		BeforeEach(func() {
			req := newReq("GET", "/api/articles", "")
			processReq(req, 200, &resp)
		})

		FIt("returns articles", func() {
			articles := resp["articles"].([]interface{})
			article := articles[0].(map[string]interface{})
			checkTime(article, "createdAt", time.Now())

			assertArticle(article)
		})
	})
})
