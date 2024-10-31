package middleware

import (
	"ddup-apis/internal/errors"
	"ddup-apis/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var appErr *errors.AppError

			// 转换为应用错误
			if e, ok := err.(*errors.AppError); ok {
				appErr = e
			} else {
				appErr = errors.Wrap(err, "未知错误")
			}

			// 记录错误日志
			logger.Error("请求处理失败",
				zap.Error(appErr.Err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)

			// 返回错误响应
			c.JSON(appErr.Code, gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
			c.Abort()
		}
	}
}
