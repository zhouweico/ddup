package service

import (
	"context"
	"ddup-apis/internal/model"
	"ddup-apis/internal/utils"
	"errors"
	"fmt"
	"log"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
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
	GetUserByUserID(ctx context.Context, userID string) (*model.User, error)
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
	ID       uint
	Username string
	UserID   string
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

	// 生成指定长度的 Nano ID (21位)
	id, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 21)
	if err != nil {
		log.Printf("生成ID失败: %v", err)
		return "", fmt.Errorf("生成ID失败: %w", err)
	}
	log.Printf("生成的UserID: %s, 长度: %d", id, len(id))

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		return "", fmt.Errorf("密码加密失败: %w", err)
	}

	newUser := model.User{
		Username: username,
		Password: hashedPassword,
		UserID:   id,
	}

	log.Printf("准备保存到数据库的UserID: %s", newUser.UserID)

	if err := s.db.Create(&newUser).Error; err != nil {
		log.Printf("创建用户失败: %v", err)
		return "", fmt.Errorf("创建用户失败: %w", err)
	}

	log.Printf("保存后的UserID: %s", newUser.UserID)

	return newUser.UserID, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	if !utils.ComparePasswords(user.Password, password) {
		return nil, errors.New("密码错误")
	}

	token, createdAt, expiresIn, expiredAt, err := utils.GenerateToken(user.ID, user.Username, user.UserID)
	if err != nil {
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

	var user model.User
	if err := s.db.First(&user, session.UserID).Error; err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	return &TokenValidationResult{
		IsValid:  true,
		ID:       user.ID,
		Username: user.Username,
		UserID:   user.UserID,
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

	result := s.db.Model(&model.User{}).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

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

func (s *UserService) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error {
	delete(updates, "id")
	delete(updates, "password")
	delete(updates, "created_at")
	delete(updates, "deleted_at")

	result := s.db.Model(&model.User{}).Where("id = ?", id).Updates(updates)
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

	if !utils.ComparePasswords(user.Password, oldPassword) {
		return errors.New("原密码错误")
	}

	if oldPassword == newPassword {
		return errors.New("新密码不能与原密码相同")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	return s.db.Model(&user).Update("password", hashedPassword).Error
}

func (s *UserService) GetUserByUserID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}
