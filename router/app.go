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

	// problem (with query rate limiting)
	r.GET("/problem-list", middlewares.QueryRateLimit(), service.GetProblemList)
	r.GET("/problem-detail", middlewares.QueryRateLimit(), service.GetProblemDetail)

	//user
	r.GET("/user-detail", middlewares.QueryRateLimit(), service.GetUserDetail)
	r.POST("/login", middlewares.LoginRateLimit(), service.Login)
	r.POST("/send-code", middlewares.SendCodeRateLimit(), service.SendCode)
	r.POST("/register", middlewares.RegisterRateLimit(), service.Register)
	//排行榜
	r.GET("/rank-list", middlewares.QueryRateLimit(), service.GetRankList)

	//submit
	r.GET("/submit-list", middlewares.QueryRateLimit(), service.GetSubmitList)
	r.GET("/submit-status", middlewares.QueryRateLimit(), service.GetSubmitStatus)

	//管理员私有方法
	authAdmin := r.Group("/admin", middlewares.AuthAdmin())

	//create problem (with admin operation rate limiting)
	authAdmin.POST("/problem-create", middlewares.AdminOperationRateLimit(), service.ProblemCreate)
	//update problem
	authAdmin.PUT("/problem-update", middlewares.AdminOperationRateLimit(), service.ProblemUpdate)

	// 分页获取分类列表
	authAdmin.GET("/category-list", middlewares.QueryRateLimit(), service.GetCategoryList)
	//create category
	authAdmin.POST("/category-create", middlewares.AdminOperationRateLimit(), service.CreateCategory)
	//delete category (dangerous operation, stricter rate limit)
	authAdmin.DELETE("/category-delete", middlewares.AdminOperationRateLimit(), service.DeleteCategory)
	//update category
	authAdmin.PUT("/category-update", middlewares.AdminOperationRateLimit(), service.UpdateCategory)

	//user private method
	authUser := r.Group("/user", middlewares.AuthUser())
	//submit code (with rate limiting)
	authUser.POST("/submit", middlewares.SubmitRateLimit(), service.SubmitCode)

	return r
}
