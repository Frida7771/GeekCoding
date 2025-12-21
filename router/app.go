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
	authAdmin := r.Group("/admin", middlewares.AuthAdmin())

	//create problem
	authAdmin.POST("/problem-create", service.ProblemCreate)
	//update problem
	authAdmin.PUT("/problem-update", service.ProblemUpdate)

	// 分页获取分类列表
	authAdmin.GET("/category-list", service.GetCategoryList)
	//create category
	authAdmin.POST("/category-create", service.CreateCategory)
	//delete category
	authAdmin.DELETE("/category-delete", service.DeleteCategory)
	//update category
	authAdmin.PUT("/category-update", service.UpdateCategory)

	//user private method
	authUser := r.Group("/user", middlewares.AuthUser())
	//submit code
	authUser.POST("/submit", service.SubmitCode)

	return r
}
