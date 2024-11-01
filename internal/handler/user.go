package handler

import (
	"ddup-apis/internal/dto"
	"ddup-apis/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// @Summary 用户注册
// @Description 用户注册
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "注册信息"
// @Success 200 {object} Response "注册成功"
// @Failure 400 {object} Response "无效的请求参数"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.userService.Register(c.Request.Context(), &req); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "注册成功", nil)
}

// @Summary 用户登录
// @Description 用户登录
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录信息"
// @Success 200 {object} Response{data=dto.LoginResponse} "登录成功"
// @Failure 400 {object} Response "无效的请求参数"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	resp, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		SendError(c, http.StatusUnauthorized, err.Error())
		return
	}

	SendSuccess(c, "登录成功", resp)
}

// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Produce json
// @Security Bearer
// @Success 200 {object} Response{data=dto.UserResponse} "获取成功"
// @Failure 401 {object} Response "未授权"
// @Router /user [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, http.StatusUnauthorized, "未授权")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID.(uint))
	if err != nil {
		SendError(c, http.StatusNotFound, err.Error())
		return
	}

	SendSuccess(c, "获取成功", user)
}

// @Summary 更新用户信息
// @Description 更新用户基本信息
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.UpdateUserRequest true "用户信息"
// @Success 200 {object} Response "更新成功"
// @Failure 400 {object} Response "无效的请求参数"
// @Failure 401 {object} Response "未授权"
// @Router /user [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, http.StatusUnauthorized, "未授权")
		return
	}

	if err := h.userService.UpdateUser(c.Request.Context(), userID.(uint), &req); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "更新用户信息成功", nil)
}

// @Summary 修改密码
// @Description 修改用户密码
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.ChangePasswordRequest true "密码信息"
// @Success 200 {object} Response "修改成功"
// @Failure 400 {object} Response "无效的请求参数"
// @Failure 401 {object} Response "未授权"
// @Router /user/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, http.StatusUnauthorized, "未授权")
		return
	}

	if err := h.userService.ChangePassword(c.Request.Context(), userID.(uint), &req); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "密码修改成功", nil)
}

// @Summary 退出登录
// @Description 用户退出登录
// @Produce json
// @Security Bearer
// @Success 200 {object} Response "退出成功"
// @Failure 401 {object} Response "未授权"
// @Router /logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		SendError(c, http.StatusUnauthorized, "未授权")
		return
	}

	if err := h.userService.Logout(c.Request.Context(), token); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "退出成功", nil)
}

// @Summary 删除用户
// @Description 删除用户账号
// @Produce json
// @Security Bearer
// @Success 200 {object} Response "删除成功"
// @Failure 401 {object} Response "未授权"
// @Router /user [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, http.StatusUnauthorized, "未授权")
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), userID.(uint)); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "删除成功", nil)
}
