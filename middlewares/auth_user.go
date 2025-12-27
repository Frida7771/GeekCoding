package middlewares

import (
	"GeekCoding/help"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		userClaims, err := help.AnalyzeToken(auth)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized Authentication",
			})
			return
		}
		if userClaims == nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "Unauthorized Authentication",
			})
			return
		}
		// 设置用户信息到 context，供后续使用
		c.Set("user", userClaims)
		c.Next()
	}
}
