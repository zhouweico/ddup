package middleware

import (
	"ddup-apis/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuth(userService service.IUserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			sendError(c, http.StatusUnauthorized, "未提供认证信息")
			c.Abort()
			return
		}

		result, err := userService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			sendError(c, http.StatusUnauthorized, "Token 验证失败")
			c.Abort()
			return
		}

		if !result.IsValid {
			sendError(c, http.StatusUnauthorized, "Token 已失效")
			c.Abort()
			return
		}

		requestUserID := c.Param("userid")
		if requestUserID != "" && requestUserID != result.UserID {
			sendError(c, http.StatusForbidden, "无权访问该资源")
			c.Abort()
			return
		}

		c.Set("userID", result.UserID)
		c.Set("username", result.Username)
		c.Next()
	}
}

func VerifyResourceOwnership(userService service.IUserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestUserID := c.Param("userid")
		token := extractToken(c)
		if token == "" {
			sendError(c, http.StatusUnauthorized, "未提供认证信息")
			c.Abort()
			return
		}

		result, err := userService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			sendError(c, http.StatusUnauthorized, "Token 验证失败")
			c.Abort()
			return
		}
		if !result.IsValid || requestUserID != result.UserID {
			sendError(c, http.StatusForbidden, "无权访问此资源")
			c.Abort()
			return
		}

		c.Next()
	}
}

func sendError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"code":    status,
		"message": message,
	})
}

func extractToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
	}
	return token
}
