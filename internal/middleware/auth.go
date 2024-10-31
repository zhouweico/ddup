package middleware

import (
	"ddup-apis/internal/service"
	"net/http"

	"ddup-apis/internal/errors"

	"github.com/gin-gonic/gin"
)

func JWTAuth(userService service.IUserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		result, err := userService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.Error(errors.Wrap(err, "Token 验证失败"))
			c.Abort()
			return
		}

		if !result.IsValid {
			c.Error(errors.New(http.StatusUnauthorized, "Token 已失效", nil))
			c.Abort()
			return
		}

		c.Set("userID", result.UserID)
		c.Set("username", result.Username)
		c.Set("IsValid", result.IsValid)
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
