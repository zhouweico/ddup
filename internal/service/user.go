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

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
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

func (s *userService) Signup(ctx context.Context, username, password string) error {
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

func (s *userService) ValidateToken(ctx context.Context, token string) (*TokenValidationResult, error) {
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
