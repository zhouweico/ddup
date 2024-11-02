package service

import (
	"context"
	"ddup-apis/internal/dto"
	"ddup-apis/internal/model"
	"ddup-apis/internal/repository"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type ProfileService struct {
	repo *repository.ProfileRepository
}

func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{
		repo: repository.NewProfileRepository(db),
	}
}

func (s *ProfileService) Create(ctx context.Context, userID uint, req *dto.CreateProfileRequest) error {
	profile := &model.Profile{
		UserID:       userID,
		Type:         model.ProfileType(req.Type),
		Title:        req.Title,
		Year:         req.Year,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Organization: req.Organization,
		Location:     req.Location,
		URL:          req.URL,
		Description:  req.Description,
		Metadata:     json.RawMessage(req.Metadata),
		Visibility:   req.Visibility,
	}
	return s.repo.Create(ctx, profile)
}

func (s *ProfileService) GetByID(ctx context.Context, userID, profileID uint) (*dto.ProfileResponse, error) {
	profile, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	if profile.UserID != userID {
		return nil, errors.New("无权访问此资料")
	}

	return s.toProfileResponse(profile), nil
}

func (s *ProfileService) GetByType(ctx context.Context, userID uint, profileType string) ([]dto.ProfileResponse, error) {
	profiles, err := s.repo.GetByType(ctx, userID, profileType)
	if err != nil {
		return nil, err
	}

	var resp []dto.ProfileResponse
	for _, p := range profiles {
		resp = append(resp, *s.toProfileResponse(&p))
	}
	return resp, nil
}

func (s *ProfileService) Update(ctx context.Context, userID, profileID uint, req *dto.UpdateProfileRequest) error {
	profile, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return err
	}

	if profile.UserID != userID {
		return errors.New("无权修改此资料")
	}

	if req.Title != "" {
		profile.Title = req.Title
	}
	if req.Year != nil {
		profile.Year = req.Year
	}
	// ... 更新其他字段

	return s.repo.Update(ctx, profile)
}

func (s *ProfileService) Delete(ctx context.Context, userID, profileID uint) error {
	profile, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return err
	}

	if profile.UserID != userID {
		return errors.New("无权删除此资料")
	}

	return s.repo.Delete(ctx, profileID)
}

func (s *ProfileService) UpdateDisplayOrder(ctx context.Context, userID uint, req *dto.UpdateDisplayOrderRequest) error {
	// 验证所有项目属于当前用户
	for _, item := range req.Items {
		profile, err := s.repo.GetByID(ctx, item.ID)
		if err != nil {
			return err
		}
		if profile.UserID != userID {
			return errors.New("无权修改此资料")
		}
	}

	var items []struct {
		ID    uint
		Order int
	}
	for _, item := range req.Items {
		items = append(items, struct {
			ID    uint
			Order int
		}{
			ID:    item.ID,
			Order: item.Order,
		})
	}

	return s.repo.UpdateDisplayOrder(ctx, items)
}

func (s *ProfileService) toProfileResponse(p *model.Profile) *dto.ProfileResponse {
	return &dto.ProfileResponse{
		ID:           p.ID,
		Type:         string(p.Type),
		Title:        p.Title,
		Year:         p.Year,
		StartDate:    p.StartDate,
		EndDate:      p.EndDate,
		Organization: p.Organization,
		Location:     p.Location,
		URL:          p.URL,
		Description:  p.Description,
		Metadata:     json.RawMessage(p.Metadata),
		DisplayOrder: p.DisplayOrder,
		Visibility:   p.Visibility,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}
