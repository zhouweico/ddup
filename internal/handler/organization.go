package handler

import (
	"ddup-apis/internal/dto"
	"ddup-apis/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	service *service.OrganizationService
}

func NewOrganizationHandler(service *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

// @Summary 创建组织
// @Description 创建新的组织，创建者默认成为管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateOrganizationRequest true "组织信息"
// @Success 200 {object} Response{data=dto.OrganizationResponse}
// @Router /api/v1/organizations [post]
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req dto.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	userID := c.GetUint("userID")
	org, err := h.service.CreateOrganization(c.Request.Context(), &req, userID)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "创建组织成功", org)
}

// @Summary 获取用户组织列表
// @Description 获取当前用户所属的所有组织
// @Produce json
// @Security Bearer
// @Success 200 {object} Response{data=[]dto.OrganizationResponse}
// @Router /api/v1/organizations [get]
func (h *OrganizationHandler) GetUserOrganization(c *gin.Context) {
	userID := c.GetUint("userID")
	orgs, err := h.service.GetUserOrganizations(c.Request.Context(), userID)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "获取组织列表成功", orgs)
}

// @Summary 更新组织信息
// @Description 更新组织基本信息（仅管理员可操作）
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "组织ID"
// @Param request body dto.UpdateOrganizationRequest true "更新信息"
// @Success 200 {object} Response
// @Router /api/v1/organizations/{id} [put]
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := c.GetUint("userID")

	// 检查权限
	role, err := h.service.CheckMemberRole(c.Request.Context(), uint(id), userID)
	if err != nil || role != "admin" {
		SendError(c, http.StatusForbidden, "没有权限执行此操作")
		return
	}

	var req dto.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.service.UpdateOrganization(c.Request.Context(), uint(id), &req); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "更新组织信息成功", nil)
}

// @Summary 删除组织
// @Description 删除组织（仅管理员可操作）
// @Produce json
// @Security Bearer
// @Param id path int true "组织ID"
// @Success 200 {object} Response
// @Router /api/v1/organizations/{id} [delete]
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := c.GetUint("userID")

	// 检查权限
	role, err := h.service.CheckMemberRole(c.Request.Context(), uint(id), userID)
	if err != nil || role != "admin" {
		SendError(c, http.StatusForbidden, "没有权限执行此操作")
		return
	}

	if err := h.service.DeleteOrganization(c.Request.Context(), uint(id)); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "删除组织成功", nil)
}

// @Summary 加入组织
// @Description 加入已存在的组织
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.JoinOrganizationRequest true "组织ID"
// @Success 200 {object} Response
// @Router /api/v1/organizations/{id}/join [post]
func (h *OrganizationHandler) JoinOrganization(c *gin.Context) {
	var req dto.JoinOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	userID := c.GetUint("userID")
	if err := h.service.AddMember(c.Request.Context(), req.OrganizationID, &dto.AddMemberRequest{
		UserID: userID,
		Role:   "member",
	}); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "加入组织成功", nil)
}

// @Summary 添加组织成员
// @Description 添加新成员到组织（仅管理员可操作）
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "组织ID"
// @Param request body dto.AddMemberRequest true "成员信息"
// @Success 200 {object} Response
// @Router /api/v1/organizations/{id}/members [post]
func (h *OrganizationHandler) AddOrganizationMember(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := c.GetUint("userID")

	// 检查权限
	role, err := h.service.CheckMemberRole(c.Request.Context(), uint(id), userID)
	if err != nil || role != "admin" {
		SendError(c, http.StatusForbidden, "没有权限执行此操作")
		return
	}

	var req dto.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.service.AddMember(c.Request.Context(), uint(id), &req); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "添加成员成功", nil)
}

// @Summary 获取组织成员列表
// @Description 获取组织的所有成员
// @Produce json
// @Security Bearer
// @Param id path int true "组织ID"
// @Success 200 {object} Response{data=[]dto.MemberResponse}
// @Router /api/v1/organizations/{id}/members [get]
func (h *OrganizationHandler) GetOrganizationMembers(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	members, err := h.service.GetMembers(c.Request.Context(), uint(id))
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "获取成员列表成功", members)
}

// @Summary 更新组织成员
// @Description 更新组织成员信息（仅管理员可操作）
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "组织ID"
// @Param userid path int true "用户ID"
// @Param request body dto.UpdateMemberRequest true "更新信息"
// @Success 200 {object} Response
// @Router /api/v1/organizations/{id}/members/{userid} [put]
func (h *OrganizationHandler) UpdateOrganizationMember(c *gin.Context) {
	orgID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := c.GetUint("userID")
	targetUserID, _ := strconv.ParseUint(c.Param("userid"), 10, 64)

	// 检查权限
	role, err := h.service.CheckMemberRole(c.Request.Context(), uint(orgID), userID)
	if err != nil || role != "admin" {
		SendError(c, http.StatusForbidden, "没有权限执行此操作")
		return
	}

	var req dto.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.service.UpdateMember(c.Request.Context(), uint(orgID), uint(targetUserID), &req); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "更新成员信息成功", nil)
}

// @Summary 移除组织成员
// @Description 从组织中移除成员（仅管理员可操作）
// @Produce json
// @Security Bearer
// @Param id path int true "组织ID"
// @Param userid path int true "用户ID"
// @Success 200 {object} Response
// @Router /api/v1/organizations/{id}/members/{userid} [delete]
func (h *OrganizationHandler) RemoveOrganizationMember(c *gin.Context) {
	orgID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := c.GetUint("userID")
	targetUserID, _ := strconv.ParseUint(c.Param("userid"), 10, 64)

	// 检查权限
	role, err := h.service.CheckMemberRole(c.Request.Context(), uint(orgID), userID)
	if err != nil || role != "admin" {
		SendError(c, http.StatusForbidden, "没有权限执行此操作")
		return
	}

	if err := h.service.RemoveMember(c.Request.Context(), uint(orgID), uint(targetUserID)); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "移除成员成功", nil)
}

// @Summary 设置成员角色
// @Description 设置组织成员的角色（仅管理员可操作）
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "组织ID"
// @Param userid path int true "用户ID"
// @Param request body dto.UpdateMemberRequest true "角色信息"
// @Success 200 {object} Response
// @Router /api/v1/organizations/{id}/members/{userid}/role [put]
func (h *OrganizationHandler) SetOrganizationMemberRole(c *gin.Context) {
	orgID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := c.GetUint("userID")
	targetUserID, _ := strconv.ParseUint(c.Param("userid"), 10, 64)

	// 检查权限
	role, err := h.service.CheckMemberRole(c.Request.Context(), uint(orgID), userID)
	if err != nil || role != "admin" {
		SendError(c, http.StatusForbidden, "没有权限执行此操作")
		return
	}

	var req dto.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.service.UpdateMember(c.Request.Context(), uint(orgID), uint(targetUserID), &req); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "设置成员角色成功", nil)
}

// 其他处理方法继续...
