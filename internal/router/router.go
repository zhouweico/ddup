package router

import (
	"ddup-apis/internal/db"
	"ddup-apis/internal/handler"
	"ddup-apis/internal/middleware"
	"ddup-apis/internal/service"

	_ "ddup-apis/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 初始化 services
	userService := service.NewUserService(db.DB)

	// 初始化 handlers
	userHandler := handler.NewUserHandler(userService)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 公开路由
		v1.POST("/sign-up", userHandler.Signup)
		v1.POST("/login", userHandler.Login)

		// 需要认证的路由
		auth := v1.Group("")
		auth.Use(middleware.JWTAuth(userService))
		{
			// 添加需要认证的路由
		}
	}

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
