package handler

import (
	"ddup-apis/internal/dto"
	"ddup-apis/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	service     *service.OrganizationService
	userService *service.UserService
}

func NewOrganizationHandler(service *service.OrganizationService, userService *service.UserService) *OrganizationHandler {
	return &OrganizationHandler{service: service, userService: userService}
}

// @Tags 组织
// @Summary 创建组织
// @Description 创建新的组织，创建者默认成为管理员
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateOrganizationRequest true "组织信息"
// @Success 200 {object} Response
// @Router /api/v1/organizations [post]
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req dto.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	userID := c.GetUint("userID")
	err := h.service.CreateOrganization(c.Request.Context(), userID, &req)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "创建组织成功", nil)
}

// @Tags 组织
// @Summary 获取用户组织列表
// @Description 获取当前用户所属的所有组织
// @Produce json
// @Security Bearer
// @Success 200 {object} Response{data=[]dto.OrganizationResponse} "组织列表"
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

// @Tags 组织
// @Summary 更新组织信息
// @Description 更新组织基本信息（仅管理员可操作）
// @Accept json
// @Produce json
// @Security Bearer
// @Param org_name path string true "组织名称"
// @Param request body dto.UpdateOrganizationRequest true "更新信息"
// @Success 200 {object} Response
// @Router /api/v1/organizations/{org_name} [put]
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	orgName := c.Param("org_name")
	userID := c.GetUint("userID")

	// 验证组织名称格式
	if err := h.service.ValidateOrgName(orgName); err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 获取组织信息
	org, err := h.service.GetOrgByName(c.Request.Context(), orgName)
	if err != nil {
		SendError(c, http.StatusNotFound, "组织不存在")
		return
	}

	// 检查权限
	role, err := h.service.CheckMemberRole(c.Request.Context(), org.ID, userID)
	if err != nil || role != "admin" {
		SendError(c, http.StatusForbidden, "没有权限执行此操作")
		return
	}

	var req dto.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.service.UpdateOrganization(c.Request.Context(), org.ID, &req); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "更新组织信息成功", nil)
}

// @Tags 组织
// @Summary 删除组织
// @Description 删除组织（仅管理员可操作）
// @Produce json
// @Security Bearer
// @Param org_name path string true "组织名称"
// @Success 200 {object} Response
// @Router /api/v1/organizations/{org_name} [delete]
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	orgName := c.Param("org_name")
	userID := c.GetUint("userID")

	// 验证组织名称格式
	if err := h.service.ValidateOrgName(orgName); err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 获取组织信息
	org, err := h.service.GetOrgByName(c.Request.Context(), orgName)
	if err != nil {
		SendError(c, http.StatusNotFound, "组织不存在")
		return
	}

	// 检查权限
	role, err := h.service.CheckMemberRole(c.Request.Context(), org.ID, userID)
	if err != nil || role != "admin" {
		SendError(c, http.StatusForbidden, "没有权限执行此操作")
		return
	}

	if err := h.service.DeleteOrganization(c.Request.Context(), org.ID); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "删除组织成功", nil)
}

// @Tags 组织成员
// @Summary 加入组织
// @Description 加入已存在的组织
// @Accept json
// @Produce json
// @Security Bearer
// @Param org_name path string true "组织名称"
// @Success 200 {object} Response
// @Router /api/v1/organizations/{org_name}/join [post]
func (h *OrganizationHandler) JoinOrganization(c *gin.Context) {
	orgName := c.Param("org_name")
	username := c.GetString("username")

	// 验证组织名称格式
	if err := h.service.ValidateOrgName(orgName); err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 获取组织信息
	org, err := h.service.GetOrgByName(c.Request.Context(), orgName)
	if err != nil {
		SendError(c, http.StatusNotFound, "组织不存在")
		return
	}

	if err := h.service.AddMember(c.Request.Context(), org.ID, username, "member"); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "加入组织成功", nil)
}

