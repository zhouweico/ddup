package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

type LoggerConfig struct {
	SkipPaths []string
}

func Logger(config ...LoggerConfig) gin.HandlerFunc {
	var conf LoggerConfig
	if len(config) > 0 {
		conf = config[0]
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method

		// 跳过不需要记录日志的路径
		for _, skipPath := range conf.SkipPaths {
			if path == skipPath {
				c.Next()
				return
			}
		}

		// 处理请求
		c.Next()

		// 收集请求信息
		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()
		if bodySize < 0 {
			bodySize = 0
		}

		if raw != "" {
			path = path + "?" + raw
		}

		// 根据状态码选择颜色
		var statusColor string
		switch {
		case statusCode >= 500:
			statusColor = red
		case statusCode >= 400:
			statusColor = yellow
		case statusCode >= 300:
			statusColor = white
		case statusCode >= 200:
			statusColor = green
		default:
			statusColor = cyan
		}

		// 选择请求方法的颜色
		var methodColor string
		switch method {
		case "GET":
			methodColor = blue
		case "POST":
			methodColor = cyan
		case "PUT":
			methodColor = yellow
		case "DELETE":
			methodColor = red
		case "PATCH":
			methodColor = green
		case "HEAD":
			methodColor = magenta
		case "OPTIONS":
			methodColor = white
		default:
			methodColor = reset
		}

		// 格式化日志输出
		fmt.Printf("[DDUP] %v |%s %3d %s| %13v | %15s | %s %s %s | %s | %d bytes\n",
			end.Format("2006/01/02 - 15:04:05"),
			statusColor, statusCode, reset,
			latency,
			clientIP,
			methodColor, method, reset,
			path,
			bodySize,
		)

		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			fmt.Printf("[GIN-ERROR] %v\n", c.Errors.String())
		}
	}
}
