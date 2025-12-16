package service

import (
	"net/http"

	"GeekCoding/models"

	"github.com/gin-gonic/gin"
)

func GetUserDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "user identity is required",
		})

		return
	}
	data := new(models.User_Basic)
	err := models.DB.Omit("password").Where("identity = ?", identity).First(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get user detail by identity: " + identity + " Error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}
