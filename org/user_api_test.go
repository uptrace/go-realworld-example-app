package org_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"
	"github.com/uptrace/go-realworld-example-app/xconfig"
)

func TestOrg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "org")
}

func init() {
	ctx := context.Background()

	cfg, err := xconfig.LoadConfig("test")
	if err != nil {
		panic(err)
	}

	ctx = rwe.Init(ctx, cfg)
}

var _ = Describe("createUser", func() {
	var resp struct{ User org.User }

	BeforeEach(func() {
		rwe.PGMain().Exec("TRUNCATE users;")

		data := `{"username": "wangzitian0","email": "wzt@gg.cn","password": "jakejxke", "image": "img"}`
		req := newReq("POST", "/api/users", data)

		processReq(req, 200, &resp)
	})

	It("creates new user", func() {
		Expect(resp.User.Email).To(Equal("wzt@gg.cn"))
		Expect(resp.User.Username).To(Equal("wangzitian0"))
		Expect(resp.User.Bio).To(Equal(""))
		Expect(resp.User.Image).To(Equal("img"))
		Expect(resp.User.Token).NotTo(BeEmpty())
	})

	Describe("loginUser", func() {
		var resp struct{ User org.User }
		var token string

		BeforeEach(func() {
			data := `{"username": "wangzitian0","email": "wzt@gg.cn","password": "jakejxke"}`
			req := newReq("POST", "/api/users", data)

			processReq(req, 200, &resp)
			token = resp.User.Token
		})

		It("returns user with JWT token", func() {
			Expect(resp.User.Email).To(Equal("wzt@gg.cn"))
			Expect(resp.User.Username).To(Equal("wangzitian0"))
			Expect(resp.User.Bio).To(Equal(""))
			Expect(resp.User.Image).To(Equal(""))
			Expect(resp.User.Token).NotTo(BeEmpty())
		})

		Describe("currentUser", func() {
			var resp struct{ User org.User }

			BeforeEach(func() {
				req := newReqWithToken("GET", "/api/user", "", token)
				processReq(req, 200, &resp)
			})

			It("returns logged in user", func() {
				Expect(resp.User.Email).To(Equal("wzt@gg.cn"))
				Expect(resp.User.Username).To(Equal("wangzitian0"))
				Expect(resp.User.Bio).To(Equal(""))
				Expect(resp.User.Image).To(Equal(""))
				Expect(resp.User.Token).NotTo(BeEmpty())
			})
		})

		Describe("updateUser", func() {
			var resp struct{ User org.User }

			BeforeEach(func() {
				token := "fix"
				data := `{"username": "hello","email": "foo@bar.com"}`
				req := newReqWithToken("PUT", "/api/users", data, token)

				processReq(req, 200, &resp)
			})

			FIt("returns updated user", func() {
				Expect(resp.User.Email).To(Equal("foo@bar.com"))
				Expect(resp.User.Username).To(Equal("hello"))
				Expect(resp.User.Bio).To(Equal(""))
				Expect(resp.User.Image).To(Equal(""))
				Expect(resp.User.Token).NotTo(BeEmpty())
			})
		})
	})
})
