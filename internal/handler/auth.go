package handler

import (
	"ddup-apis/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService service.UserService
}

func NewHandler(userService service.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	result, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		sendError(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	sendSuccess(c, "登录成功", LoginResponse{
		TokenInfo: TokenInfo{
			Token:     result.Token,
			CreatedAt: result.CreatedAt,
			ExpiresIn: result.ExpiresIn,
			ExpiredAt: result.ExpiredAt,
		},
		UserInfo: User{
			Username: result.User.Username,
		},
	})
}

type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// @Summary 用户注册
// @Description 用户注册接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body SignupRequest true "注册信息"
// @Success 200 {object} Response{data=SignupResponse} "注册成功"
// @Failure 400 {object} Response "请求错误"
// @Router /sign-up [post]
func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	err := h.userService.Signup(c.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		switch err.Error() {
		case "用户名已存在", "邮箱已被使用":
			sendError(c, http.StatusBadRequest, err.Error())
		default:
			sendError(c, http.StatusInternalServerError, "注册失败")
		}
		return
	}

	sendSuccess(c, "注册成功", SignupResponse{
		UserInfo: User{
			Username: req.Username,
		},
	})
}

func sendError(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Code:    status,
		Message: message,
	})
}
