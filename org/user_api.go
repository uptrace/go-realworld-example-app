package org

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"golang.org/x/crypto/bcrypt"

	"github.com/uptrace/go-realworld-example-app/rwe"
)

var errUserNotFound = errors.New("Not Registered email or invalid password")

type UserOut struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

func newUserOut(user *User) (*UserOut, error) {
	token, err := createUserToken(user.ID, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &UserOut{
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
		Image:    user.Image,
		Token:    token,
	}, nil
}

func currentUser(c *gin.Context) {
	user, _ := c.Get("user")

	userOut, err := newUserOut(user.(*User))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"user": userOut})
}

func createUser(c *gin.Context) {
	user := new(User)
	if err := c.BindJSON(user); err != nil {
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

	userOut, err := newUserOut(user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"user": userOut})
}

func hashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func loginUser(c *gin.Context) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&in); err != nil {
		return
	}

	user := new(User)
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

	if err := comparePasswords(user.PasswordHash, in.Password); err != nil {
		c.Error(err)
		return
	}

	userOut, err := newUserOut(user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"user": userOut})
}

func comparePasswords(hash, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return errUserNotFound
	}
	return nil
}
