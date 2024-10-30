package handler

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`    // 响应码
	Message string      `json:"message"` // 响应信息
	Data    interface{} `json:"data"`    // 响应数据
}

// 响应相关的方法
func SendSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

func SendError(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Code:    status,
		Message: message,
	})
}

// TokenInfo Token详细信息
type TokenInfo struct {
	Token     string    `json:"token"`     // JWT token
	CreatedAt time.Time `json:"createdAt"` // 创建时间
	ExpiresIn int64     `json:"expiresIn"` // 有效时长(秒)
	ExpiredAt time.Time `json:"expiredAt"` // 过期时间
}

// User 用户信息
type User struct {
	Username string `json:"username"` // 用户名
}

// API 响应数据结构
type (
	// LoginResponse 登录响应数据
	LoginResponse struct {
		TokenInfo TokenInfo `json:"tokenInfo"` // Token信息
		UserInfo  User      `json:"userInfo"`  // 用户信息
	}

	// RegisterResponse 注册响应数据
	RegisterResponse struct {
		UserInfo User `json:"userInfo"` // 用户信息
	}

	// ErrorResponse 错误响应结构
	ErrorResponse struct {
		Code    int    `json:"code"`    // 错误码
		Message string `json:"message"` // 错误信息
	}
)

// UserDetail 用户详细信息
type UserDetail struct {
	ID          uint      `json:"id"`          // 用户ID
	Username    string    `json:"username"`    // 用户名
	Email       string    `json:"email"`       // 邮箱
	Password    string    `json:"-"`           // 密码，不返回给前端
	Nickname    string    `json:"nickname"`    // 昵称
	Bio         string    `json:"bio"`         // 简介
	Gender      string    `json:"gender"`      // 性别
	Avatar      string    `json:"avatar"`      // 头像
	Status      int       `json:"status"`      // 状态
	LastLoginAt time.Time `json:"lastLoginAt"` // 最后登录时间
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
}
