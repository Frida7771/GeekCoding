package router

import (
	"GeekCoding/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", service.Ping)

	r.GET("/problem-list", service.GetProblemList)

	return r
}
