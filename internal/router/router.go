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

	// 添加全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Cors())
	r.Use(gin.Recovery())

	// 初始化 services
	userService := service.NewUserService(db.DB)

	// 初始化 handlers
	userHandler := handler.NewUserHandler(userService)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 公开路由
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)

		// 需要认证的路由
		auth := v1.Group("")
		auth.Use(middleware.JWTAuth(userService))
		{
			auth.POST("/logout", userHandler.Logout)
			auth.PUT("/user", userHandler.UpdateUser)
			auth.DELETE("/user", userHandler.DeleteUser)
			auth.GET("/users", userHandler.GetUsers)               // 获取用户列表
			auth.GET("/users/:uuid", userHandler.GetUser)          // 使用 uuid 参数
			auth.PUT("/user/password", userHandler.ChangePassword) // 添加修改密码路由
		}
	}

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
