package middleware

import (
	"ddup-apis/internal/db"
	"log"
	"time"
)

func PeriodicHealthCheck(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			if err := db.Ping(); err != nil {
				log.Printf("数据库健康检查失败: %v", err)
				// 这里可以添加告警通知逻辑
			}
		}
	}()
}
