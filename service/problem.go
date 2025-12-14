package service

import (
	"GeekCoding/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProblemList
// @Summary      Get Problem List
// @Description  Get a list of problems with pagination
// @Tags         problems
// @Accept       json
// @Produce      json
// @Param        page  query     int     false  "Page number"  default(1)
// @Param        size  query     int     false  "Page size"    default(10)
// @Success      200   {object}  map[string]interface{}  "Success response"
// @Router       /problem-list [get]
func GetProblemList(c *gin.Context) {
	models.GetProblemList()
	c.String(http.StatusOK, "Get Problem List")
}
