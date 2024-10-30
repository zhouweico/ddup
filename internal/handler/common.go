package handler

import (
	"time"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`    // 响应码
	Message string      `json:"message"` // 响应信息
	Data    interface{} `json:"data"`    // 响应数据
}

// TokenInfo Token详细信息
type TokenInfo struct {
	Token     string    `json:"token"`     // JWT token
	CreatedAt time.Time `json:"createdAt"` // 创建时间
	ExpiresIn int64     `json:"expiresIn"` // 有效时长(秒)
	ExpiredAt time.Time `json:"expiredAt"` // 过期时间
}

// LoginResponse 登录响应数据
type LoginResponse struct {
	TokenInfo TokenInfo `json:"tokenInfo"` // Token信息
	UserInfo  User      `json:"userInfo"`  // 用户信息
}

// SignupResponse 注册响应数据
type SignupResponse struct {
	UserInfo User `json:"userInfo"` // 用户信息
}

// User 用户信息
type User struct {
	Username string `json:"username"` // 用户名
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
}
