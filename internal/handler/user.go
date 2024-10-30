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

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Login 处理用户登录
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

// Signup 处理用户注册
func (h *UserHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.userService.Signup(c.Request.Context(), req.Username, req.Password, req.Email); err != nil {
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
