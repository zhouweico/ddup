package handler

import (
	"ddup-apis/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignupRequest 注册请求参数
type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// @Summary 用户登录
// @Description 用户登录
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} LoginResponse "登录成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Failure 401 {object} ErrorResponse "用户名或密码错误"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
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

// @Summary 用户注册
// @Description 用户注册
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} SignupResponse "注册成功"
// @Failure 400 {object} ErrorResponse "无效的请求参数"
// @Router /signup [post]
func (h *UserHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.userService.Signup(c.Request.Context(), req.Username, req.Password); err != nil {
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	sendSuccess(c, "注册成功", SignupResponse{
		UserInfo: User{
			Username: req.Username,
		},
	})
}

func sendSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

func sendError(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}
