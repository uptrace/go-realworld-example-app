package httperror

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	json []byte

	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorJSON Error

func (e *Error) MarshalJSON() ([]byte, error) {
	if len(e.json) != 0 {
		return e.json, nil
	}

	b, err := json.Marshal((*errorJSON)(e))
	if err != nil {
		return nil, err
	}
	e.json = b

	return b, nil
}

func (e *Error) Error() string {
	return e.Message
}

func New(status int, code, msg string, args ...interface{}) *Error {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return &Error{
		Status:  status,
		Code:    code,
		Message: msg,
	}
}

func NotFound(msg string, args ...interface{}) *Error {
	return New(http.StatusNotFound, "not_found", msg, args...)
}

func Unauthorized(msg string, args ...interface{}) *Error {
	return New(http.StatusUnauthorized, "unauthorized", msg, args...)
}

func Forbidden(msg string, args ...interface{}) *Error {
	return New(http.StatusForbidden, "forbidden", msg, args...)
}

func BadRequest(code, msg string, args ...interface{}) *Error {
	return New(http.StatusBadRequest, code, msg, args...)
}
