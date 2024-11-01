package repository

import (
	"context"
	"ddup-apis/internal/model"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(ctx context.Context, org *model.Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id uint) (*model.Organization, error) {
	var org model.Organization
	err := r.db.WithContext(ctx).First(&org, id).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrganizationRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Organization{}).Where("id = ?", id).Updates(updates).Error
}

func (r *OrganizationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Organization{}, id).Error
}

func (r *OrganizationRepository) GetUserOrganizations(ctx context.Context, userID uint) ([]model.Organization, error) {
	var orgs []model.Organization
	err := r.db.WithContext(ctx).
		Joins("JOIN organization_members ON organizations.id = organization_members.organization_id").
		Where("organization_members.user_id = ?", userID).
		Find(&orgs).Error
	return orgs, err
}

func (r *OrganizationRepository) AddMember(ctx context.Context, member *model.OrganizationMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *OrganizationRepository) GetMember(ctx context.Context, orgID, userID uint) (*model.OrganizationMember, error) {
	var member model.OrganizationMember
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *OrganizationRepository) UpdateMember(ctx context.Context, orgID, userID uint, role string) error {
	return r.db.WithContext(ctx).
		Model(&model.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Update("role", role).Error
}

func (r *OrganizationRepository) RemoveMember(ctx context.Context, orgID, userID uint) error {
	return r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Delete(&model.OrganizationMember{}).Error
}

func (r *OrganizationRepository) GetMembers(ctx context.Context, orgID uint) ([]model.OrganizationMember, error) {
	var members []model.OrganizationMember
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Find(&members).Error
	return members, err
}
