package org_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func newPOSTReq(url, data string, token string) *http.Request {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	Expect(err).NotTo(HaveOccurred())

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Token "+token)
	}

	return req
}

func processReq(req *http.Request, code int, v interface{}) {
	w := httptest.NewRecorder()
	rwe.Router.ServeHTTP(w, req)

	Expect(w.Code).To(Equal(code))

	err := json.Unmarshal(w.Body.Bytes(), v)
	Expect(err).NotTo(HaveOccurred())
}

var _ = Describe("createUser", func() {
	var resp struct{ User org.UserOut }

	BeforeEach(func() {
		rwe.PGMain().Exec("TRUNCATE users;")

		data := `{"username": "wangzitian0","email": "wzt@gg.cn","password": "jakejxke"}`
		req := newPOSTReq("/api/users", data, "")

		processReq(req, 200, &resp)
	})

	It("creates new user", func() {
		Expect(resp.User.Email).To(Equal("wzt@gg.cn"))
		Expect(resp.User.Username).To(Equal("wangzitian0"))
		Expect(resp.User.Bio).To(Equal(""))
		Expect(resp.User.Image).To(Equal(""))
		Expect(resp.User.Token).NotTo(BeEmpty())
	})

	Describe("loginUser", func() {
		var resp struct{ User org.UserOut }
		var token string

		BeforeEach(func() {
			data := `{"username": "wangzitian0","email": "wzt@gg.cn","password": "jakejxke"}`
			req := newPOSTReq("/api/users", data, "")

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
			var resp struct{ User org.UserOut }

			BeforeEach(func() {
				data := `{"username": "wangzitian0","email": "wzt@gg.cn","password": "jakejxke"}`
				req := newPOSTReq("/api/users", data, "Token "+token)

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
	})
})
