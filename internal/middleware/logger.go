package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		if len(c.Errors) > 0 {
			// 记录错误日志
			c.Error(c.Errors.Last())
		}

		// TODO: 使用你的日志库记录请求信息
		log.Printf("[%d] %s %s %v", statusCode, method, path, latency)
	}
}
