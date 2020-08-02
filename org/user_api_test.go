package org_test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/uptrace/go-realworld-example-app/rwe"
	"github.com/uptrace/go-realworld-example-app/xconfig"
)

var listenFlag = flag.String("listen", ":8888", "listen address")

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
	var w *httptest.ResponseRecorder

	BeforeEach(func() {
		rwe.PGMain().Exec("TRUNCATE users;")

		data := `{"username": "wangzitian0","email": "wzt@gg.cn","password": "jakejxke"}`
		req, err := http.NewRequest("POST", "/api/users", bytes.NewBufferString(data))
		Expect(err).NotTo(HaveOccurred())

		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		rwe.Router.ServeHTTP(w, req)
	})

	It("create new user", func() {
		Expect(w.Code).To(Equal(200))
		Expect(w.Body.String()).To(ContainSubstring(`"username":"wangzitian0","email":"wzt@gg.cn","bio":"","password":"jakejxke"}`))
	})

	Describe("loginUser", func() {
		var resp struct {
			User struct {
				Email string `json:"email"`
				Token string `json:"token"`
			} `json:"user"`
		}

		BeforeEach(func() {
			data := `{"email": "wzt@gg.cn","password": "jakejxke"}`
			req, err := http.NewRequest("POST", "/api/users/login", bytes.NewBufferString(data))
			Expect(err).NotTo(HaveOccurred())

			req.Header.Set("Content-Type", "application/json")

			w = httptest.NewRecorder()
			rwe.Router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(200))

			err = json.Unmarshal(w.Body.Bytes(), &resp)
			Expect(err).NotTo(HaveOccurred())
		})

		It("login user", func() {
			Expect(resp.User.Email).To(Equal("wzt@gg.cn"))
			Expect(resp.User.Token).NotTo(BeEmpty())
		})

		Describe("currentUser", func() {
			BeforeEach(func() {
				req, err := http.NewRequest("POST", "/api/users/current", nil)
				Expect(err).NotTo(HaveOccurred())

				var bearer = "TOKEN " + resp.User.Token
				req.Header.Set("Authorization", bearer)
				req.Header.Set("Content-Type", "application/json")

				w = httptest.NewRecorder()
				rwe.Router.ServeHTTP(w, req)
			})

			It("returns logined user", func() {
				Expect(w.Code).To(Equal(200))
				Expect(w.Body.String()).To(ContainSubstring(`"username":"wangzitian0"`))
			})
		})
	})
})
