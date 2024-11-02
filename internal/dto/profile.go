package dto

import (
	"encoding/json"
	"time"
)

// CreateProfileRequest 创建个人资料请求
type CreateProfileRequest struct {
	Type         string          `json:"type" binding:"required" example:"education"`
	Title        string          `json:"title" binding:"required" example:"测试大学"`
	Year         *int            `json:"year" example:"2020"`
	StartDate    *time.Time      `json:"start_date" example:"2020-09-01T00:00:00Z"`
	EndDate      *time.Time      `json:"end_date" example:"2024-06-30T00:00:00Z"`
	Organization string          `json:"organization" example:"测试大学"`
	Location     string          `json:"location" example:"北京"`
	URL          string          `json:"url" example:"https://example.com"`
	Description  string          `json:"description" example:"这是一段描述"`
	Metadata     json.RawMessage `json:"metadata" swaggertype:"string" example:"{\"degree\":\"学士\"}"`
	Visibility   string          `json:"visibility" example:"public"`
}

// UpdateProfileRequest 更新个人资料请求
type UpdateProfileRequest struct {
	Title        string          `json:"title"`
	Year         *int            `json:"year"`
	StartDate    *time.Time      `json:"start_date"`
	EndDate      *time.Time      `json:"end_date"`
	Organization string          `json:"organization"`
	Location     string          `json:"location"`
	URL          string          `json:"url"`
	Description  string          `json:"description"`
	Metadata     json.RawMessage `json:"metadata"`
	Visibility   string          `json:"visibility"`
}

// UpdateDisplayOrderRequest 更新显示顺序请求
type UpdateDisplayOrderRequest struct {
	Items []struct {
		ID    uint `json:"id"`
		Order int  `json:"order"`
	} `json:"items" binding:"required"`
}

// ProfileResponse 个人资料响应
type ProfileResponse struct {
	ID           uint            `json:"id" example:"1"`
	Type         string          `json:"type" example:"education"`
	Title        string          `json:"title" example:"测试大学"`
	Year         *int            `json:"year" example:"2020"`
	StartDate    *time.Time      `json:"start_date" example:"2020-09-01T00:00:00Z"`
	EndDate      *time.Time      `json:"end_date" example:"2024-06-30T00:00:00Z"`
	Organization string          `json:"organization" example:"测试大学"`
	Location     string          `json:"location" example:"北京"`
	URL          string          `json:"url" example:"https://example.com"`
	Description  string          `json:"description" example:"这是一段描述"`
	Metadata     json.RawMessage `json:"metadata" swaggertype:"string" example:"{\"degree\":\"学士\"}"`
	DisplayOrder int             `json:"display_order" example:"0"`
	Visibility   string          `json:"visibility" example:"public"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
