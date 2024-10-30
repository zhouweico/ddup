package service

import (
	"context"
	"ddup-apis/internal/model"
	"ddup-apis/internal/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type IUserService interface {
	Register(ctx context.Context, username, password string) (string, error)
	Login(ctx context.Context, username, password string) (*LoginResult, error)
	ValidateToken(ctx context.Context, token string) (*TokenValidationResult, error)
	Logout(ctx context.Context, token string) error
	UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, userID uint) error
	GetUsers(ctx context.Context, page, pageSize int) ([]model.User, int64, error)
	GetUserByID(ctx context.Context, userID uint) (*model.User, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	GetUserByUUID(ctx context.Context, uuid string) (*model.User, error)
}

type LoginResult struct {
	Token     string
	CreatedAt time.Time
	ExpiresIn int64
	ExpiredAt time.Time
	User      *model.User
}

type TokenValidationResult struct {
	IsValid  bool
	UserID   uint
	Username string
	UUID     string
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Register(ctx context.Context, username, password string) (string, error) {
	var existingUser model.User
	result := s.db.Where("username = ?", username).First(&existingUser)

	if result.Error == nil {
		return "", errors.New("用户名已存在")
	}

	if result.Error != gorm.ErrRecordNotFound {
		return "", errors.New("系统错误")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}
	newUser := model.User{
		Username: username,
		Password: hashedPassword,
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return "", fmt.Errorf("创建用户失败: %w", err)
	}

	return newUser.UUID, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, fmt.Errorf("系统错误: %w", err)
	}

	if !utils.ComparePasswords(user.Password, password) {
		return nil, errors.New("用户名或密码错误")
	}

	token, createdAt, expiresIn, expiredAt, err := utils.GenerateToken(user.ID, user.Username, user.UUID)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	if err := s.db.Model(&user).Update("last_login", time.Now()).Error; err != nil {
		return nil, fmt.Errorf("更新登录时间失败: %w", err)
	}

	// 生成token后，保存会话信息
	session := model.UserSession{
		UserID:    user.ID,
		Token:     token,
		IsValid:   true,
		ExpiredAt: expiredAt,
	}

	if err := s.db.Create(&session).Error; err != nil {
		return nil, fmt.Errorf("保存会话信息失败: %w", err)
	}

	return &LoginResult{
		Token:     token,
		CreatedAt: createdAt,
		ExpiresIn: expiresIn,
		ExpiredAt: expiredAt,
		User:      &user,
	}, nil
}

func (s *UserService) ValidateToken(ctx context.Context, token string) (*TokenValidationResult, error) {
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

	// 获取用户信息
	var user model.User
	if err := s.db.First(&user, session.UserID).Error; err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	return &TokenValidationResult{
		IsValid:  true,
		UserID:   user.ID,
		Username: user.Username,
		UUID:     user.UUID,
	}, nil
}

func (s *UserService) Logout(ctx context.Context, token string) error {
	return s.db.Model(&model.UserSession{}).
		Where("token = ?", token).
		Update("is_valid", false).Error
}

func (s *UserService) GetUsers(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// 计算总数
	result := s.db.Model(&model.User{}).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// 分页查询
	offset := (page - 1) * pageSize
	result = s.db.Offset(offset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) error {
	delete(updates, "id")
	delete(updates, "password")
	delete(updates, "created_at")
	delete(updates, "deleted_at")

	result := s.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates)
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return result.Error
}

func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	result := s.db.Delete(&model.User{}, userID)
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return result.Error
}

func (s *UserService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !utils.ComparePasswords(user.Password, oldPassword) {
		return errors.New("原密码错误")
	}

	// 检查新密码是否与旧密码相同
	if oldPassword == newPassword {
		return errors.New("新密码不能与原密码相同")
	}

	// 生成新密码哈希
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	return s.db.Model(&user).Update("password", hashedPassword).Error
}

func (s *UserService) GetUserByUUID(ctx context.Context, uuid string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}
