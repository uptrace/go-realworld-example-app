package org_test

import (
	"time"

	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("createArticle", func() {
	var resp map[string]interface{}

	var checkTime = func(article map[string]interface{}, key string, expectedTime time.Time) {
		tm, err := time.Parse(time.RFC3339, article[key].(string))
		Expect(err).NotTo(HaveOccurred())

		Expect(tm).To(BeTemporally("~", expectedTime, time.Second))
		delete(article, key)
	}

	BeforeEach(func() {
		user := &org.User{
			Username:     "hello",
			Email:        "foo@bar.com",
			PasswordHash: "hash",
		}

		_, err := rwe.PGMain().Model(user).Insert()
		Expect(err).NotTo(HaveOccurred())

		rwe.PGMain().Exec("TRUNCATE articles;")
		rwe.PGMain().Exec("TRUNCATE article_tags;")

		data := `{"slug": "how-to-train-your-dragon", "title": "How to train your dragon", "description": "Ever wonder how?", "body": "You have to believe", "tagList": ["reactjs", "angularjs", "dragons"]}`

		req := newReqWithToken("POST", "/api/articles", data, user.ID)

		processReq(req, 200, &resp)
	})

	It("creates new article", func() {
		article := resp["article"].(map[string]interface{})
		checkTime(article, "createdAt", time.Now())

		Expect(article).To(MatchAllKeys(Keys{
			"description":    Equal("Ever wonder how?"),
			"body":           Equal("You have to believe"),
			"author":         Equal(map[string]interface{}{"following": false, "username": "hello", "bio": "", "image": ""}),
			"tagList":        Equal([]interface{}{"reactjs", "angularjs", "dragons"}),
			"favoritesCount": Equal(float64(0)),
			"favorited":      Equal(false),
			"slug":           Equal("how-to-train-your-dragon"),
			"title":          Equal("How to train your dragon"),
			"updatedAt":      Equal("0001-01-01T00:00:00Z"),
		}))
	})
})
