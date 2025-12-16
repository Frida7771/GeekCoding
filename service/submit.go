package service

import (
	"GeekCoding/define"
	"GeekCoding/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetSubmitList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("Get Submit List Page error: ", err)
		return
	}
	page = (page - 1) * size
	var count int64
	list := make([]models.Submit_Basic, 0)

	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	status, _ := strconv.Atoi(c.Query("status"))

	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	err = tx.Count(&count).Offset(page).Limit(size).Find(&list).Error

	if err != nil {
		log.Println("Get Submit List error: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get submit list error: " + err.Error(),
		})
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
