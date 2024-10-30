package service

import (
	"context"
	"time"

	"ddup-apis/internal/model"

	"gorm.io/gorm"
)

type UserService interface {
	Signup(ctx context.Context, username, password, email string) error
	Login(ctx context.Context, username, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
}

type userService struct {
	db        *gorm.DB
	jwtSecret string
}

// 验证 token
func (s *userService) ValidateToken(ctx context.Context, token string) (bool, error) {
	var session model.UserSession
	if err := s.db.Where("token = ?", token).First(&session).Error; err != nil {
		return false, err
	}

	// 检查 token 是否有效
	if !session.IsValid || time.Now().After(session.ExpiredAt) {
		// 标记 token 为失效
		if err := s.db.Model(&model.UserSession{}).Where("token = ?", token).Update("is_valid", false).Error; err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

// 修复会话相关的数据库操作
func (s *userService) GetUserSessionByToken(token string) (*model.UserSession, error) {
	var session model.UserSession
	if err := s.db.Where("token = ?", token).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *userService) InvalidateUserSession(token string) error {
	return s.db.Where("token = ?", token).Delete(&model.UserSession{}).Error
}
