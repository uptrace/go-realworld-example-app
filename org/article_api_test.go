package org_test

import (
	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("createArticle", func() {
	var resp struct{ Article org.Article }

	BeforeEach(func() {
		rwe.PGMain().Exec("TRUNCATE articles;")
		rwe.PGMain().Exec("TRUNCATE article_tags;")

		data := `{"title": "How to train your dragon", "description": "Ever wonder how?", "body": "You have to believe", "tagList": ["reactjs", "angularjs", "dragons"]}`
		req := newReqWithToken("POST", "/api/articles", data, "token")

		processReq(req, 200, &resp)
	})

	FIt("creates new article", func() {
		Expect(resp.Article).To(Equal("wzt@gg.cn"))
	})
})
