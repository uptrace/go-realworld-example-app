package rwe

import "github.com/gin-gonic/gin"

var (
	Router = gin.Default()
	API    = Router.Group("/api")
)
