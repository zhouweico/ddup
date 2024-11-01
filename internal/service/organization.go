package service

import (
	"context"
	"ddup-apis/internal/dto"
	"ddup-apis/internal/errors"
	"ddup-apis/internal/model"
	"ddup-apis/internal/repository"

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

func (s *OrganizationService) CreateOrganization(ctx context.Context, req *dto.CreateOrganizationRequest, userID uint) (*dto.OrganizationResponse, error) {
	org := &model.Organization{
		Name:        req.Name,
		Email:       req.Email,
		Avatar:      req.Avatar,
		Description: req.Description,
		Location:    req.Location,
		Website:     req.Website,
	}

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, errors.Wrap(err, "创建组织失败")
	}

	member := &model.OrganizationMember{
		OrganizationID: org.ID,
		UserID:         userID,
		Role:           "admin",
	}

	if err := s.orgRepo.AddMember(ctx, member); err != nil {
		return nil, errors.Wrap(err, "添加组织管理员失败")
	}

	return &dto.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Email:       org.Email,
		Avatar:      org.Avatar,
		Description: org.Description,
		Location:    org.Location,
		Website:     org.Website,
		Role:        "admin",
	}, nil
}

func (s *OrganizationService) GetUserOrganizations(ctx context.Context, userID uint) ([]dto.OrganizationResponse, error) {
	orgs, err := s.orgRepo.GetUserOrganizations(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "获取用户组织列表失败")
	}

	var resp []dto.OrganizationResponse
	for _, org := range orgs {
		member, err := s.orgRepo.GetMember(ctx, org.ID, userID)
		if err != nil {
			return nil, errors.Wrap(err, "获取成员信息失败")
		}

		resp = append(resp, dto.OrganizationResponse{
			ID:          org.ID,
			Name:        org.Name,
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

func (s *OrganizationService) UpdateOrganization(ctx context.Context, id uint, req *dto.UpdateOrganizationRequest) error {
	updates := map[string]interface{}{
		"name":        req.Name,
		"email":       req.Email,
		"avatar":      req.Avatar,
		"description": req.Description,
		"location":    req.Location,
		"website":     req.Website,
	}

	return s.orgRepo.Update(ctx, id, updates)
}

func (s *OrganizationService) DeleteOrganization(ctx context.Context, id uint) error {
	return s.orgRepo.Delete(ctx, id)
}

func (s *OrganizationService) AddMember(ctx context.Context, orgID uint, req *dto.AddMemberRequest) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return errors.Wrap(err, "用户不存在")
	}

	// 检查是否已经是成员
	member, err := s.orgRepo.GetMember(ctx, orgID, user.ID)
	if err == nil && member != nil {
		return errors.New(400, "用户已经是组织成员", nil)
	}

	newMember := &model.OrganizationMember{
		OrganizationID: orgID,
		UserID:         user.ID,
		Role:           req.Role,
	}

	return s.orgRepo.AddMember(ctx, newMember)
}

func (s *OrganizationService) GetMembers(ctx context.Context, orgID uint) ([]dto.MemberResponse, error) {
	members, err := s.orgRepo.GetMembers(ctx, orgID)
	if err != nil {
		return nil, errors.Wrap(err, "获取成员列表失败")
	}

	var resp []dto.MemberResponse
	for _, member := range members {
		user, err := s.userRepo.GetByID(ctx, member.UserID)
		if err != nil {
			return nil, errors.Wrap(err, "获取用户信息失败")
		}

		resp = append(resp, dto.MemberResponse{
			UserID:    user.ID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Role:      member.Role,
			CreatedAt: member.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

func (s *OrganizationService) UpdateMember(ctx context.Context, orgID, userID uint, req *dto.UpdateMemberRequest) error {
	return s.orgRepo.UpdateMember(ctx, orgID, userID, req.Role)
}

func (s *OrganizationService) RemoveMember(ctx context.Context, orgID, userID uint) error {
	return s.orgRepo.RemoveMember(ctx, orgID, userID)
}

func (s *OrganizationService) CheckMemberRole(ctx context.Context, orgID, userID uint) (string, error) {
	member, err := s.orgRepo.GetMember(ctx, orgID, userID)
	if err != nil {
		return "", errors.Wrap(err, "获取成员信息失败")
	}
	return member.Role, nil
}
