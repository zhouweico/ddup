package repository

import (
	"context"
	"ddup-apis/internal/model"

	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Create(ctx context.Context, profile *model.Profile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *ProfileRepository) GetByID(ctx context.Context, id uint) (*model.Profile, error) {
	var profile model.Profile
	err := r.db.WithContext(ctx).First(&profile, id).Error
	return &profile, err
}

func (r *ProfileRepository) GetByUserID(ctx context.Context, userID uint) ([]model.Profile, error) {
	var profiles []model.Profile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("display_order asc").Find(&profiles).Error
	return profiles, err
}

func (r *ProfileRepository) GetByType(ctx context.Context, userID uint, profileType string) ([]model.Profile, error) {
	var profiles []model.Profile
	err := r.db.WithContext(ctx).Where("user_id = ? AND type = ?", userID, profileType).
		Order("display_order asc").Find(&profiles).Error
	return profiles, err
}

func (r *ProfileRepository) Update(ctx context.Context, profile *model.Profile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

func (r *ProfileRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Profile{}, id).Error
}

func (r *ProfileRepository) UpdateDisplayOrder(ctx context.Context, items []struct {
	ID    uint
	Order int
}) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&model.Profile{}).Where("id = ?", item.ID).
				Update("display_order", item.Order).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
