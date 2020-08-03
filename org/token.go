package org

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/uptrace/go-realworld-example-app/rwe"
)

// var errTokenEmpty = errors.New("token is missing or empty")

func decodeUserToken(jwtToken string) (uint64, error) {
	if len(jwtToken) == 0 {
		return 0, errors.New("token is missing or empty")
	}

	token, err := jwt.ParseWithClaims(jwtToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(rwe.Config.SecretKey), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims := token.Claims.(*jwt.StandardClaims)

	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func createUserToken(userID uint64, ttl time.Duration) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   strconv.FormatUint(userID, 10),
		ExpiresAt: time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key := []byte(rwe.Config.SecretKey)
	return token.SignedString(key)
}

// // -----------------------
// func newToken(id uint64) string {
// 	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))
// 	jwt_token.Claims = jwt.MapClaims{
// 		"id":  id,
// 		"exp": time.Now().Add(time.Hour * 24).Unix(),
// 	}

// 	token, _ := jwt_token.SignedString([]byte(rwe.Config.SecretKey))
// 	return token
// }
