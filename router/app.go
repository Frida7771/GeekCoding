package router

import (
	_ "GeekCoding/docs"
	"GeekCoding/middlewares"

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

	//管理员私有方法
	r.POST("/problem-create", middlewares.AuthAdmin(), service.ProblemCreate)

	// 分页获取分类列表
	r.GET("/category-list", middlewares.AuthAdmin(), service.GetCategoryList)
	//create category
	r.POST("/category-create", middlewares.AuthAdmin(), service.CreateCategory)
	//delete category
	r.DELETE("/category-delete", middlewares.AuthAdmin(), service.DeleteCategory)
	//update category
	r.PUT("/category-update", middlewares.AuthAdmin(), service.UpdateCategory)

	return r
}
