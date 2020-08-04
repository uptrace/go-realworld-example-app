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

func setUserToken(user *User) (*User, error) {
	token, err := createUserToken(user.ID, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	user.Token = token
	return user, nil
}

func currentUser(c *gin.Context) {
	user, _ := c.Get("user")

	user, err := setUserToken(user.(*User))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"user": user})
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

	if _, err = rwe.PGMain().
		ModelContext(c, user).
		Insert(); err != nil {
		c.Error(err)
		return
	}

	user, err = setUserToken(user)
	if err != nil {
		c.Error(err)
		return
	}

	user.Password = ""
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
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&in); err != nil {
		return
	}

	user := new(User)
	if err := rwe.PGMain().
		ModelContext(c.Request.Context(), user).
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

	user, err := setUserToken(user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"user": user})
}

func comparePasswords(hash, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return errUserNotFound
	}
	return nil
}

func updateUser(c *gin.Context) {
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

	authUser, _ := c.Get("user")
	if _, err = rwe.PGMain().
		ModelContext(c, user).
		Set("email = ?", user.Email).
		Set("username = ?", user.Username).
		Set("password_hash = ?", user.PasswordHash).
		Set("image = ?", user.Image).
		Set("bio = ?", user.Bio).
		Where("id = ?", authUser.(*User).ID).
		Update(); err != nil {
		c.Error(err)
		return
	}

	user, err = setUserToken(user)
	if err != nil {
		c.Error(err)
		return
	}

	user.Password = ""
	c.JSON(200, gin.H{"user": user})
}
