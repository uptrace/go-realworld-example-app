package org_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/uptrace/go-realworld-example-app/org"
	"github.com/uptrace/go-realworld-example-app/rwe"
	. "github.com/uptrace/go-realworld-example-app/testbed"
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
		"username":  Equal("wangzitian0"),
		"email":     Equal("wzt@gg.cn"),
		"bio":       Equal("bar"),
		"image":     Equal("img"),
		"token":     Not(BeEmpty()),
		"following": Equal(false),
	}))
}

var _ = Describe("createUser", func() {
	var data map[string]interface{}

	BeforeEach(func() {
		_, err := rwe.PGMain().Exec("TRUNCATE follow_users, users;")
		Expect(err).NotTo(HaveOccurred())

		json := `{"username": "wangzitian0","email": "wzt@gg.cn","password": "jakejxke", "image": "img", "bio": "bar"}`
		resp := Post("/api/users", json)

		data = ParseJSON(resp, http.StatusOK)
	})

	It("creates new user", func() {
		assertUser(data["user"].(map[string]interface{}))
	})

	Describe("loginUser", func() {
		var user *org.User

		BeforeEach(func() {
			json := `{"email": "wzt@gg.cn","password": "jakejxke"}`
			resp := Post("/api/users/login", json)

			data = ParseJSON(resp, http.StatusOK)

			username := data["user"].(map[string]interface{})["username"].(string)
			var err error
			user, err = org.SelectUserByUsername(context.Background(), username)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns user with JWT token", func() {
			assertUser(data["user"].(map[string]interface{}))
		})

		Describe("currentUser", func() {
			BeforeEach(func() {
				resp := GetWithToken("/api/user", user.ID)
				data = ParseJSON(resp, http.StatusOK)
			})

			It("returns logged in user", func() {
				assertUser(data["user"].(map[string]interface{}))
			})
		})

		Describe("updateUser", func() {
			BeforeEach(func() {
				json := `{"username": "hello","email": "foo@bar.com", "image": "bar", "bio": "foo"}`
				resp := PutWithToken("/api/users", json, user.ID)
				data = ParseJSON(resp, http.StatusOK)
			})

			It("returns updated user", func() {
				user := data["user"].(map[string]interface{})
				Expect(user).To(MatchAllKeys(Keys{
					"username":  Equal("hello"),
					"email":     Equal("foo@bar.com"),
					"bio":       Equal("foo"),
					"image":     Equal("bar"),
					"token":     Not(BeEmpty()),
					"following": Equal(false),
				}))
			})
		})

		Describe("followUser", func() {
			var username string

			BeforeEach(func() {
				json := `{"username": "hello","email": "foo@bar.com","password": "pwd"}`
				resp := Post("/api/users", json)

				data = ParseJSON(resp, http.StatusOK)

				username = data["user"].(map[string]interface{})["username"].(string)

				url := fmt.Sprintf("/api/profiles/%s/follow", username)
				resp = PostWithToken(url, "", user.ID)
				data = ParseJSON(resp, 200)
			})

			It("returns followed profile", func() {
				profile := data["profile"].(map[string]interface{})
				Expect(profile).To(MatchAllKeys(Keys{
					"username":  Equal("hello"),
					"bio":       Equal(""),
					"image":     Equal(""),
					"following": Equal(true),
				}))
			})

			Describe("unfollowUser", func() {
				BeforeEach(func() {
					url := fmt.Sprintf("/api/profiles/%s/follow", username)
					resp := DeleteWithToken(url, user.ID)
					data = ParseJSON(resp, 200)
				})

				It("returns profile", func() {
					profile := data["profile"].(map[string]interface{})
					Expect(profile).To(MatchAllKeys(Keys{
						"username":  Equal("hello"),
						"bio":       Equal(""),
						"image":     Equal(""),
						"following": Equal(false),
					}))
				})
			})
		})
	})
})
