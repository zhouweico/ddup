package service

import (
	"context"
	"ddup-apis/internal/dto"
	"ddup-apis/internal/model"
	"ddup-apis/internal/repository"
	"fmt"
	"regexp"
	"strings"

	"ddup-apis/internal/errors"

	"gorm.io/gorm"
)

type OrganizationService struct {
	orgRepo  *repository.OrganizationRepository
	userRepo repository.IUserRepository
}

func NewOrganizationService(db *gorm.DB) *OrganizationService {
	return &OrganizationService{
		orgRepo:  repository.NewOrganizationRepository(db),
		userRepo: repository.NewUserRepository(db),
	}
}

// CreateOrganization 创建组织
func (s *OrganizationService) CreateOrganization(ctx context.Context, userID uint, req *dto.CreateOrganizationRequest) error {
	// 验证组织名称
	if err := s.ValidateOrgName(req.Name); err != nil {
		return err
	}

	org := &model.Organization{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Avatar:      req.Avatar,
		Description: req.Description,
		Location:    req.Location,
		Website:     req.Website,
	}

	// 使用事务确保创建组织和添加管理员成员是原子操作
	err := s.orgRepo.DB().Transaction(func(tx *gorm.DB) error {
		txOrgRepo := s.orgRepo.WithTransaction(tx)
		// 创建组织
		if err := txOrgRepo.Create(ctx, org); err != nil {
			return err
		}

		// 添加创建者为管理员
		member := &model.OrganizationMember{
			OrganizationID: org.ID,
			UserID:         userID,
			Role:           "admin",
		}
		return txOrgRepo.AddMember(ctx, member)
	})

	return err
}

// GetUserOrganizations 获取用户的组织列表
func (s *OrganizationService) GetUserOrganizations(ctx context.Context, userID uint) ([]dto.OrganizationResponse, error) {
	orgs, err := s.orgRepo.GetUserOrganizations(ctx, userID)
	if err != nil {
		return nil, err
	}

	var resp []dto.OrganizationResponse
	for _, org := range orgs {
		// 获取用户在组织中的角色
		member, err := s.orgRepo.GetMember(ctx, org.ID, userID)
		if err != nil {
			continue
		}

		resp = append(resp, dto.OrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
			DisplayName: org.DisplayName,
			Email:       org.Email,
			Avatar:      org.Avatar,
			Description: org.Description,
			Location:    org.Location,
			Website:     org.Website,
			Role:        member.Role,
		})
	}

	return resp, nil
}

// UpdateOrganization 更新组织信息
func (s *OrganizationService) UpdateOrganization(ctx context.Context, orgID uint, req *dto.UpdateOrganizationRequest) error {
	if req.Name != "" {
		if err := s.ValidateOrgName(req.Name); err != nil {
			return err
		}
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.DisplayName != "" {
		updates["display_name"] = req.DisplayName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Website != "" {
		updates["website"] = req.Website
	}

	return s.orgRepo.Update(ctx, orgID, updates)
}

// DeleteOrganization 删除组织
func (s *OrganizationService) DeleteOrganization(ctx context.Context, orgID uint) error {
	return s.orgRepo.Delete(ctx, orgID)
}

// GetMembers 获取组织成员列表
func (s *OrganizationService) GetMembers(ctx context.Context, orgID uint) ([]dto.MemberResponse, error) {
	members, err := s.orgRepo.GetMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var resp []dto.MemberResponse
	for _, member := range members {
		user, err := s.userRepo.GetByID(ctx, member.UserID)
		if err != nil {
			continue
		}

		resp = append(resp, dto.MemberResponse{
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Avatar:   user.Avatar,
			Role:     member.Role,
			JoinedAt: member.CreatedAt,
		})
	}

	return resp, nil
}

// AddMember 添加成员
func (s *OrganizationService) AddMember(ctx context.Context, orgID uint, username string, role string) error {
	// 获取用户信息
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return errors.New(404, "用户不存在", err)
	}

	// 检查是否已经是成员
	exists, err := s.orgRepo.IsMember(ctx, orgID, user.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New(400, "用户已经是组织成员", nil)
	}

	member := &model.OrganizationMember{
		OrganizationID: orgID,
		UserID:         user.ID,
		Role:           role,
	}
	return s.orgRepo.AddMember(ctx, member)
}

// RemoveMember 从组织中移除成员
func (s *OrganizationService) RemoveMember(ctx context.Context, orgID uint, username string) error {
	// 根据用户名获取用户ID
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return err
	}

	// 从组织成员表中删除记录
	return s.orgRepo.RemoveMember(ctx, orgID, user.ID)
}

func (s *OrganizationService) ValidateOrgName(name string) error {
	// 检查组织名称格式
	pattern := `^[a-z0-9][a-z0-9-]{0,38}[a-z0-9]$`
	matched, _ := regexp.MatchString(pattern, name)
	if !matched {
		return errors.New(400, "组织名称格式不正确，只能包含小写字母、数字和中划线，长度在2-40之间", nil)
	}

	// 检查保留名称
	reservedNames := []string{"admin", "api", "system", "root", "support", "help", "about"}
	for _, reserved := range reservedNames {
		if strings.ToLower(name) == reserved {
			return errors.New(400, "该组织名称已被系统保留", nil)
		}
	}

	return nil
}

func (s *OrganizationService) GetOrgByName(ctx context.Context, orgName string) (*model.Organization, error) {
	var org model.Organization
	if err := s.orgRepo.DB().Where("name = ?", orgName).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (s *OrganizationService) CheckMemberRole(ctx context.Context, orgID, userID uint) (string, error) {
	var member model.OrganizationMember
	if err := s.orgRepo.DB().Where("organization_id = ? AND user_id = ?", orgID, userID).First(&member).Error; err != nil {
		return "", err
	}
	return member.Role, nil
}

func (s *OrganizationService) UpdateMember(ctx context.Context, orgID uint, username string, req *dto.UpdateMemberRequest) error {
	// 根据用户名获取用户ID
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("用户不存在")
	}

	// 更新成员信息
	updates := map[string]interface{}{
		"role": req.Role,
	}
	if err := s.orgRepo.UpdateMember(ctx, orgID, user.ID, updates); err != nil {
		return err
	}

	return nil
}
