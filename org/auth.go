package org

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func authToken(req *http.Request) string {
	const prefix = "Token "
	v := req.Header.Get("Authorization")
	if strings.HasPrefix(v, prefix) {
		return v[len(prefix):]
	}

	return ""
}

func AuthMiddleware(c *gin.Context) {
	token := authToken(c.Request)
	userID, err := decodeUserToken(token)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	user, err := SelectUser(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.Set("user", user)
}
