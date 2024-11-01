package repository

import (
	"context"
	"ddup-apis/internal/model"

	"gorm.io/gorm"
)

type ISessionRepository interface {
	CreateSession(ctx context.Context, session *model.Session) error
	GetSessionByToken(ctx context.Context, token string) (*model.Session, error)
	InvalidateSession(ctx context.Context, token string) error
}

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) ISessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *model.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *SessionRepository) GetSessionByToken(ctx context.Context, token string) (*model.Session, error) {
	var session model.Session
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&session).Error
	return &session, err
}

func (r *SessionRepository) InvalidateSession(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&model.Session{}).
		Where("token = ?", token).
		Update("is_valid", false).Error
}
