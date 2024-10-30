package handler

import (
	"ddup-apis/internal/service"
	"ddup-apis/internal/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SignupRequest 注册请求参数
type SignupRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// @Summary 用户注册
// @Description 用户注册
// @Accept json
// @Produce json
// @Param request body SignupRequest true "注册信息"
// @Success 200 {object} response.Response{data=SignupResponse} "注册成功"
// @Failure 400 {object} response.Response "无效的请求参数"
// @Router /sign-up [post]
func (h *UserHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	if err := h.userService.Signup(c.Request.Context(), req.Username, req.Password); err != nil {
		response.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SendSuccess(c, "注册成功", SignupResponse{
		UserInfo: User{
			Username: req.Username,
		},
	})
}

// @Summary 用户登录
// @Description 用户登录
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=LoginResponse} "登录成功"
// @Failure 400 {object} response.Response "无效的请求参数"
// @Failure 401 {object} response.Response "用户名或密码错误"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendError(c, http.StatusBadRequest, "无效的请求参数")
		return
	}

	result, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.SendError(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	response.SendSuccess(c, "登录成功", LoginResponse{
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
