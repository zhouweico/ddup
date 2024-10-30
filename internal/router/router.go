package router

import (
	"ddup-apis/internal/handler"

	_ "ddup-apis/docs" // 这行很重要，需要导入生成的 docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由
		v1.POST("/sign-up", handler.Signup) // 注册
		v1.POST("/login", handler.Login)    // 登录
		// ... 其他路由
	}

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
