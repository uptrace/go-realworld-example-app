package rwe

import (
	"errors"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var (
	Router *gin.Engine
	API    *gin.RouterGroup
)

func init() {
	Router = gin.Default()
	Router.Use(corsMiddleware)
	Router.Use(errorMiddleware)
	Router.Use(rateLimitMiddleware)
	Router.Use(otelgin.Middleware("rwe"))

	API = Router.Group("/api")
}

func errorMiddleware(c *gin.Context) {
	c.Next()

	ginErr := c.Errors.Last()
	if ginErr == nil {
		return
	}

	switch err := ginErr.Err.(type) {
	default:
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
}

func corsMiddleware(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	if origin == "" {
		return
	}

	if c.Request.Method == http.MethodOptions {
		h := c.Writer.Header()
		h.Set("Access-Control-Allow-Origin", origin)
		h.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,HEAD")
		h.Set("Access-Control-Allow-Headers", "authorization,content-type")
		h.Set("Access-Control-Max-Age", "86400")
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	h := c.Writer.Header()
	h.Set("Access-Control-Allow-Origin", origin)
}

func rateLimitMiddleware(c *gin.Context) {
	if c.Request.Method == http.MethodOptions {
		return
	}

	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		c.Error(err)
		return
	}

	rateKey := "rl:" + host
	limit := redis_rate.PerMinute(100)

	res, err := RateLimiter().Allow(c.Request.Context(), rateKey, limit)
	if err != nil {
		c.Error(err)
		return
	}
	if res.Allowed == 0 {
		c.AbortWithError(http.StatusTooManyRequests, errors.New("rate limited"))
		return
	}
}
