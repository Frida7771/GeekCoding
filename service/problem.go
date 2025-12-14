package service

import (
	"GeekCoding/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProblemList(c *gin.Context) {
	models.GetProblemList()
	c.String(http.StatusOK, "Get Problem List")

}
