package handler

import (
	"ddup-apis/internal/db"
	"ddup-apis/internal/model"
	"ddup-apis/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body LoginRequest true "登录信息"
// @Success 200 {object} Response{data=LoginResponse} "登录成功"
// @Failure 400 {object} Response "请求错误"
// @Failure 401 {object} Response "未授权"
// @Router /login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 查询用户
	var user model.User
	if err := db.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if !utils.ComparePasswords(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成 Token
	token, createdAt, expiresIn, expiredAt, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 Token 失败"})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "登录成功",
		Data: LoginResponse{
			TokenInfo: TokenInfo{
				Token:     token,
				CreatedAt: createdAt,
				ExpiresIn: expiresIn,
				ExpiredAt: expiredAt,
			},
			UserInfo: User{
				Username: user.Username,
			},
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
func Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 同时检查用户名和邮箱是否已存在
	var existingUser model.User
	result := db.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser)

	// 只有在找到记录时才检查重复
	if result.Error == nil {
		if existingUser.Username == req.Username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已被使用"})
		}
		return
	}

	// 如果是 "record not found" 错误，则继续创建用户
	if result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
		return
	}

	// 创建新用户
	hashedPassword := utils.HashPassword(req.Password)
	newUser := model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if err := db.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "注册成功",
		Data: SignupResponse{
			UserInfo: User{
				Username: newUser.Username,
			},
		},
	})
}
