package org

import (
	"github.com/gin-gonic/gin"
	"github.com/uptrace/go-realworld-example-app/rwe"
	"golang.org/x/crypto/bcrypt"
)

func listUsers(c *gin.Context) error {
	var users []User
	if err := rwe.PGMain().
		Model(&users).
		Select(); err != nil {
		return err
	}

	c.JSON(200, gin.H{"users": users})
	return nil
}

func showUser(c *gin.Context) error {
	user := new(User)
	if err := rwe.PGMain().
		Model(user).
		Where("user_id = ?", c.Param("user_id")).
		Select(); err != nil {
		return err
	}

	c.JSON(200, gin.H{"users": user})
	return nil
}

func createUser(c *gin.Context) error {
	user := new(User)
	c.BindJSON(user)

	var err error
	user.PasswordHash, err = hashPassword(user.PasswordHash)
	if err != nil {
		return err
	}

	_, err = rwe.PGMain().
		Model(user).
		Insert()
	if err != nil {
		return err
	}

	c.JSON(200, gin.H{"user": user})
	return nil
}

func hashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
