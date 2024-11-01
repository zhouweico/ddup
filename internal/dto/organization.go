package dto

type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Email       string `json:"email" binding:"omitempty,email"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Website     string `json:"website"`
}

type UpdateOrganizationRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Email       string `json:"email" binding:"omitempty,email"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Website     string `json:"website"`
}

type OrganizationResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Website     string `json:"website"`
	Role        string `json:"role"` // 当前用户在组织中的角色
}

type AddMemberRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=admin member"`
}

type UpdateMemberRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member"`
}

type MemberResponse struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type JoinOrganizationRequest struct {
	OrganizationID uint `json:"organization_id" binding:"required"`
}
