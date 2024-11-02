package dto

import "time"

type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	DisplayName string `json:"display_name" binding:"required,min=2,max=100"`
	Email       string `json:"email" binding:"omitempty,email"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Website     string `json:"website"`
}

type UpdateOrganizationRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	DisplayName string `json:"display_name" binding:"omitempty,min=2,max=100"`
	Email       string `json:"email" binding:"omitempty,email"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Website     string `json:"website"`
}

type OrganizationResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Website     string `json:"website"`
	Role        string `json:"role"` // 当前用户在组织中的角色
}

type UpdateMemberRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member"`
}

type MemberResponse struct {
	Username string    `json:"username"`
	Nickname string    `json:"nickname"`
	Email    string    `json:"email"`
	Avatar   string    `json:"avatar"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type JoinOrganizationRequest struct {
	OrganizationID uint `json:"organization_id" binding:"required"`
}

type AddOrganizationMemberRequest struct {
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin member"`
}
