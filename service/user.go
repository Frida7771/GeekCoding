package service

import (
	"log"
	"net/http"
	"strconv"

	"GeekCoding/define"
	"GeekCoding/models"

	"GeekCoding/help"

	"time"

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

// SendCode
// @Tags         Public Method
// @Summary      Send Code
// @Param        email  formData     string     false  "email"
// @Success      200   {string}    json "{"code": 200, "data": ""}"
// @Router       /send-code [post]
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "email is required",
		})
		return
	}
	code := help.GetRandomCode()
	models.RDB.Set(c, email, code, time.Second*300)
	err := help.SendCode(email, code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "send code error: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "send code success",
	})
}

// Register
// @Tags         Public Method
// @Summary      Register
// @Param        email  formData     string     true  "email"
// @Param        code  formData     string     true  "code"
// @Param        name  formData     string     true  "name"
// @Param        password  formData     string     true  "password"
// @Param        phone  formData     string     false  "phone"
// @Success      200   {string}    json "{"code": 200, "data": ""}"
// @Router       /register [post]
func Register(c *gin.Context) {
	email := c.PostForm("email")
	code := c.PostForm("code")
	name := c.PostForm("name")
	password := c.PostForm("password")
	phone := c.PostForm("phone")

	if email == "" || code == "" || name == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  " parameters are required",
		})
		return
	}
	//check code
	systemCode, err := models.RDB.Get(c, email).Result()
	if err != nil {
		log.Printf("get code error: %v \n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码错误，请重新发送",
		})
		return
	}
	if systemCode != code {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "code is incorrect",
		})
		return
	}

	//判断邮箱是否已经注册
	var cnt int64
	err = models.DB.Where("email = ?", email).Model(&models.User_Basic{}).Count(&cnt).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get User Errot " + err.Error(),
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱已注册，请直接登录",
		})
		return
	}

	//input data
	userIdentity := help.GetUUID()
	data := &models.User_Basic{
		Identity: userIdentity,
		Name:     name,
		Email:    email,
		Password: help.MD5(password),
		Phone:    phone,
	}

	err = models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "create user error: " + err.Error(),
		})
		return
	}

	//generate token
	token, err := help.GenerateToken(userIdentity, name)
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

// GetRankList
// @Tags         Public Method
// @Summary      Get Rank List
// @Param        page  query     int     false  "page"
// @Param        size  query     int     false  "size"
// @Success      200   {string}    json "{"code":"200","msg","","data": ""}"
// @Router       /rank-list [get]
func GetRankList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "invalid page parameter",
		})
		return
	}

	offset := (page - 1) * size

	var count int64
	list := make([]*models.User_Basic, 0)

	db := models.DB.Model(&models.User_Basic{})

	db.Count(&count)

	err = db.
		Order("finish_problem_num DESC, submit_num ASC").
		Offset(offset).
		Limit(size).
		Find(&list).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Rank List Error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"list":  list,
			"count": count,
		},
	})
}
