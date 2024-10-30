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

type UserService interface {
	Signup(ctx context.Context, username, password string) error
	Login(ctx context.Context, username, password string) (*LoginResult, error)
	ValidateToken(ctx context.Context, token string) (*TokenValidationResult, error)
	Logout(ctx context.Context, token string) error
	UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, userID uint) error
	GetUsers(ctx context.Context, page, pageSize int) ([]model.User, int64, error)
	GetUserByID(ctx context.Context, userID uint) (*model.User, error)
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
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
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

	token, createdAt, expiresIn, expiredAt, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	return &LoginResult{
		Token:     token,
		CreatedAt: createdAt,
		ExpiresIn: expiresIn,
		ExpiredAt: expiredAt,
		User:      &user,
	}, nil
}

func (s *UserService) Signup(ctx context.Context, username, password string) error {
	var existingUser model.User
	result := s.db.Where("username = ?", username).First(&existingUser)

	if result.Error == nil {
		return errors.New("用户名已存在")
	}

	if result.Error != gorm.ErrRecordNotFound {
		return errors.New("系统错误")
	}

	hashedPassword := utils.HashPassword(password)
	newUser := model.User{
		Username: username,
		Password: hashedPassword,
	}

	return s.db.Create(&newUser).Error
}

func (s *UserService) ValidateToken(ctx context.Context, token string) (*TokenValidationResult, error) {
	var session model.UserSession
	if err := s.db.Where("token = ?", token).First(&session).Error; err != nil {
		return nil, err
	}

	if !session.IsValid || time.Now().After(session.ExpiredAt) {
		if err := s.db.Model(&model.UserSession{}).
			Where("token = ?", token).
			Update("is_valid", false).Error; err != nil {
			return nil, err
		}
		return &TokenValidationResult{IsValid: false}, nil
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	return &TokenValidationResult{
		IsValid:  true,
		UserID:   claims.UserID,
		Username: claims.Username,
	}, nil
}

func (s *UserService) Logout(ctx context.Context, token string) error {
	return s.db.Model(&model.UserSession{}).
		Where("token = ?", token).
		Update("is_valid", false).Error
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
