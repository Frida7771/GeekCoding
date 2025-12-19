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

	// problem
	r.GET("/problem-list", service.GetProblemList)
	r.GET("/problem-detail", service.GetProblemDetail)

	//user
	r.GET("/user-detail", service.GetUserDetail)
	r.POST("/login", service.Login)
	r.POST("/send-code", service.SendCode)
	r.POST("/register", service.Register)
	//排行榜
	r.GET("/rank-list", service.GetRankList)

	//submit
	r.GET("/submit-list", service.GetSubmitList)

	return r
}
