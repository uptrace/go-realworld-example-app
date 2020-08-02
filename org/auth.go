package org

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/uptrace/go-realworld-example-app/rwe"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

func stripBearerPrefixFromTokenString(tok string) (string, error) {
	if len(tok) > 5 && strings.ToUpper(tok[0:6]) == "TOKEN " {
		return tok[6:], nil
	}
	return tok, nil
}

var AuthorizationHeaderExtractor = &request.PostExtractionFilter{
	request.HeaderExtractor{"Authorization"},
	stripBearerPrefixFromTokenString,
}

var MyAuth2Extractor = &request.MultiExtractor{
	AuthorizationHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

func UpdateContextUserModel(c *gin.Context, userID uint64) {
	user, err := SelectUser(userID)
	if err != nil && err != pg.ErrNoRows {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.Set("userID", userID)
	c.Set("user", user)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		UpdateContextUserModel(c, 0)

		token, err := request.ParseFromRequest(c.Request, MyAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			b := ([]byte(rwe.Config.SecretKey))
			return b, nil
		})
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := uint64(claims["id"].(float64))
			UpdateContextUserModel(c, userID)
		}
	}
}

func newToken(id uint64) string {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))
	jwt_token.Claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token, _ := jwt_token.SignedString([]byte(rwe.Config.SecretKey))
	return token
}
