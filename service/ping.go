package service

import "github.com/gin-gonic/gin"

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func Root(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to GeekCoding API",
		"endpoints": gin.H{
			"ping":        "/ping",
			"problemList": "/problem-list",
			"swagger":     "/swagger/index.html",
		},
	})
}
