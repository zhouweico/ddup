package main

import (
	"log"

	"ddup-apis/internal/config"
	"ddup-apis/internal/db"
	"ddup-apis/internal/middleware"
	"ddup-apis/internal/router"
)

// @title DDUP API
// @version 1.0
// @description DDUP 服务 API 文档
// @host localhost:8080
// @BasePath /api/v1

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	config.SetConfig(*cfg)

	// 初始化数据库
	if err := db.InitDB(cfg); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化路由
	r := router.SetupRouter()

	// 启动定期健康检查
	middleware.PeriodicHealthCheck(cfg.HealthCheck.Interval)

	// 启动服务器
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
