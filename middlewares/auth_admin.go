package middlewares

import (
	"GeekCoding/help"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthAdmin() gin.HandlerFunc {
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
		if userClaims == nil || userClaims.IsAdmin != 1 {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusForbidden,
				"msg":  "Unauthorized Admin",
			})
			return
		}
		c.Next()
	}
}
