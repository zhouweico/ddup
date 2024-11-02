package repository

import (
	"context"
	"ddup-apis/internal/errors"
	"ddup-apis/internal/model"
	stderrors "errors"

	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// 组织相关方法
// GetByName 通过名称获取组织
func (r *OrganizationRepository) GetByName(ctx context.Context, name string) (*model.Organization, error) {
	var org model.Organization
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&org).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(404, "组织不存在", err)
		}
		return nil, errors.Wrap(err, "查询组织失败")
	}
	return &org, nil
}

// Create 创建组织
func (r *OrganizationRepository) Create(ctx context.Context, org *model.Organization) error {
	// 检查名称是否已存在
	exists, err := r.IsNameExists(ctx, org.Name)
	if err != nil {
		return err
	}
	if exists {
		return errors.New(400, "组织名称已存在", nil)
	}

	return r.db.WithContext(ctx).Create(org).Error
}

// IsNameExists 检查组织名称是否存在
func (r *OrganizationRepository) IsNameExists(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Organization{}).
		Where("name = ?", name).
		Count(&count).Error
	return count > 0, err
}

// Update 更新组织信息
func (r *OrganizationRepository) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 如果要更新名称，先检查是否存在
	if name, ok := updates["name"]; ok {
		exists, err := r.IsNameExists(ctx, name.(string))
		if err != nil {
			return err
		}
		if exists {
			return errors.New(400, "组织名称已存在", nil)
		}
	}

	return r.db.WithContext(ctx).Model(&model.Organization{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// Delete 删除组织
func (r *OrganizationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Organization{}, id).Error
}

// 成员相关方法
// GetMember 获取成员信息
func (r *OrganizationRepository) GetMember(ctx context.Context, orgID, userID uint) (*model.OrganizationMember, error) {
	var member model.OrganizationMember
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		First(&member).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(404, "成员不存在", err)
		}
		return nil, errors.Wrap(err, "查询成员失败")
	}
	return &member, nil
}

// GetMembers 获取组织所有成员
func (r *OrganizationRepository) GetMembers(ctx context.Context, orgID uint) ([]model.OrganizationMember, error) {
	var members []model.OrganizationMember
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Find(&members).Error
	return members, err
}

// GetMembersByRole 获取指定角色的成员
func (r *OrganizationRepository) GetMembersByRole(ctx context.Context, orgID uint, role string) ([]model.OrganizationMember, error) {
	var members []model.OrganizationMember
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND role = ?", orgID, role).
		Find(&members).Error
	return members, err
}

// IsMember 检查用户是否是组织成员
func (r *OrganizationRepository) IsMember(ctx context.Context, orgID, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Count(&count).Error
	return count > 0, err
}

// AddMember 添加成员
func (r *OrganizationRepository) AddMember(ctx context.Context, member *model.OrganizationMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

// UpdateMember 更新成员信息
func (r *OrganizationRepository) UpdateMember(ctx context.Context, orgID, userID uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Updates(updates).Error
}

// RemoveMember 移除成员
func (r *OrganizationRepository) RemoveMember(ctx context.Context, orgID, userID uint) error {
	return r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Delete(&model.OrganizationMember{}).Error
}

// GetUserOrganizations 获取用户所属的所有组织
func (r *OrganizationRepository) GetUserOrganizations(ctx context.Context, userID uint) ([]model.Organization, error) {
	var orgs []model.Organization
	err := r.db.WithContext(ctx).
		Joins("JOIN organization_members ON organizations.id = organization_members.organization_id").
		Where("organization_members.user_id = ?", userID).
		Find(&orgs).Error
	return orgs, err
}

func (r *OrganizationRepository) WithTransaction(tx *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: tx}
}

func (r *OrganizationRepository) DB() *gorm.DB {
	return r.db
}
