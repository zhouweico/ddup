package middleware

import (
	"ddup-apis/internal/db"
	"ddup-apis/internal/logger"
	"time"

	"go.uber.org/zap"
)

func PeriodicHealthCheck(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			if err := db.Ping(); err != nil {
				logger.Error("数据库健康检查失败", zap.Error(err))
				// 这里可以添加告警通知逻辑
			}
		}
	}()
}
