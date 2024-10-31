package middleware

import (
	"ddup-apis/internal/logger"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		statusCode := c.Writer.Status()

		if query != "" {
			path = fmt.Sprintf("%s?%s", path, query)
		}

		logger.Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("ip", clientIP),
			zap.Int("size", c.Writer.Size()),
		)

		if len(c.Errors) > 0 {
			logger.Error("HTTP Request Error",
				zap.String("errors", c.Errors.String()),
			)
		}
	}
}
