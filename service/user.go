package service

import (
	"net/http"

	"GeekCoding/models"

	"GeekCoding/help"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

// GetUserDetail
// @Tags         Public Method
// @Summary      Get User Detail
// @Param        identity  query     string     true  "user identity"
// @Success      200   {string}    json "{"code": 200, "data": ""}"
// @Router       /user-detail [get]
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

// Login
// @Tags         Public Method
// @Summary      User Login
// @Param        username  formData     string     false  "username"
// @Param        password  formData     string     false  "password"
// @Success      200   {string}    json "{"code": 200, "data": ""}"
// @Router       /login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "username and password are required",
		})
		return
	}
	//md5
	password = help.MD5(password)
	print(username, password)

	data := new(models.User_Basic)

	err := models.DB.Where("username = ? AND password = ?", username, password).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "username or password is incorrect",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get UserBasic Errot " + err.Error(),
		})
		return
	}

	token, err := help.GenerateToken(data.Identity, data.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "generate token error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}