// @Tags 组织成员
// @Summary 添加组织成员
// @Description 添加新成员到组织（仅管理员可操作）
// @Accept json
// @Produce json
// @Security Bearer
// @Param org_name path string true "组织名称"
// @Param request body dto.AddOrganizationMemberRequest true "成员信息"
// @Success 200 {object} Response
// @Router /api/v1/organizations/{org_name}/members [post]
func (h *OrganizationHandler) AddOrganizationMember(c *gin.Context) {
	orgName := c.Param("org_name")
	userID := c.GetUint("userID")

	// 验证组织名称格式
	if err := h.service.ValidateOrgName(orgName); err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 获取组织信息
	org, err := h.service.GetOrgByName(c.Request.Context(), orgName)
	if err != nil {
		SendError(c, http.StatusNotFound, "组织不存在")
		return
	}

	// 检查权限
	role, err := h.service.CheckMemberRole(c.Request.Context(), org.ID, userID)
	if err != nil || role != "admin" {
		SendError(c, http.StatusForbidden, "没有权限执行此操作")
		return
	}

	var req dto.AddOrganizationMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.service.AddMember(c.Request.Context(), org.ID, req.Username, req.Role); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "添加成员成功", nil)
}

// @Tags 组织成员
// @Summary 获取组织成员列表
// @Description 获取组织的所有成员
// @Produce json
// @Security Bearer
// @Param org_name path string true "组织名称"
// @Success 200 {object} Response{data=[]dto.MemberResponse}
// @Router /api/v1/organizations/{org_name}/members [get]
func (h *OrganizationHandler) GetOrganizationMembers(c *gin.Context) {
	orgName := c.Param("org_name")

	// 验证组织名称格式
	if err := h.service.ValidateOrgName(orgName); err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 获取组织信息
	org, err := h.service.GetOrgByName(c.Request.Context(), orgName)
	if err != nil {
		SendError(c, http.StatusNotFound, "组织不存在")
		return
	}

	members, err := h.service.GetMembers(c.Request.Context(), org.ID)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "获取成员列表成功", members)
}

// @Tags 组织成员
// @Summary 更新组织成员
// @Description 更新组织成员信息（仅管理员可操作）
// @Accept json
// @Produce json
// @Security Bearer
// @Param org_name path string true "组织名称"
// @Param username path string true "用户名"
// @Param request body dto.UpdateMemberRequest true "更新信息"
// @Success 200 {object} Response
// @Router /api/v1/organizations/{org_name}/members/{username} [put]
func (h *OrganizationHandler) UpdateOrganizationMember(c *gin.Context) {
	orgName := c.Param("org_name")
	username := c.Param("username")
	userID := c.GetUint("userID")

	// 验证组织名称格式
	if err := h.service.ValidateOrgName(orgName); err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 获取组织信息
	org, err := h.service.GetOrgByName(c.Request.Context(), orgName)
	if err != nil {
		SendError(c, http.StatusNotFound, "组织不存在")
		return
	}

	// 检查权限
	role, err := h.service.CheckMemberRole(c.Request.Context(), org.ID, userID)
	if err != nil || role != "admin" {
		SendError(c, http.StatusForbidden, "没有权限执行此操作")
		return
	}

	var req dto.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.service.UpdateMember(c.Request.Context(), org.ID, username, &req); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "更新成员信息成功", nil)
}

// @Tags 组织成员
// @Summary 移除组织成员
// @Description 从组织中移除成员（仅管理员可操作）
// @Produce json
// @Security Bearer
// @Param org_name path string true "组织名称"
// @Param username path string true "用户名"
// @Success 200 {object} Response
// @Router /api/v1/organizations/{org_name}/members/{username} [delete]
func (h *OrganizationHandler) RemoveOrganizationMember(c *gin.Context) {
	orgName := c.Param("org_name")
	username := c.Param("username")
	userID := c.GetUint("userID")

	// 验证组织名称格式
	if err := h.service.ValidateOrgName(orgName); err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 获取组织信息
	org, err := h.service.GetOrgByName(c.Request.Context(), orgName)
	if err != nil {
		SendError(c, http.StatusNotFound, "组织不存在")
		return
	}

	// 检查权限
	role, err := h.service.CheckMemberRole(c.Request.Context(), org.ID, userID)
	if err != nil || role != "admin" {
		SendError(c, http.StatusForbidden, "没有权限执行此操作")
		return
	}

	if err := h.service.RemoveMember(c.Request.Context(), org.ID, username); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "移除成员成功", nil)
}
