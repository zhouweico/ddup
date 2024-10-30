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
		raw := c.Request.URL.RawQuery
		method := c.Request.Method

		clientIP := c.ClientIP()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		var logFunc func(string, ...interface{})
		if statusCode >= 500 {
			logFunc = log.Printf
		} else if statusCode >= 400 {
			logFunc = log.Printf
		} else {
			logFunc = log.Printf
		}

		logFunc("[GIN] %s | %3d | %13v | %15s | %s | %d bytes",
			method,
			statusCode,
			latency,
			clientIP,
			path,
			bodySize,
		)

		if len(c.Errors) > 0 {
			log.Printf("[ERROR] %v", c.Errors.String())
		}
	}
}
