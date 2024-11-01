package service

import (
	"context"
	"ddup-apis/internal/model"
	"ddup-apis/internal/repository"

	"gorm.io/gorm"
)

type ISocialRepository interface {
	Create(ctx context.Context, social *model.Social) error
	GetByUserID(ctx context.Context, userID uint) ([]model.Social, error)
	Delete(ctx context.Context, userID uint, id string) error
	Update(ctx context.Context, userID uint, id string, social *model.Social) error
}

type ISocialService interface {
	Create(ctx context.Context, userID uint, social *model.Social) error
	GetByUserID(ctx context.Context, userID uint) ([]model.Social, error)
	Delete(ctx context.Context, userID uint, id string) error
	Update(ctx context.Context, userID uint, id string, social *model.Social) error
}

type SocialService struct {
	repo ISocialRepository
}

func NewSocialService(db *gorm.DB) *SocialService {
	return &SocialService{
		repo: repository.NewSocialRepository(db),
	}
}

func (s *SocialService) Create(ctx context.Context, userID uint, social *model.Social) error {
	social.UserID = userID
	return s.repo.Create(ctx, social)
}

func (s *SocialService) GetByUserID(ctx context.Context, userID uint) ([]model.Social, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *SocialService) Delete(ctx context.Context, userID uint, id string) error {
	return s.repo.Delete(ctx, userID, id)
}

func (s *SocialService) Update(ctx context.Context, userID uint, id string, social *model.Social) error {
	return s.repo.Update(ctx, userID, id, social)
}
