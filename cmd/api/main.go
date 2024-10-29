package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zhouweico/ddup-apis/internal/config"
	"github.com/zhouweico/ddup-apis/internal/handler"
	"github.com/zhouweico/ddup-apis/internal/middleware"
)

func main() {
	// 初始化配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 打印配置信息（不包含敏感信息）
	log.Printf("Configuration loaded successfully: %s", cfg)

	// 初始化 Gin
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// 路由组
	api := r.Group("/api/v1")
	{
		// 公开路由
		public := api.Group("/")
		{
			public.POST("/login", handler.Login)
			public.POST("/register", handler.Register)
		}

		// 需要认证的路由
		protected := api.Group("/").Use(middleware.Auth())
		{
			protected.GET("/users", handler.GetUsers)
			protected.GET("/users/:id", handler.GetUser)
			protected.PUT("/users/:id", handler.UpdateUser)
			protected.DELETE("/users/:id", handler.DeleteUser)
		}
	}

	// 启动服务器
	r.Run(cfg.Server.Address)
}
