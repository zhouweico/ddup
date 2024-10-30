package handler

import (
	"ddup-apis/internal/service"
	"ddup-apis/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname,omitempty"` // 用户昵称
	Email    string `json:"email,omitempty"`    // 邮箱
	Mobile   string `json:"mobile,omitempty"`   // 手机号
	Location string `json:"location,omitempty"` // 位置
	Bio      string `json:"bio,omitempty"`      // 用户简介
	Gender   string `json:"gender,omitempty"`   // 性别
	Avatar   string `json:"avatar,omitempty"`   // 头像URL
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int64  `json:"total"` // 总记录数
	Users []User `json:"users"` // 用户列表
}

// UserDetailResponse 用户详情响应
type UserDetailResponse struct {
	UserID    string     `json:"userid"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Mobile    string     `json:"mobile"`   // 手机号
	Location  string     `json:"location"` // 位置
	Nickname  string     `json:"nickname"`
	Bio       string     `json:"bio"`
	Gender    string     `json:"gender"`
	Birthday  *time.Time `json:"birthday"`
	Avatar    string     `json:"avatar"`
	LastLogin *time.Time `json:"lastLogin"`
}

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// @Summary 用户注册
// @Description 用户注册
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 200 {object} handler.Response{data=handler.RegisterResponse} "注册成功"
// @Failure 400 {object} handler.Response "无效的请求参数"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	UserID, err := h.userService.Register(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	SendSuccess(c, "注册成功", RegisterResponse{
		UserInfo: User{
			Username: req.Username,
			UserID:   UserID,
		},
	})
}

// @Summary 用户登录
// @Description 用户登录
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} handler.Response{data=handler.LoginResponse} "登录成功"
// @Failure 400 {object} handler.Response "无效的请求参数"
// @Failure 401 {object} handler.Response "用户名或密码错误"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	result, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		SendError(c, http.StatusUnauthorized, err.Error())
		return
	}

	SendSuccess(c, "登录成功", LoginResponse{
		TokenInfo: TokenInfo{
			Token:     result.Token,
			CreatedAt: result.CreatedAt,
			ExpiresIn: result.ExpiresIn,
			ExpiredAt: result.ExpiredAt,
		},
		UserInfo: User{
			Username: result.User.Username,
			UserID:   result.User.UserID,
		},
	})
}

// @Summary 用户退出
// @Description 用户退出登录
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} handler.Response "退出成功"
// @Failure 401 {object} handler.Response "未授权"
// @Router /logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if err := h.userService.Logout(c.Request.Context(), token); err != nil {
		SendError(c, http.StatusInternalServerError, "退出失败")
		return
	}
	SendSuccess(c, "退出成功", nil)
}

// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码(默认1)" default(1)
// @Param page_size query int false "每页数量(默认10)" default(10)
// @Success 200 {object} handler.Response{data=UserListResponse} "获取成功"
// @Failure 401 {object} handler.Response "未授权"
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	// 获取分页参数
	page := utils.StringToInt(c.DefaultQuery("page", "1"))
	pageSize := utils.StringToInt(c.DefaultQuery("page_size", "10"))

	users, total, err := h.userService.GetUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		SendError(c, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	// 转换为 UserDetail 列表
	var userDetails []UserDetail
	for _, user := range users {
		userDetails = append(userDetails, UserDetail{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Bio:       user.Bio,
			Gender:    user.Gender,
			Avatar:    user.Avatar,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	SendSuccess(c, "获取用户列表成功", gin.H{
		"total": total,
		"users": userDetails,
	})
}

// @Summary 获取用户详情
// @Description 获取用户详细信息（仅允许获取自己的信息）
// @Accept json
// @Produce json
// @Security Bearer
// @Param userid path string true "用户ID"
// @Success 200 {object} handler.Response{data=UserDetailResponse} "获取成功"
// @Failure 401 {object} handler.Response "未授权"
// @Failure 403 {object} handler.Response "禁止访问"
// @Router /users/{userid} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("userid")
	user, err := h.userService.GetUserByUserID(c.Request.Context(), userID)
	if err != nil {
		SendError(c, http.StatusNotFound, err.Error())
		return
	}

	SendSuccess(c, "获取成功", UserDetailResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Mobile:    user.Mobile,
		Location:  user.Location,
		Nickname:  user.Nickname,
		Bio:       user.Bio,
		Gender:    user.Gender,
		Birthday:  user.Birthday,
		Avatar:    user.Avatar,
		LastLogin: user.LastLogin,
	})
}

// @Summary 更新用户信息
// @Description 更新用户基本信息
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body UpdateUserRequest true "用户信息"
// @Success 200 {object} handler.Response "更新成功"
// @Failure 400 {object} handler.Response "无效的请求参数"
// @Failure 401 {object} handler.Response "未授权"
// @Router /user [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	userID := c.Param("userid")
	user, err := h.userService.GetUserByUserID(c.Request.Context(), userID)
	if err != nil {
		SendError(c, http.StatusNotFound, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Mobile != "" {
		updates["mobile"] = req.Mobile
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if err := h.userService.UpdateUser(c.Request.Context(), user.ID, updates); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "更新成功", nil)
}

// @Summary 删除用户
// @Description 软除用户账号
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} handler.Response "删除成功"
// @Failure 401 {object} handler.Response "未授权"
// @Failure 500 {object} handler.Response "系统错误"
// @Router /user [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("userid")
	user, err := h.userService.GetUserByUserID(c.Request.Context(), userID)
	if err != nil {
		SendError(c, http.StatusNotFound, err.Error())
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), user.ID); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "删除成功", nil)
}

// @Summary 修改密码
// @Description 用户修改密码
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} handler.Response "密码修改成功"
// @Failure 400 {object} handler.Response "请求参数错误"
// @Failure 401 {object} handler.Response "未授权"
// @Router /user/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	userID := c.Param("userid")
	user, err := h.userService.GetUserByUserID(c.Request.Context(), userID)
	if err != nil {
		SendError(c, http.StatusNotFound, err.Error())
		return
	}

	err = h.userService.ChangePassword(c.Request.Context(), user.ID, req.OldPassword, req.NewPassword)
	if err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	SendSuccess(c, "密码修改成功", nil)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}
