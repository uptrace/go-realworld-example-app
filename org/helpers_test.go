package org_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/uptrace/go-realworld-example-app/org"

	"github.com/uptrace/go-realworld-example-app/rwe"

	. "github.com/onsi/gomega"
)

func newReqWithToken(method, url, data string, userID uint64) *http.Request {
	req := newReq(method, url, data)

	token, err := org.CreateUserToken(userID, time.Hour)
	Expect(err).NotTo(HaveOccurred())

	req.Header.Set("Authorization", "Token "+token)
	return req
}

func newReq(method, url, data string) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(data))
	Expect(err).NotTo(HaveOccurred())

	req.Header.Set("Content-Type", "application/json")

	return req
}

func processReq(req *http.Request, code int, v interface{}) {
	w := httptest.NewRecorder()
	rwe.Router.ServeHTTP(w, req)

	Expect(w.Code).To(Equal(code))

	err := json.Unmarshal(w.Body.Bytes(), v)
	Expect(err).NotTo(HaveOccurred())
}
