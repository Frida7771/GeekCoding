package router

import (
	_ "GeekCoding/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"GeekCoding/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	//swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/", service.Root)

	r.GET("/ping", service.Ping)

	r.GET("/problem-list", service.GetProblemList)

	return r
}
