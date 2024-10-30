package response

import (
	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`    // 响应码
	Message string      `json:"message"` // 响应信息
	Data    interface{} `json:"data"`    // 响应数据
}

// SendSuccess 发送成功响应
func SendSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// SendError 发送错误响应
func SendError(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Code:    status,
		Message: message,
	})
}
