package handler

import (
	"ddup-apis/internal/dto"
	"ddup-apis/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	service *service.ProfileService
}

func NewProfileHandler(service *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{service: service}
}

// @Tags 个人资料
// @Summary 创建个人资料项
// @Description 创建新的个人资料项
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateProfileRequest true "个人资料信息"
// @Success 200 {object} Response
// @Router /api/v1/profiles [post]
func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	var req dto.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	userID := c.GetUint("userID")
	if err := h.service.Create(c.Request.Context(), userID, &req); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "创建成功", nil)
}

// @Tags 个人资料
// @Summary 获取个人资料列表
// @Description 获取指定类型的个人资料列表
// @Produce json
// @Security Bearer
// @Param type query string true "资料类型"
// @Success 200 {object} Response{data=[]dto.ProfileResponse}
// @Router /api/v1/profiles [get]
func (h *ProfileHandler) GetProfiles(c *gin.Context) {
	profileType := c.Query("type")
	userID := c.GetUint("userID")

	profiles, err := h.service.GetByType(c.Request.Context(), userID, profileType)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "获取成功", profiles)
}

// @Tags 个人资料
// @Summary 更新个人资料
// @Description 更新指定的个人资料项
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path uint true "资料ID"
// @Param request body dto.UpdateProfileRequest true "更新信息"
// @Success 200 {object} Response
// @Router /api/v1/profiles/{id} [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	userID := c.GetUint("userID")
	profileID, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		SendError(c, http.StatusBadRequest, "无效的ID参数")
		return
	}

	if err := h.service.Update(c.Request.Context(), userID, uint(profileID), &req); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "更新成功", nil)
}

// @Tags 个人资料
// @Summary 删除个人资料
// @Description 删除指定的个人资料项
// @Produce json
// @Security Bearer
// @Param id path uint true "资料ID"
// @Success 200 {object} Response
// @Router /api/v1/profiles/{id} [delete]
func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	userID := c.GetUint("userID")
	profileID, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		SendError(c, http.StatusBadRequest, "无效的ID参数")
		return
	}

	if err := h.service.Delete(c.Request.Context(), userID, uint(profileID)); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "删除成功", nil)
}

// @Tags 个人资料
// @Summary 更新显示顺序
// @Description 更新个人资料项的显示顺序
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.UpdateDisplayOrderRequest true "顺序信息"
// @Success 200 {object} Response
// @Router /api/v1/profiles/display-order [put]
func (h *ProfileHandler) UpdateDisplayOrder(c *gin.Context) {
	var req dto.UpdateDisplayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	userID := c.GetUint("userID")
	if err := h.service.UpdateDisplayOrder(c.Request.Context(), userID, &req); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccess(c, "更新显示顺序成功", nil)
}
