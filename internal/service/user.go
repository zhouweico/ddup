package service

import (
	"context"
	"ddup-apis/internal/logger"
	"ddup-apis/internal/model"
	"ddup-apis/internal/utils"
	"fmt"
	"time"

	"ddup-apis/internal/errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IUserService interface {
	Register(ctx context.Context, username, password string) error
	Login(ctx context.Context, username, password string) (*LoginResult, error)
	ValidateToken(token string) (*TokenValidationResult, error)
	Logout(ctx context.Context, token string) error
	GetUserByID(ctx context.Context, userID uint) (*model.User, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, userID uint) error
}

type LoginResult struct {
	Token     string
	CreatedAt time.Time
	ExpiresIn int64
	ExpiredAt time.Time
	User      *model.User
}

type TokenValidationResult struct {
	UserID   uint
	Username string
	IsValid  bool
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Register(ctx context.Context, username, password string) error {
	var existingUser model.User
	result := s.db.Where("username = ?", username).First(&existingUser)

	if result.Error == nil {
		return errors.New(400, "用户名已存在", nil)
	}

	if result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("查询用户失败: %w", result.Error)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		logger.Error("密码加密失败", zap.Error(err))
		return fmt.Errorf("密码加密失败: %w", err)
	}

	newUser := model.User{
		Username: username,
		Password: hashedPassword,
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		logger.Error("创建用户失败", zap.Error(err))
		return fmt.Errorf("创建用户失败: %w", err)
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(404, "用户不存在", nil)
		}
		return nil, err
	}

	if !utils.ComparePasswords(user.Password, password) {
		return nil, errors.New(401, "密码错误", nil)
	}

	token, createdAt, expiresIn, expiredAt, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		logger.Error("生成令牌失败", zap.Error(err))
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	s.db.Model(&user).Update("last_login", time.Now())

	return &LoginResult{
		Token:     token,
		User:      &user,
		CreatedAt: createdAt,
		ExpiresIn: expiresIn,
		ExpiredAt: expiredAt,
	}, nil
}

func (s *UserService) ValidateToken(token string) (*TokenValidationResult, error) {
	var session model.UserSession
	if err := s.db.Where("token = ?", token).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &TokenValidationResult{IsValid: false}, nil
		}
		return nil, fmt.Errorf("验证token失败: %w", err)
	}

	if !session.IsValid || time.Now().After(session.ExpiredAt) {
		if err := s.db.Model(&session).Update("is_valid", false).Error; err != nil {
			return nil, fmt.Errorf("更新会话状态失败: %w", err)
		}
		return &TokenValidationResult{IsValid: false}, nil
	}

	var user model.User
	if err := s.db.First(&user, session.UserID).Error; err != nil {
		return nil, errors.New(404, "获取用户信息失败", err)
	}

	return &TokenValidationResult{
		UserID:   user.ID,
		Username: user.Username,
		IsValid:  true,
	}, nil
}

func (s *UserService) Logout(ctx context.Context, token string) error {
	return s.db.Model(&model.UserSession{}).
		Where("token = ?", token).
		Update("is_valid", false).Error
}

func (s *UserService) GetUserByID(ctx context.Context, userID uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(404, "用户不存在", nil)
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error {
	delete(updates, "id")
	delete(updates, "password")
	delete(updates, "created_at")
	delete(updates, "deleted_at")

	result := s.db.Model(&model.User{}).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		return errors.New(404, "用户不存在", nil)
	}
	return result.Error
}

func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	result := s.db.Delete(&model.User{}, userID)
	if result.RowsAffected == 0 {
		return errors.New(404, "用户不存在", nil)
	}
	return result.Error
}

func (s *UserService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New(404, "用户不存在", nil)
	}

	if !utils.ComparePasswords(user.Password, oldPassword) {
		return errors.New(401, "原密码错误", nil)
	}

	if oldPassword == newPassword {
		return errors.New(400, "新密码不能与原密码相同", nil)
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	return s.db.Model(&user).Update("password", hashedPassword).Error
}
