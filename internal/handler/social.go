package handler

import (
	"ddup-apis/internal/model"
	"ddup-apis/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SocialRequest 创建/更新社交媒体账号请求
type SocialRequest struct {
	Platform    string `json:"platform" binding:"required"` // 平台名称
	Username    string `json:"username" binding:"required"` // 平台用户名
	URL         string `json:"url"`                         // 个人主页链接
	Description string `json:"description"`                 // 描述
}

// SocialResponse 社交媒体账号响应
type SocialResponse struct {
	ID          uint   `json:"id"`
	Platform    string `json:"platform"`
	Username    string `json:"username"`
	URL         string `json:"url"`
	Verified    bool   `json:"verified"`
	Description string `json:"description"`
}

type SocialHandler struct {
	service *service.SocialService
}

func NewSocialHandler(service *service.SocialService) *SocialHandler {
	return &SocialHandler{
		service: service,
	}
}

// CreateSocial 创建社交媒体账号
// @Summary 创建社交媒体账号
// @Description 为当前用户创建社交媒体账号
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body SocialRequest true "社交媒体账号信息"
// @Success 200 {object} Response
// @Router /api/v1/users/socials [post]
func (h *SocialHandler) CreateSocial(c *gin.Context) {
	var req SocialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, http.StatusUnauthorized, "未授权")
		return
	}

	social := &model.Social{
		UserID:      userID.(uint),
		Platform:    req.Platform,
		Username:    req.Username,
		URL:         req.URL,
		Description: req.Description,
	}

	if err := h.service.Create(c.Request.Context(), userID.(uint), social); err != nil {
		SendError(c, http.StatusInternalServerError, "创建失败")
		return
	}

	SendSuccess(c, "创建成功", nil)
}

// GetUserSocial 获取用户的社交媒体账号列表
// @Summary 获取社交媒体账号列表
// @Description 获取当前用户的所有社交媒体账号
// @Produce json
// @Security Bearer
// @Success 200 {object} Response{data=[]SocialResponse}
// @Router /api/v1/users/socials [get]
func (h *SocialHandler) GetUserSocial(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, http.StatusUnauthorized, "未授权")
		return
	}
	s, err := h.service.GetByUserID(c.Request.Context(), userID.(uint))
	if err != nil {
		SendError(c, http.StatusInternalServerError, "获取失败")
		return
	}

	var resp []SocialResponse
	for _, m := range s {
		resp = append(resp, SocialResponse{
			ID:          m.ID,
			Platform:    m.Platform,
			Username:    m.Username,
			URL:         m.URL,
			Verified:    m.Verified,
			Description: m.Description,
		})
	}

	SendSuccess(c, "获取成功", resp)
}

// UpdateSocial 更新社交媒体账号
// @Summary 更新社交媒体账号
// @Description 更新指定的社交媒体账号
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "社交媒体账号ID"
// @Param request body SocialRequest true "社交媒体账号信息"
// @Success 200 {object} Response
// @Router /api/v1/users/socials/{id} [put]
func (h *SocialHandler) UpdateSocial(c *gin.Context) {
	var req SocialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, http.StatusUnauthorized, "未授权")
		return
	}

	id := c.Param("id")

	if err := h.service.Update(c.Request.Context(), userID.(uint), id, &model.Social{
		Platform:    req.Platform,
		Username:    req.Username,
		URL:         req.URL,
		Description: req.Description,
	}); err != nil {
		SendError(c, http.StatusInternalServerError, "更新失败")
		return
	}

	SendSuccess(c, "更新成功", nil)
}

// DeleteSocial 删除社交媒体账号
// @Summary 删除社交媒体账号
// @Description 删除指定的社交媒体账号
// @Produce json
// @Security Bearer
// @Param id path string true "社交媒体账号ID"
// @Success 200 {object} Response
// @Router /api/v1/users/socials/{id} [delete]
func (h *SocialHandler) DeleteSocial(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, http.StatusUnauthorized, "未授权")
		return
	}

	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), userID.(uint), id); err != nil {
		SendError(c, http.StatusInternalServerError, "删除失败")
		return
	}

	SendSuccess(c, "删除成功", nil)
}
