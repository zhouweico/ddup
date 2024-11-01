package repository

import (
	"context"

	"ddup-apis/internal/model"

	"gorm.io/gorm"
)

type ISocialRepository interface {
	Create(ctx context.Context, social *model.Social) error
	GetByUserID(ctx context.Context, userID uint) ([]model.Social, error)
	Delete(ctx context.Context, userID uint, id string) error
	Update(ctx context.Context, userID uint, id string, social *model.Social) error
}

type SocialRepository struct {
	db *gorm.DB
}

func NewSocialRepository(db *gorm.DB) *SocialRepository {
	return &SocialRepository{
		db: db,
	}
}

func (r *SocialRepository) Delete(ctx context.Context, userID uint, id string) error {
	return r.db.Where("user_id = ? AND id = ?", userID, id).Delete(&model.Social{}).Error
}

func (r *SocialRepository) Create(ctx context.Context, social *model.Social) error {
	return r.db.Create(social).Error
}

func (r *SocialRepository) GetByUserID(ctx context.Context, userID uint) ([]model.Social, error) {
	var s []model.Social
	err := r.db.Where("user_id = ?", userID).Find(&s).Error
	return s, err
}

func (r *SocialRepository) Update(ctx context.Context, userID uint, id string, social *model.Social) error {
	return r.db.Where("user_id = ? AND id = ?", userID, id).Updates(social).Error
}
