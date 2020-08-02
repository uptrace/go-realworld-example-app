package org

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"golang.org/x/crypto/bcrypt"

	"github.com/uptrace/go-realworld-example-app/rwe"
)

var errUserNotFound = errors.New("Not Registered email or invalid password")

func listUsers(c *gin.Context) {
	var users []User
	if err := rwe.PGMain().
		Model(&users).
		Select(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"users": users})
}

func currentUser(c *gin.Context) {
	user := c.MustGet("user").(*User)
	c.JSON(200, gin.H{"users": user})
}

func createUser(c *gin.Context) {
	user := new(User)
	if err := c.BindJSON(user); err != nil {
		c.Error(err)
		return
	}

	var err error
	user.PasswordHash, err = hashPassword(user.Password)
	if err != nil {
		c.Error(err)
		return
	}

	_, err = rwe.PGMain().
		Model(user).
		Insert()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"user": user})
}

func hashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func loginUser(c *gin.Context) {
	user := new(User)
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := c.BindJSON(&in)

	if err := rwe.PGMain().
		Model(user).
		Where("email = ?", in.Email).
		Select(); err != nil {
		if err == pg.ErrNoRows {
			err = errUserNotFound
		}

		c.Error(err)
		return
	}

	if err = comparePasswords(user.PasswordHash, in.Password); err != nil {
		c.Error(err)
		return
	}

	UpdateContextUserModel(c, user.ID)

	var out struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}

	out.Email = user.Email
	out.Token = newToken(user.ID)

	c.JSON(http.StatusOK, gin.H{"user": out})
}

func comparePasswords(hash, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return errUserNotFound
	}
	return nil
}
