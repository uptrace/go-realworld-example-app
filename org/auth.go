package org

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func authToken(req *http.Request) string {
	const prefix = "Token "
	v := req.Header.Get("Authorization")
	v = strings.TrimPrefix(v, prefix)
	return v
}

func AuthMiddleware(c *gin.Context) {
	token := authToken(c.Request)
	userID, err := decodeUserToken(token)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	user, err := SelectUser(c, userID)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	user.Token, err = CreateUserToken(user.ID, 24*time.Hour)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.Set("user", user)
}
