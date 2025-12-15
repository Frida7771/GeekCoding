package service

import (
	"GeekCoding/define"
	"GeekCoding/models"
	"log"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProblemList
//@Tafs Public Method
// @Summary      Problem List
// @Param        page  query     int     false  "page"
// @Param        size  query     int     false  "size"
// @Param        keyword  query     string     false  "keyword"
// @Success      200   {string}    json "{"code": 200, "message""data": ""}"
// @Router       /problem-list [get]

func GetProblemList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery(define.DefaultSize, define.DefaultSize))
	if err != nil {
		log.Println("Get Problem List error: ", err)
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery(define.DefaultPage, define.DefaultPage))
	if err != nil {
		log.Println("Get Problem List Page error: ", err)
		return
	}
	page = (page - 1) * size

	var count int64

	keyword := c.Query("keyword")

	list := make([]models.Problem, 0)

	tx := models.GetProblemList(keyword)

	err = tx.Count(&count).Limit(size).Find(&list).Error
	if err != nil {
		log.Println("Get Problem List error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
		"count": count,
	})
}
