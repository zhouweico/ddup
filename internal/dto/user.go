package dto

import "time"

// 请求DTO
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname" binding:"required,min=2,max=20"`
	Email    string `json:"email,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
	Location string `json:"location,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Gender   string `json:"gender,omitempty"`
	Avatar   string `json:"avatar"`
	Language string `json:"language" binding:"omitempty,oneof=zh-CN en-US"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// 响应DTO
type UserResponse struct {
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Mobile    string     `json:"mobile"`
	Location  string     `json:"location"`
	Nickname  string     `json:"nickname"`
	Bio       string     `json:"bio"`
	Gender    string     `json:"gender"`
	Birthday  *time.Time `json:"birthday"`
	Avatar    string     `json:"avatar"`
	LastLogin *time.Time `json:"lastLogin"`
	Language  string     `json:"language"`
}

type LoginResponse struct {
	Token     string       `json:"token"`
	CreatedAt time.Time    `json:"createdAt"`
	ExpiresIn int64        `json:"expiresIn"`
	ExpiredAt time.Time    `json:"expiredAt"`
	User      UserResponse `json:"user"`
}
