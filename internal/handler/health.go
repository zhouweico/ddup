package handler

import (
	"ddup-apis/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check godoc
// @Summary 健康检查
// @Description 检查服务和数据库连接状态
// @Tags 系统
// @Produce json
// @Success 200 {object} Response "服务正常"
// @Failure 503 {object} Response "服务异常"
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	if err := db.Ping(); err != nil {
		SendError(c, http.StatusServiceUnavailable, "数据库连接异常")
		return
	}

	SendSuccess(c, "服务正常", nil)
}
