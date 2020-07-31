package org

import (
	"github.com/gin-gonic/gin"
	"github.com/uptrace/go-realworld-example-app/rwe"
	"golang.org/x/crypto/bcrypt"
)

func listUsers(c *gin.Context) {
	var users []User
	if err := rwe.PGMain().
		Model(&users).
		Select(); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"users": users})
}

func showUser(c *gin.Context) {
	user := new(User)
	if err := rwe.PGMain().
		Model(user).
		Where("user_id = ?", c.Param("user_id")).
		Select(); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"users": user})
}

func createUser(c *gin.Context) {
	user := new(User)
	c.BindJSON(user)

	var err error
	user.PasswordHash, err = hashPassword(user.PasswordHash)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = rwe.PGMain().
		Model(user).
		Insert()
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
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
