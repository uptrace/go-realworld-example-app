package rwe

import (
	"github.com/gin-gonic/gin"
	"github.com/uptrace/go-realworld-example-app/httperror"
)

var (
	Router = gin.Default()
	API    = Router.Group("/api")
)

func init() {
	API.Use(errorHandler)
}

func errorHandler(c *gin.Context) {
	c.Next()

	ginErr := c.Errors.Last()
	if ginErr == nil {
		return
	}

	switch err := ginErr.Err.(type) {
	case *httperror.Error:
		c.JSON(err.Status, err)
	default:
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	return
}
