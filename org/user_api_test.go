package org_test

import (
	"context"
	"testing"

	"github.com/uptrace/go-realworld-example-app/rwe"
	"github.com/uptrace/go-realworld-example-app/xconfig"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
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

func assertUser(user map[string]interface{}) {
	Expect(user).To(MatchAllKeys(Keys{
		"username": Equal("wangzitian0"),
		"email":    Equal("wzt@gg.cn"),
		"bio":      Equal("bar"),
		"image":    Equal("img"),
		"token":    Not(BeEmpty()),
	}))
}

var _ = Describe("createUser", func() {
	var resp map[string]interface{}

	BeforeEach(func() {
		rwe.PGMain().Exec("TRUNCATE users;")

		data := `{"username": "wangzitian0","email": "wzt@gg.cn","password": "jakejxke", "image": "img", "bio": "bar"}`
		req := newReq("POST", "/api/users", data)

		processReq(req, 200, &resp)
	})

	It("creates new user", func() {
		assertUser(resp["user"].(map[string]interface{}))
	})

	Describe("loginUser", func() {
		var resp map[string]interface{}
		var token string

		BeforeEach(func() {
			data := `{"email": "wzt@gg.cn","password": "jakejxke"}`
			req := newReq("POST", "/api/users/login", data)

			processReq(req, 200, &resp)
			token = resp["user"].(map[string]interface{})["token"].(string)
		})

		It("returns user with JWT token", func() {
			assertUser(resp["user"].(map[string]interface{}))
		})

		Describe("currentUser", func() {
			var resp map[string]interface{}

			BeforeEach(func() {
				req := newReq("GET", "/api/user", "")
				req.Header.Set("Authorization", "Token "+token)
				processReq(req, 200, &resp)
			})

			It("returns logged in user", func() {
				assertUser(resp["user"].(map[string]interface{}))
			})
		})

		Describe("updateUser", func() {
			var resp map[string]interface{}

			BeforeEach(func() {
				data := `{"username": "hello","email": "foo@bar.com", "image": "bar", "bio": "foo"}`
				req := newReq("PUT", "/api/users", data)

				req.Header.Set("Authorization", "Token "+token)
				processReq(req, 200, &resp)
			})

			It("returns updated user", func() {
				user := resp["user"].(map[string]interface{})
				Expect(user).To(MatchAllKeys(Keys{
					"username": Equal("hello"),
					"email":    Equal("foo@bar.com"),
					"bio":      Equal("foo"),
					"image":    Equal("bar"),
					"token":    Not(BeEmpty()),
				}))
			})
		})
	})
})
