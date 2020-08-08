package testbed

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

func setToken(req *http.Request, userID uint64) {
	token, err := org.CreateUserToken(userID, time.Hour)
	Expect(err).NotTo(HaveOccurred())

	req.Header.Set("Authorization", "Token "+token)
}

func DoReqWithToken(method, url, data string, userID uint64) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(data))
	Expect(err).NotTo(HaveOccurred())

	req.Header.Set("Content-Type", "application/json")
	setToken(req, userID)

	return serve(req)
}

func DoReq(method, url, data string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(data))
	Expect(err).NotTo(HaveOccurred())

	req.Header.Set("Content-Type", "application/json")

	return serve(req)
}

func ParseJSON(resp *httptest.ResponseRecorder, code int) map[string]interface{} {
	res := make(map[string]interface{})
	err := json.Unmarshal(resp.Body.Bytes(), &res)
	Expect(err).NotTo(HaveOccurred())

	Expect(resp.Code).To(Equal(code))

	return res
}

func serve(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	rwe.Router.ServeHTTP(resp, req)
	return resp
}

func Get(url string) *httptest.ResponseRecorder {
	return DoReq("GET", url, "")
}

func GetWithToken(url string, userID uint64) *httptest.ResponseRecorder {
	return DoReqWithToken("GET", url, "", userID)
}

func PostWithToken(url, data string, userID uint64) *httptest.ResponseRecorder {
	return DoReqWithToken("POST", url, data, userID)
}

func PutWithToken(url, data string, userID uint64) *httptest.ResponseRecorder {
	return DoReqWithToken("PUT", url, data, userID)
}

func DeleteWithToken(url string, userID uint64) *httptest.ResponseRecorder {
	return DoReqWithToken("DELETE", url, "", userID)
}
