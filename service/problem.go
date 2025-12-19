package service

import (
	"GeekCoding/define"
	"GeekCoding/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProblemList
// @Tags         Public Method
// @Summary      Get Problem List
// @Param        page  query     int     false  "page"
// @Param        size  query     int     false  "size"
// @Param        keyword  query     string     false  "keyword"
// @Success      200   {string}    json "{"code":"200","msg","","data": ""}"
// @Router       /problem-list [get]
func GetProblemList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetProblemList Page strconv Error:", err)
		return
	}
	page = (page - 1) * size
	var count int64
	keyword := c.Query("keyword")
	categoryIdentity := c.Query("category_identity")

	list := make([]*models.ProblemBasic, 0)
	err = models.GetProblemList_Basic(keyword, categoryIdentity).Distinct("`problem_basic`.`id`").Count(&count).Error
	if err != nil {
		log.Println("GetProblemList Count Error:", err)
		return
	}
	err = models.GetProblemList_Basic(keyword, categoryIdentity).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get Problem List Error:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// GetProblemDetail
// @Tags         Public Method
// @Summary      Get Problem Detail
// @Param        identity  query     string     true  "problem identity"
// @Success      200   {string}    json "{"code": 200, "data": ""}"
// @Router       /problem-detail [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "problem identity is required",
		})
		return
	}

	data := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?", identity).
		Preload("ProblemCategorys").Preload("ProblemCategorys.CategoryBasic").First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "problem not found",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "get problem detail error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}
